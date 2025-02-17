package types

import (
	"net/http"
	"net/url"
	"sort"

	"golang.org/x/net/context"
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

type ServiceContext struct {
	Method       string        // http method
	Request      *http.Request // http request
	URL          *url.URL      // request url
	Header       http.Header   // http request header
	WriterHeader http.Header   // http writer header
	ClientIP     string        // client ip
	UserAgent    string        // user agent

	SessionId string // session id
	Username  string // currrent login user.
	UserId    string // currrent login user id
	Context   context.Context

	RequestId string
	TraceId   string
	SpanId    string
	PSpanId   string
}

type ControllerContext struct {
	Username string // currrent login user.
	UserId   string // currrent login user id

	RequestId string
	TraceId   string
	SpanId    string
	PSpanId   string
}

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
}
