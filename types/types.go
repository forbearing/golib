package types

import (
	"net/http"
	"net/url"
	"sort"

	"github.com/forbearing/golib/types/consts"
)

type sortByCreatedTime[M Model] []M

func (s sortByCreatedTime[M]) Len() int      { return len(s) }
func (s sortByCreatedTime[M]) Swap(x, y int) { s[x], s[y] = s[y], s[x] }
func (s sortByCreatedTime[M]) Less(x, y int) bool {
	return s[x].GetCreatedAt().After(s[y].GetCreatedAt())
}

type sortByUpdatedTime[M Model] []M

func (s sortByUpdatedTime[M]) Len() int      { return len(s) }
func (s sortByUpdatedTime[M]) Swap(x, y int) { s[x], s[y] = s[y], s[x] }
func (s sortByUpdatedTime[M]) Less(x, y int) bool {
	return s[x].GetUpdatedAt().After(s[y].GetUpdatedAt())
}

type Order int

const (
	UpdatedTime Order = iota
	CreatedTime
)

func Sort[M Model](order Order, data []M, reverse ...bool) {
	var _reverse bool
	if len(reverse) > 0 {
		_reverse = reverse[0]
	}
	_sort := func(data sort.Interface) {
		if _reverse {
			sort.Sort(sort.Reverse(data))
		} else {
			sort.Sort(data)
		}
	}

	switch order {
	case CreatedTime:
		_sort(sortByCreatedTime[M](data))
	case UpdatedTime:
		_sort(sortByUpdatedTime[M](data))
	}
}

type ControllerContext struct {
	Username string // currrent login user.
	UserId   string // currrent login user id
	Route    string
	Params   map[string]string
	Query    map[string][]string

	RequestId string
	TraceId   string
	PSpanId   string
	SpanId    string
	Seq       int
}

type DatabaseContext struct {
	Username string // currrent login user.
	UserId   string // currrent login user id
	Route    string
	Params   map[string]string
	Query    map[string][]string

	RequestId string
	TraceId   string
	PSpanId   string
	SpanId    string
	Seq       int
}

type ServiceContext struct {
	Method       string        // http method
	Request      *http.Request // http request
	URL          *url.URL      // request url
	Header       http.Header   // http request header
	WriterHeader http.Header   // http writer header
	ClientIP     string        // client ip
	UserAgent    string        // user agent

	// route parameters,
	//
	// eg: PUT /api/gists/:id/star
	// Params: map[string]string{"id": "xxxxx-mygistid-xxxxx"}
	//
	// eg: DELETE /api/user/:userid/shelf/shelfid/book
	// Params: map[string]string{"userid": "xxxxx-myuserid-xxxxx", "shelfid": "xxxxx-myshelfid-xxxxx"}
	Params map[string]string
	Query  map[string][]string

	SessionId string // session id
	Username  string // currrent login user.
	UserId    string // currrent login user id
	Route     string

	RequestId string
	TraceId   string
	PSpanId   string
	SpanId    string
	Seq       int

	phase consts.Phase

	request  any // custom http request.
	response any // custom http response.
}

// SetRequest is called in the controller layer when the model has custom request type.
// The custom request only take effect for CREATE(POST), UPDATE(PUT), UPDATE_PARTIAL(PATCH)
// CreateMany(POST), UpdateMany(PUT) and PatchMany(PATCH) operations.
func (sc *ServiceContext) SetRequest(m any) { sc.request = m }

// SetResponse is called in the controller layer if the model has custom response type.
// The custom response only take effect for Create(POST), Update(PUT), Patch(PATCH)
// CreateMany(POST), UpdateMany(PUT) and PatchMany(PATCH) operations.
func (sc *ServiceContext) SetResponse(m any) { sc.response = m }

// GetRequest is called in the service layer when the model has a custom request type.
//
// It returns the custom request object unmarshaled from the HTTP request body.
// This method is only supported in the following service hooks:
//
//	CREATE_BEFORE (POST), CREATE_AFTER (POST)
//	UPDATE_BEFORE (PUT), UPDATE_AFTER (PUT)
//	PATCH_BEFORE (PATCH), PATCH_AFTER (PATCH)
//	CREATE_MANY_BEFORE (POST), CREATE_MANY_AFTER (POST)
//	UPDATE_MANY_BEFORE (PUT), UPDATE_MANY_AFTER (PUT)
//	PATCH_MANY_BEFORE (PATCH), PATCH_MANY_AFTER (PATCH)
//
// For all other service hooks, GetRequest always returns nil, including:
//
//	DELETE_BEFORE (DELETE), DELETE_AFTER (DELETE)
//	LIST_BEFORE (GET), LIST_AFTER (GET)
//	GET_BEFORE (GET), GET_AFTER (GET)
func (sc *ServiceContext) GetRequest() any { return sc.request }

// GetResponse is typically called in the controller layer to return a custom response object to the client.
// The response object is set in the following service hooks:
//
//	CREATE_BEFORE (POST), CREATE_AFTER (POST)
//	UPDATE_BEFORE (PUT), UPDATE_AFTER (PUT)
//	PATCH_BEFORE (PATCH), PATCH_AFTER (PATCH)
//	CREATE_MANY_BEFORE (POST), CREATE_MANY_AFTER (POST)
//	UPDATE_MANY_BEFORE (PUT), UPDATE_MANY_AFTER (PUT)
//	PATCH_MANY_BEFORE (PATCH), PATCH_MANY_AFTER (PATCH)
func (sc *ServiceContext) GetResponse() any { return sc.response }

func (sc *ServiceContext) SetPhase(phase consts.Phase) { sc.phase = phase }
func (sc *ServiceContext) GetPhase() consts.Phase      { return sc.phase }
func (sc *ServiceContext) WithPhase(phase consts.Phase) *ServiceContext {
	sc.phase = phase
	return sc
}

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
}
