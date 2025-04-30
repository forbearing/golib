package types

import (
	"net/http"
	"net/url"
	"sort"
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

	requestBody  any
	responseBody any
}

func (sc *ServiceContext) SetRequestBody(m any)  { sc.requestBody = m }
func (sc *ServiceContext) SetResponseBody(m any) { sc.responseBody = m }
func (sc *ServiceContext) GetRequestBody() any   { return sc.requestBody }
func (sc *ServiceContext) GetResponseBody() any  { return sc.responseBody }

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
}
