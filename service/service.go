package service

import (
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	// WARN: 一定不要使用 reflect.TypeOf(*new(M)).Name(), 因为可能存在 model.User, model2.User 的情况,
	// 这样就会导致 key 重复.
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
		logger.Service.Debugz(ErrNotFoundService.Error(), zap.String("model", reflect.TypeOf(*new(M)).String()))
		return new(Base[M])
	}
	return svc.(types.Service[M])
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
func (Base[M]) FilterRaw(_ *types.ServiceContext) string              { return "" }
