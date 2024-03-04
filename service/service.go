package service

import (
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	serviceMap = make(map[string]any)
	mu         sync.Mutex
)

var (
	ErrNotFoundService   = errors.New("no service instant matches the give Model interface, skip processing service layer")
	ErrNotFoundServiceId = errors.New("not found service id in assetIdMap")
)

// Register service instance into serviceMap.
// pass parameters to replace the default service instance.
// If the passed parameters is nil, skip replace.
func Register[S types.Service[M], M types.Model](s ...S) {
	mu.Lock()
	defer mu.Unlock()
	key := reflect.TypeOf(*new(M)).String()
	val := reflect.New(reflect.TypeOf(*new(S)).Elem()).Interface()
	serviceMap[key] = val
	if len(s) > 0 {
		// replace the default service instance if the passed service instance isn't nil.
		if !reflect.ValueOf(s[0]).IsNil() {
			serviceMap[key] = s[0]
		}
	}
}

func Init() error {
	// Init all service types.Logger
	{
		for _, s := range serviceMap {
			typ := reflect.TypeOf(s).Elem()
			val := reflect.ValueOf(s).Elem()
			for i := 0; i < typ.NumField(); i++ {
				switch strings.ToLower(typ.Field(i).Name) {
				case "logger": // service object has itself types.Logger
					if val.Field(i).IsZero() {
						val.Field(i).Set(reflect.ValueOf(logger.Service))
					}
				case "base": // service object's types.Logger extends from 'base' struct.
					fieldLogger := val.Field(i).FieldByName("Logger")
					if fieldLogger.IsZero() {
						fieldLogger.Set(reflect.ValueOf(logger.Service))
					}
				}
			}
		}
	}

	return nil
}

// Factory is a service factory used to product service instance.
// The servicei instance should registered by function `register()` in init()
//
// The service defined by user should be unexported (structure name is lowercase).
// service instance are only returns by the `Factory`.
type Factory[M types.Model] struct{}

func (f Factory[M]) Service() types.Service[M] {
	svc, ok := serviceMap[reflect.TypeOf(*new(M)).String()]
	if !ok {
		logger.Service.Debugw(ErrNotFoundService.Error(), "model", reflect.TypeOf(*new(M)).String())
		return new(Base[M])
	}
	return svc.(types.Service[M])
}

// GinContext build *types.ServiceContext from *gin.Context.
func GinContext(c *gin.Context) *types.ServiceContext {
	var requestId string
	val, _ := c.Get(types.REQUEST_ID)
	switch v := val.(type) {
	case string:
		requestId = v
	}

	return &types.ServiceContext{
		Method:    c.Request.Method,
		URL:       c.Request.URL,
		Header:    c.Writer.Header(),
		Username:  c.GetString(types.CTX_USERNAME),
		UserId:    c.GetString(types.CTX_USER_ID),
		SessionId: c.GetString(types.CTX_SESSION_ID),
		RequestId: requestId,
	}
}

type Base[M types.Model] struct{ types.Logger }

func (Base[M]) CreateBefore(*types.ServiceContext, ...M) error        { return nil }
func (Base[M]) CreateAfter(*types.ServiceContext, ...M) error         { return nil }
func (Base[M]) DeleteBefore(*types.ServiceContext, ...M) error        { return nil }
func (Base[M]) DeleteAfter(*types.ServiceContext, ...M) error         { return nil }
func (Base[M]) UpdateBefore(*types.ServiceContext, ...M) error        { return nil }
func (Base[M]) UpdateAfter(*types.ServiceContext, ...M) error         { return nil }
func (Base[M]) UpdatePartialBefore(*types.ServiceContext, ...M) error { return nil }
func (Base[M]) UpdatePartialAfter(*types.ServiceContext, ...M) error  { return nil }
func (Base[M]) ListBefore(*types.ServiceContext, *[]M) error          { return nil }
func (Base[M]) ListAfter(*types.ServiceContext, *[]M) error           { return nil }
func (Base[M]) GetBefore(*types.ServiceContext, ...M) error           { return nil }
func (Base[M]) GetAfter(*types.ServiceContext, ...M) error            { return nil }
func (Base[M]) Import(*types.ServiceContext, io.Reader) ([]M, error)  { return make([]M, 0), nil }
func (Base[M]) Export(*types.ServiceContext, ...M) ([]byte, error)    { return make([]byte, 0), nil }
func (Base[M]) Filter(_ *types.ServiceContext, m M) M                 { return m }
