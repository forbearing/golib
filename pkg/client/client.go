package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/types"
	"github.com/google/go-querystring/query"
	"golang.org/x/time/rate"
)

type Client struct {
	addr       string
	httpClient *http.Client
	username   string
	password   string
	token      string

	header      http.Header
	query       *model.Base
	queryRaw    string
	param       string
	debug       bool
	maxRetries  int
	retryWait   time.Duration
	rateLimiter *rate.Limiter

	ctx context.Context

	types.Logger
}

type action int

const (
	create action = iota
	delete_
	update
	update_partial
	list
	get
)

type response struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// New creates a new client instance with given base URL and options.
// The base URL must start with "http://" or "https://".
func New(addr string, opts ...Option) (*Client, error) {
	client := &Client{
		httpClient: http.DefaultClient,
		header:     http.Header{},
		addr:       strings.TrimRight(addr, "/"),
		ctx:        context.Background(),
		Logger:     zap.New(),
	}
	client.header.Set("User-Agent", "golib")
	client.header.Set("Content-Type", "application/json")

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(client)
	}

	return client, nil
}

// QueryString build the query string from structured query parameters
// and raw query string.
func (c *Client) QueryString() (string, error) {
	if c.query == nil && len(c.queryRaw) == 0 {
		return "", nil
	}
	if c.query == nil {
		return c.queryRaw, nil
	}

	val, err := query.Values(c.query)
	if err != nil {
		return "", err
	}

	encoded := val.Encode()
	if len(encoded) == 0 {
		return c.queryRaw, nil
	}
	if len(c.queryRaw) != 0 {
		return c.queryRaw + "&" + encoded, nil
	}

	return encoded, nil
}

// RequestURL constructs the full request URL including base URL and query parameters.
func (c *Client) RequestURL() (string, error) {
	if !strings.HasPrefix(c.addr, "http://") && !strings.HasPrefix(c.addr, "https://") {
		return "", errors.New("addr must start with http:// or https://")
	}
	query, err := c.QueryString()
	if err != nil {
		return "", err
	}
	if len(query) > 0 {
		return fmt.Sprintf("%s?%s", c.addr, query), nil
	}
	return c.addr, nil
}

// Create send a POST request to create a new resource.
// payload can be []byte or struct/pointer that can be marshaled to JSON.
func (c *Client) Create(payload any) ([]byte, error) {
	return c.request(create, payload)
}

// Delete send a DELETE request to delete a resource.
func (c *Client) Delete(id string) ([]byte, error) {
	if len(id) == 0 {
		return nil, errors.New("id is required")
	}
	c.param = id
	return c.request(delete_, nil)
}

// Update send a PUT request to fully update a resource.
func (c *Client) Update(payload any) ([]byte, error) {
	return c.request(update, payload)
}

// UpdatePartial send a PATCH request to partially update a resource.
func (c *Client) UpdatePartial(payload any) ([]byte, error) {
	return c.request(update_partial, payload)
}

// List send a GET request to retrieve a list of resources.
// items must be a pointer to slice where items will be unmarshaled into.
// total will be set to the total number of items available.
func (c *Client) List(items any, total *int64) error {
	if items == nil {
		return errors.New("items cannot be nil")
	}
	if total == nil {
		return errors.New("total cannot be nil")
	}

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Ptr {
		return errors.New("items must be a pointer to slice")
	}
	if val.Elem().Kind() != reflect.Slice {
		return errors.New("items must be a pointer to slice")
	}
	data, err := c.request(list, nil)
	if err != nil {
		return err
	}
	responseList := new(struct {
		Items json.RawMessage `json:"items"`
		Total int64           `json:"total"`
	})
	if err := json.Unmarshal(data, responseList); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}
	if err := json.Unmarshal(responseList.Items, items); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}
	*total = responseList.Total
	return nil
}

// ListRaw send a GET request to retrieve a list of resources without unmarshaling.
func (c *Client) ListRaw() ([]byte, error) {
	return c.request(list, nil)
}

// Get send a GET request to get one resource by given id.
// The id parameter specifies which resource to retrieve.
// The dst parameter must be a pointer to struct where the resource will be unmarshaled into.
func (c *Client) Get(id string, dst any) error {
	if len(id) == 0 {
		return errors.New("id is required")
	}
	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer to struct")
	}
	if val.Elem().Kind() != reflect.Struct {
		return errors.New("dst must be a pointer to struct")
	}
	if !val.Elem().IsZero() {
		newVal := reflect.New(reflect.TypeOf(dst).Elem())
		val.Elem().Set(newVal.Elem())
	}
	c.param = id
	data, err := c.request(get, nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}
	return nil
}

// GetRaw send a GET request to get one resource by given id without unmarshaling.
func (c *Client) GetRaw() ([]byte, error) {
	return c.request(get, nil)
}

// request send a request to backend server.
// action determines the type of request,
// payload can be []byte or struct/pointer that can be marshaled to JSON.
func (c *Client) request(action action, payload any) ([]byte, error) {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(c.ctx); err != nil {
			return nil, errors.Wrap(err, "rate limit exceeded")
		}
	}

	var url string
	var err error
	var method string
	switch action {
	case create:
		method = http.MethodPost
		url = c.addr
	case delete_:
		method = http.MethodDelete
		url = fmt.Sprintf("%s/%s", c.addr, c.param)
	case update:
		method = http.MethodPut
		url = c.addr
	case update_partial:
		method = http.MethodPatch
		url = c.addr
	case list:
		method = http.MethodGet
		url, err = c.RequestURL()
	case get:
		method = http.MethodGet
		url = fmt.Sprintf("%s/%s", c.addr, c.param)
	}
	if err != nil {
		return nil, errors.Wrap(err, "invalid request url")
	}

	var _payload io.Reader
	if payload != nil {
		switch v := payload.(type) {
		case []byte:
			_payload = bytes.NewReader(v)
		default:
			var data []byte
			if data, err = json.Marshal(v); err != nil {
				return nil, errors.Wrap(err, "failed to marshal payload")
			}
			_payload = bytes.NewReader(data)
		}
	}
	req, err := http.NewRequest(method, url, _payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	}
	if len(c.username) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}
	if len(c.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	maps.Copy(req.Header, c.header)

	if c.debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, errors.Wrap(err, "failed to copy response body")
	}
	if c.debug {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(dump))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return buf.Bytes(), fmt.Errorf("response status code: %d, body: %s", resp.StatusCode, buf.String())
	}

	res := new(response)
	if err := json.Unmarshal(buf.Bytes(), res); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("response status code: %d, code: %d, msg: %s, body: %s", resp.StatusCode, res.Code, res.Msg, buf.String())
	}
	return res.Data, nil
}
