package httpwrapper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/forbearing/gst/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrappedRequest(t *testing.T) {
	body := []byte("hello world")
	req, err := http.NewRequest("POST", "http://example.com", bytes.NewReader(body))
	assert.NoError(t, err)
	testWrappedRequest(t, req, body)

	body = []byte{}
	req, err = http.NewRequest("POST", "http://example.com", bytes.NewReader(body))
	require.NoError(t, err)
	testWrappedRequest(t, req, body)
}

func testWrappedRequest(t *testing.T, req *http.Request, body []byte) {
	reqWrapper := &WrappedRequest{Request: req}
	data, err := json.Marshal(reqWrapper)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	t.Log(util.StringAny(data))

	reqWrapper = new(WrappedRequest)
	err = json.Unmarshal(data, reqWrapper)
	assert.NoError(t, err)
	data, err = io.ReadAll(reqWrapper.Body)
	assert.NoError(t, err)
	assert.Equal(t, data, body)
}

func TestWrappedResponse(t *testing.T) {
	domains := []string{
		"http://example.com",
		"https://example.com",
	}

	for _, domain := range domains {
		req, err := http.NewRequest(http.MethodGet, domain, nil)
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		require.NoError(t, err)
		testWrappedResponse(t, resp, body)
		resp.Body.Close()
	}
}

func testWrappedResponse(t *testing.T, resp *http.Response, body []byte) {
	respWrapper := &WrappedResponse{Response: resp}
	data, err := json.Marshal(respWrapper)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	// t.Log(internal.String(data))

	respWrapper = new(WrappedResponse)
	err = json.Unmarshal(data, respWrapper)
	assert.NoError(t, err)
	data, err = io.ReadAll(respWrapper.Body)
	assert.NoError(t, err)
	assert.Equal(t, data, body)
}
