package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/ds/queue/circularbuffer"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
	model_log "github.com/forbearing/golib/model/log"
	"github.com/forbearing/golib/pkg/filetype"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/types/helper"
	"github.com/forbearing/golib/util"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

// TODO: 记录失败的操作.

/*
1.Model 层处理单个 types.Model, 功能: 数据预处理
2.Service 层处理多个 types.Model, 功能: 具体的业务逻辑
3.Database 层处理多个 types.Model, 功能: 数据库的增删改查,redis缓存等.
4.这三层都能对源数据进行修改, 因为:
  - Model 的实现对象必须是结构体指针
  - types.Service[M types.Model]: types.Service 泛型接口的类型约束是 types.Model 接口
  - types.Database[M types.Model]: types.Database 泛型接口的类型约束就是 types.Model 接口
  以上这三个条件自己慢慢体会吧.
5.用户自定义的 model:
  必须继承 model.Base 结构体, 因为这个结构体实现了 types.Model 接口
  用户只需要添加自己的字段和相应的 tag 和方法即可.
  如果想要给 types.Model 在数据库中创建对象的表, 请在 init() 中调用 register 函数注册一下即可, 比如 register[*Asset]()
  如果需要在创建表格的同时创建记录, 也可以通过 register 函数来做, 比如 register[*Asset](asset01, asset02, asset03)
  这里的 asset01, asset02, asset03 的类型是 *model.Asset.
6.用户自定义 service
  必须继承 service.base 结构体, 因为这个结构体实现了 types.Service[types.Model] 接口
  用户只需要覆盖默认的方法就行了
如果有额外的业务逻辑, 在 init() 中调用 register 函数注册一下自己定义的 service, 例如: register[*asset, *model.Asset](new(asset))
如果 service.Asset 有自定义字段, 可以这样: register[*asset, *model.Asset](&asset{SheetName: "资产类别清单"})

处理资源顺序:
    通用流程: Request -> ServiceBefore -> ModelBefore -> Database -> ModelAfter -> ServiceAfter -> Response.
	导入数据: Request -> ServiceBefore -> Import -> ModelBefore ->  Database -> ModelAfter -> ServiceAfter -> Response.
	导出数据: Request -> ServiceBefore -> ModelBefore -> Database -> ModelAfter -> ServiceAfter -> Export -> Response.

    Import 逻辑类似于 Update 逻辑
	Import 的 Model 的 UpdateBefore() 在 service 层里面处理, ServiceBefore 是可选的
	Export 逻辑类似于 List 逻辑, 只是比 Update 逻辑多了 Export 步骤

其他:
	1.记录操作日志也在 controller 层
*/

var (
	pluralizeCli = pluralize.NewClient()

	cb *circularbuffer.CircularBuffer[*model_log.OperationLog]
)

func Init() (err error) {
	if cb, err = circularbuffer.New(int(config.App.Server.CircularBuffer.SizeOperationLog), circularbuffer.WithSafe[*model_log.OperationLog]()); err != nil {
		return err
	}

	// Consume operation log.
	go func() {
		operationLogs := make([]*model_log.OperationLog, 0, config.App.Server.CircularBuffer.SizeOperationLog)
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			operationLogs = operationLogs[:0]
			for !cb.IsEmpty() {
				ol, _ := cb.Dequeue()
				operationLogs = append(operationLogs, ol)
			}
			if len(operationLogs) > 0 {
				if err := database.Database[*model_log.OperationLog]().WithLimit(-1).WithBatchSize(1000).Create(operationLogs...); err != nil {
					zap.S().Error(err)
				}
			}
		}
	}()

	return nil
}

func Clean() {
	operationLogs := make([]*model_log.OperationLog, 0, config.App.Server.CircularBuffer.SizeOperationLog)
	for !cb.IsEmpty() {
		ol, _ := cb.Dequeue()
		operationLogs = append(operationLogs, ol)
	}
	if len(operationLogs) > 0 {
		if err := database.Database[*model_log.OperationLog]().WithLimit(-1).WithBatchSize(100).Create(operationLogs...); err != nil {
			zap.S().Error(err)
		}
	}
}

// Create is a generic function to product gin handler to create one resource.
// The resource type depends on the type of interface types.Model.
func Create[M types.Model](c *gin.Context) {
	CreateFactory[M]()(c)
}

// CreateFactory is a factory function to product gin handler to create one resource.
func CreateFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		var err error
		var reqErr error
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_CREATE)
		typ := reflect.TypeOf(*new(M)).Elem()
		req := reflect.New(typ).Interface().(M)
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request.
			req := model.NewRequest[M]()
			if reqErr = c.ShouldBindJSON(req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeInvalidParam.WithErr(reqErr))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				// Has custom response
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeInvalidParam.WithErr(reqErr))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			} else {
				req.SetCreatedBy(c.GetString(consts.CTX_USERNAME))
				req.SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
				log.Infoz("create", zap.Object(reflect.TypeOf(*new(M)).Elem().String(), req))
			}
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before create resource.
		if err = svc.CreateBefore(ctx.WithPhase(consts.PHASE_CREATE_BEFORE), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Create resource in database.
		// database.Database().Delete just set "deleted_at" field to current time, not really delete.
		// We should update it instead of creating it, and update the "created_at" and "updated_at" field.
		// NOTE: WithExpand(req.Expands()...) is not a good choices.
		// if err := database.Database[M]().WithExpand(req.Expands()...).Update(req); err != nil {
		if reqErr != io.EOF && !hasCustomReq {
			if err = handler(helper.NewDatabaseContext(c)).WithExpand(req.Expands()).Create(req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after create resource
		if err = svc.CreateAfter(ctx.WithPhase(consts.PHASE_CREATE_AFTER), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetRequest())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeCreate,
			Model:     typ.Name(),
			Table:     tableName,
			RecordId:  req.GetID(),
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated), ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated))
		} else {
			// Does not have custom response, the response type is the same as the resource type.
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated), req)
		}
	}
}

// Delete is a generic function to product gin handler to delete one or multiple resources.
// The resource type depends on the type of interface types.Model.
//
// Resource id must be specify and all resources that id matched will be deleted in database.
//
// Delete one resource:
// - specify resource `id` in "router parameter", eg: localhost:9000/api/myresource/myid
// - specify resource `id` in "query parameter", eg: localhost:9000/api/myresource?id=myid
//
// Delete multiple resources:
// - specify resource `id` slice in "http body data".
func Delete[M types.Model](c *gin.Context) {
	DeleteFactory[M]()(c)
}

// DeleteFactory is a factory function to product gin handler to delete one or multiple resources.
func DeleteFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_DELETE)
		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		typ := reflect.TypeOf(*new(M)).Elem()
		ml := make([]M, 0)
		idsSet := make(map[string]struct{})

		addId := func(id string) {
			if len(id) == 0 {
				return
			}
			if _, exists := idsSet[id]; exists {
				return
			}
			// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)
			ml = append(ml, m)
			idsSet[id] = struct{}{}
		}

		// Delete one record accoding to "query parameter `id`".
		if id, ok := c.GetQuery(consts.QUERY_ID); ok {
			addId(id)
		}
		// Delete one record accoding to "route parameter `id`".
		if id := c.Param(consts.PARAM_ID); len(id) != 0 {
			addId(id)
		}
		// Delete multiple records accoding to "http body data".
		bodyIds := make([]string, 0)
		if err := c.ShouldBindJSON(&bodyIds); err == nil && len(bodyIds) > 0 {
			for _, id := range bodyIds {
				addId(id)
			}
		}

		ids := make([]string, 0, len(idsSet))
		for id := range idsSet {
			ids = append(ids, id)
		}
		log.Info(fmt.Sprintf("%s delete %v", typ.Name(), ids))

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before delete resources.
		// TODO: Should there be one service hook(DeleteBefore), or multiple?
		for _, m := range ml {
			if err := svc.DeleteBefore(helper.NewServiceContext(c).WithPhase(consts.PHASE_DELETE_BEFORE), m); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}

		// find out the records and record to operation log.
		copied := make([]M, len(ml))
		for i := range ml {
			m := reflect.New(typ).Interface().(M)
			m.SetID(ml[i].GetID())
			if err := handler(helper.NewDatabaseContext(c)).WithExpand(m.Expands()).Get(m, ml[i].GetID()); err != nil {
				log.Error(err)
			}
			copied[i] = m
		}

		// 2.Delete resources in database.
		if err := handler(helper.NewDatabaseContext(c)).Delete(ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 3.Perform business logic processing after delete resources.
		// TODO: Should there be one service hook(DeleteAfter), or multiple?
		for _, m := range ml {
			if err := svc.DeleteAfter(helper.NewServiceContext(c).WithPhase(consts.PHASE_DELETE_AFTER), m); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		for i := range ml {
			record, _ := json.Marshal(copied[i])
			cb.Enqueue(&model_log.OperationLog{
				Op:        model_log.OperationTypeDelete,
				Model:     typ.Name(),
				Table:     tableName,
				RecordId:  ml[i].GetID(),
				Record:    util.BytesToString(record),
				IP:        c.ClientIP(),
				User:      c.GetString(consts.CTX_USERNAME),
				RequestId: c.GetString(consts.REQUEST_ID),
				URI:       c.Request.RequestURI,
				Method:    c.Request.Method,
				UserAgent: c.Request.UserAgent(),
			})
		}

		ResponseJSON(c, CodeSuccess.WithStatus(http.StatusNoContent))
	}
}

// Update is a generic function to product gin handler to update one resource.
// The resource type depends on the type of interface types.Model.
//
// Update will update one resource and resource "ID" must be specified,
// which can be specify in "router parameter `id`" or "http body data".
//
// "router parameter `id`" has more priority than "http body data".
// It will skip decode id from "http body data" if "router parameter `id`" not empty.
func Update[M types.Model](c *gin.Context) {
	UpdateFactory[M]()(c)
}

// UpdateFactory is a factory function to product gin handler to update one resource.
func UpdateFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_UPDATE)
		typ := reflect.TypeOf(*new(M)).Elem()
		req := reflect.New(typ).Interface().(M)
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request.
			req := model.NewRequest[M]()
			if err := c.ShouldBindJSON(req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}

			// param id has more priority than http body data id
			paramId := c.Param(consts.PARAM_ID)
			bodyId := req.GetID()
			var id string
			log.Infoz("update from request",
				zap.String("param_id", paramId),
				zap.String("body_id", bodyId),
				zap.Object(reflect.TypeOf(*new(M)).Elem().String(), req),
			)
			if paramId != "" {
				req.SetID(paramId)
				id = paramId
			} else if bodyId != "" {
				paramId = bodyId
				id = bodyId
			} else {
				log.Error("id missing")
				ResponseJSON(c, CodeFailure.WithErr(errors.New("id missing")))
				return
			}

			data := make([]M, 0)
			// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
			// 'typ' is the structure type, such as: model.User.
			// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)
			// Make sure the record must be already exists.
			if err := handler(helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(m).List(&data); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if len(data) != 1 {
				log.Errorz(fmt.Sprintf("the total number of records query from database not equal to 1(%d)", len(data)), zap.String("id", id))
				ResponseJSON(c, CodeNotFound)
				return
			}

			req.SetCreatedAt(data[0].GetCreatedAt())           // keep original "created_at"
			req.SetCreatedBy(data[0].GetCreatedBy())           // keep original "created_by"
			req.SetUpdatedBy(c.GetString(consts.CTX_USERNAME)) // set updated_by to current user”
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before update resource.
		if err := svc.UpdateBefore(ctx.WithPhase(consts.PHASE_UPDATE_BEFORE), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Update resource in database.
		if !hasCustomReq {
			log.Infoz("update in database", zap.Object(typ.Name(), req))
			if err := handler(helper.NewDatabaseContext(c)).Update(req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after update resource.
		if err := svc.UpdateAfter(ctx.WithPhase(consts.PHASE_UPDATE_AFTER), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetResponse())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeUpdate,
			Model:     typ.Name(),
			Table:     tableName,
			RecordId:  req.GetID(),
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess, ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess)
		} else {
			// Does not have custom response. the response type is the same as the resource type.
			ResponseJSON(c, CodeSuccess, req)
		}
	}
}

// Patch is a generic function to product gin handler to partial update one resource.
// The resource type depends on the type of interface types.Model.
//
// resource id must be specified.
// - specified in "query parameter `id`".
// - specified in "router parameter `id`".
//
// which one or multiple resources desired modify.
// - specified in "query parameter".
// - specified in "http body data".
func Patch[M types.Model](c *gin.Context) {
	PatchFactory[M]()(c)
}

// PatchFactory is a factory function to product gin handler to partial update one resource.
func PatchFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		var id string
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_PATCH)
		typ := reflect.TypeOf(*new(M)).Elem()
		req := reflect.New(typ).Interface().(M)
		oldVal := reflect.ValueOf(req).Elem()
		newVal := reflect.ValueOf(req).Elem()
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request.
			req := model.NewRequest[M]()
			if err := c.ShouldBindJSON(req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				// Has custom response.
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			id = c.Param(consts.PARAM_ID)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if len(id) == 0 {
				id = req.GetID()
			}
			data := make([]M, 0)
			// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
			// 'typ' is the structure type, such as: model.User.
			// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)

			// Make sure the record must be already exists.
			if err := handler(helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(m).List(&data); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if len(data) != 1 {
				log.Errorz(fmt.Sprintf("the total number of records query from database not equal to 1(%d)", len(data)), zap.String("id", id))
				ResponseJSON(c, CodeNotFound)
				return
			}
			// req.SetCreatedAt(data[0].GetCreatedAt())
			// req.SetCreatedBy(data[0].GetCreatedBy())
			// req.SetUpdatedBy(c.GetString(CTX_USERNAME))
			data[0].SetUpdatedBy(c.GetString(consts.CTX_USERNAME))

			newVal = reflect.ValueOf(req).Elem()
			oldVal = reflect.ValueOf(data[0]).Elem()
			patchValue(log, typ, oldVal, newVal)
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before partial update resource.
		if err := svc.PatchBefore(ctx.WithPhase(consts.PHASE_PATCH_BEFORE), oldVal.Addr().Interface().(M)); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Partial update resource in database.
		if !hasCustomReq {
			if err := handler(helper.NewDatabaseContext(c)).Update(oldVal.Addr().Interface().(M)); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after partial update resource.
		if err := svc.PatchAfter(ctx.WithPhase(consts.PHASE_PATCH_AFTER), oldVal.Addr().Interface().(M)); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		// NOTE: We should record the `req` instead of `oldVal`, the req is `newVal`.
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetResponse())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypePatch,
			Model:     typ.Name(),
			Table:     tableName,
			RecordId:  req.GetID(),
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess, ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess)
		} else {
			// Does not have custom response, the response type is the same as the resource type.
			// NOTE: You should response `oldVal` instead of `req`.
			// The req is `newVal`.
			ResponseJSON(c, CodeSuccess, oldVal.Addr().Interface())
		}
	}
}

// List is a generic function to product gin handler to list resources in backend.
// The resource type deponds on the type of interface types.Model.
//
// If you want make a structure field as query parameter, you should add a "schema"
// tag for it. for example: schema:"name"
//
// TODO:combine query parameter 'page' and 'size' into decoded types.Model
// FIX: retrieve records recursive (current not support in gorm.)
// https://stackoverflow.com/questions/69395891/get-recursive-field-values-in-gorm
// DB.Preload("Category.Category.Category").Find(&Category)
// its works for me.
//
// Query parameters:
//   - All feilds of types.Model's underlying structure but excluding some special fields,
//     such as "password", field value too large, json tag is "-", etc.
//   - `_expand`: strings (multiple items separated by ",").
//     The responsed data to frontend will expanded(retrieve data from external table accoding to foreign key)
//     For examples:
//     /department/myid?_expand=children
//     /department/myid?_expand=children,parent
//   - `_depth`: strings or interger.
//     How depth to retrieve records from datab recursivly, default to 1, value scope is [1,99].
//     For examples:
//     /department/myid?_expand=children&_depth=3
//     /department/myid?_expand=children,parent&_depth=10
//   - `_fuzzy`: bool
//     fuzzy match records in database, default to fase.
//     For examples:
//     /department/myid?_fuzzy=true
func List[M types.Model](c *gin.Context) {
	ListFactory[M]()(c)
}

// ListFactory is a factory function to product gin handler to list resources in backend.
func ListFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_LIST)
		var page, size int
		var startTime, endTime time.Time
		if pageStr, ok := c.GetQuery(consts.QUERY_PAGE); ok {
			page, _ = strconv.Atoi(pageStr)
		}
		if sizeStr, ok := c.GetQuery(consts.QUERY_SIZE); ok {
			size, _ = strconv.Atoi(sizeStr)
		}
		columnName, _ := c.GetQuery(consts.QUERY_COLUMN_NAME)
		index, _ := c.GetQuery(consts.QUERY_INDEX)
		selects, _ := c.GetQuery(consts.QUERY_SELECT)
		if startTimeStr, ok := c.GetQuery(consts.QUERY_START_TIME); ok {
			startTime, _ = time.ParseInLocation(consts.DATE_TIME_LAYOUT, startTimeStr, time.Local)
		}
		if endTimeStr, ok := c.GetQuery(consts.QUERY_END_TIME); ok {
			endTime, _ = time.ParseInLocation(consts.DATE_TIME_LAYOUT, endTimeStr, time.Local)
		}

		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		// 'm' is the structure value, such as: &model.User{ID: myid, Name: myname}.
		typ := reflect.TypeOf(*new(M)).Elem() // the real underlying structure type
		m := reflect.New(typ).Interface().(M)

		// FIXME: failed to convert value when size value is -1.
		if err := schema.NewDecoder().Decode(m, c.Request.URL.Query()); err != nil {
			log.Warn(fmt.Sprintf("failed to decode uri query parameter into model: %s", err))
		}
		log.Infoz(fmt.Sprintf("%s: list query parameter", typ.Name()), zap.Object(typ.String(), m))

		var err error
		var or bool
		var fuzzy bool
		var expands []string
		var cursorNext bool
		var nototal bool // default enable total.
		cursorValue := c.Query(consts.QUERY_CURSOR_VALUE)
		cursorFields := c.Query(consts.QUERY_CURSOR_FIELDS)
		nocache := true // default disable cache.
		depth := 1
		data := make([]M, 0)
		if nocacheStr, ok := c.GetQuery(consts.QUERY_NOCACHE); ok {
			var _nocache bool
			if _nocache, err = strconv.ParseBool(nocacheStr); err == nil {
				nocache = _nocache
			}
		}
		if orStr, ok := c.GetQuery(consts.QUERY_OR); ok {
			or, _ = strconv.ParseBool(orStr)
		}
		if fuzzyStr, ok := c.GetQuery(consts.QUERY_FUZZY); ok {
			fuzzy, _ = strconv.ParseBool(fuzzyStr)
		}
		if cursorNextStr, ok := c.GetQuery(consts.QUERY_CURSOR_NEXT); ok {
			cursorNext, _ = strconv.ParseBool(cursorNextStr)
		}
		if depthStr, ok := c.GetQuery(consts.QUERY_DEPTH); ok {
			depth, _ = strconv.Atoi(depthStr)
			if depth < 1 || depth > 99 {
				depth = 1
			}
		}
		if expandStr, ok := c.GetQuery(consts.QUERY_EXPAND); ok {
			var _expands []string
			items := strings.Split(expandStr, ",")
			if len(items) > 0 {
				if items[0] == consts.VALUE_ALL { // expand all feilds
					items = m.Expands()
				}
			}
			for _, e := range m.Expands() {
				for _, item := range items {
					if strings.EqualFold(item, e) {
						_expands = append(_expands, e)
					}
				}
			}
			// fmt.Println("_expends: ", _expands)
			fieldsMap := make(map[string]reflect.Kind)
			for i := range typ.NumField() {
				fieldsMap[typ.Field(i).Name] = typ.Field(i).Type.Kind()
			}
			for _, e := range _expands {
				// If the expanding field not exists in the structure fiedls, skip depth expand.
				kind, found := fieldsMap[e]
				if !found {
					expands = append(expands, e)
					continue
				}
				// If the expanding field exists in the structure but the kind is not slice, skip depth expand.
				if kind != reflect.Slice {
					expands = append(expands, e)
					continue
				}
				t := make([]string, depth)
				for i := range depth {
					t[i] = e
				}
				// fmt.Println("t: ", t)
				// If expand="Children" and depth=3, the depth expanded is "Children.Children.Children"
				expands = append(expands, strings.Join(t, "."))
			}
			// fmt.Println("expands: ", expands)
		}

		svc := service.Factory[M]().Service()
		svcCtx := helper.NewServiceContext(c)
		// 1.Perform business logic processing before list resources.
		if err = svc.ListBefore(svcCtx.WithPhase(consts.PHASE_LIST_BEFORE), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		sortBy, _ := c.GetQuery(consts.QUERY_SORTBY)
		// 2.List resources from database.
		cache := make([]byte, 0)
		cached := false
		if err = handler(helper.NewDatabaseContext(c)).
			WithScope(page, size).
			WithOr(or).
			WithIndex(index).
			WithSelect(strings.Split(selects, ",")...).
			WithQuery(svc.Filter(svcCtx, m), fuzzy).
			WithQueryRaw(svc.FilterRaw(svcCtx)).
			WithCursor(cursorValue, cursorNext, cursorFields).
			WithExclude(m.Excludes()).
			WithExpand(expands, sortBy).
			WithOrder(sortBy).
			WithTimeRange(columnName, startTime, endTime).
			WithCache(!nocache).
			List(&data, &cache); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		if len(cache) > 0 {
			cached = true
		}
		// 3.Perform business logic processing after list resources.
		if err := svc.ListAfter(svcCtx.WithPhase(consts.PHASE_LIST_AFTER), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		total := new(int64)
		nototalStr, _ := c.GetQuery(consts.QUERY_NOTOTAL)
		nototal, _ = strconv.ParseBool(nototalStr)
		// NOTE: Total count is not provided when using cursor-based pagination.
		if !nototal && len(cursorValue) == 0 {
			if err := handler(helper.NewDatabaseContext(c)).
				// WithScope(page, size). // NOTE: WithScope should not apply in Count method.
				// WithSelect(strings.Split(selects, ",")...). // NOTE: WithSelect should not apply in Count method.
				WithOr(or).
				WithIndex(index).
				WithQuery(svc.Filter(svcCtx, m), fuzzy).
				WithQueryRaw(svc.FilterRaw(svcCtx)).
				WithExclude(m.Excludes()).
				WithTimeRange(columnName, startTime, endTime).
				WithCache(!nocache).
				Count(total); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeList,
			Model:     typ.Name(),
			Table:     tableName,
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		log.Infoz(fmt.Sprintf("%s: length: %d, total: %d", typ.Name(), len(data), *total), zap.Object(typ.Name(), m))
		if cached {
			ResponseBytesList(c, CodeSuccess, uint64(*total), cache)
		} else {
			if !nototal {
				ResponseJSON(c, CodeSuccess, gin.H{
					"items": data,
					"total": *total,
				})
			} else {
				ResponseJSON(c, CodeSuccess, gin.H{
					"items": data,
				})
			}
		}
	}
}

// Get is a generic function to product gin handler to list resource in backend.
// The resource type deponds on the type of interface types.Model.
//
// Query parameters:
//   - `_expand`: strings (multiple items separated by ",").
//     The responsed data to frontend will expanded(retrieve data from external table accoding to foreign key)
//     For examples:
//     /department/myid?_expand=children
//     /department/myid?_expand=children,parent
//   - `_depth`: strings or interger.
//     How depth to retrieve records from datab recursivly, default to 1, value scope is [1,99].
//     For examples:
//     /department/myid?_expand=children&_depth=3
//     /department/myid?_expand=children,parent&_depth=10
//
// Route parameters:
// - id: string or integer.
func Get[M types.Model](c *gin.Context) {
	GetFactory[M]()(c)
}

// GetFactory is a factory function to product gin handler to list resource in backend.
func GetFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_GET)
		if len(c.Param(consts.PARAM_ID)) == 0 {
			log.Error(CodeNotFoundRouteID)
			ResponseJSON(c, CodeNotFoundRouteID)
			return
		}
		index, _ := c.GetQuery(consts.QUERY_INDEX)
		selects, _ := c.GetQuery(consts.QUERY_SELECT)

		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		// 'm' is the structure value, such as: &model.User{ID: myid, Name: myname}.
		typ := reflect.TypeOf(*new(M)).Elem()
		m := reflect.New(typ).Interface().(M)
		m.SetID(c.Param(consts.PARAM_ID)) // `GetBefore` hook need id.

		var err error
		var expands []string
		nocache := true // default disable cache.
		depth := 1
		if nocacheStr, ok := c.GetQuery(consts.QUERY_NOCACHE); ok {
			var _nocache bool
			if _nocache, err = strconv.ParseBool(nocacheStr); err == nil {
				nocache = _nocache
			}
		}
		if depthStr, ok := c.GetQuery(consts.QUERY_DEPTH); ok {
			depth, _ = strconv.Atoi(depthStr)
			if depth < 1 || depth > 99 {
				depth = 1
			}
		}
		if expandStr, ok := c.GetQuery(consts.QUERY_EXPAND); ok {
			var _expands []string
			items := strings.Split(expandStr, ",")
			if len(items) > 0 {
				if items[0] == consts.VALUE_ALL { // expand all feilds
					items = m.Expands()
				}
			}
			for _, e := range m.Expands() {
				for _, item := range items {
					if strings.EqualFold(item, e) {
						_expands = append(_expands, e)
					}
				}
			}
			// fmt.Println("_expends: ", _expands)
			fieldsMap := make(map[string]reflect.Kind)
			for i := range typ.NumField() {
				fieldsMap[typ.Field(i).Name] = typ.Field(i).Type.Kind()
			}
			for _, e := range _expands {
				// If the expanding field not exists in the structure fiedls, skip depth expand.
				// TODO: if the field type is the structure name, make depth expand.
				kind, found := fieldsMap[e]
				if !found {
					expands = append(expands, e)
					continue
				}
				// If the expanding field exists in the structure but the kind is not slice, skip depth expand.
				if kind != reflect.Slice {
					expands = append(expands, e)
					continue
				}
				t := make([]string, depth)
				for i := range depth {
					t[i] = e
				}
				// If expand="Children" and depth=3, the depth expanded is "Children.Children.Children"
				expands = append(expands, strings.Join(t, "."))
			}
			// fmt.Println("expands: ", expands)
		}
		log.Infoz("", zap.Object(typ.Name(), m))

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before get resource.
		if err = svc.GetBefore(helper.NewServiceContext(c).WithPhase(consts.PHASE_GET_BEFORE), m); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Get resource from database.
		cache := make([]byte, 0)
		cached := false
		if err = handler(helper.NewDatabaseContext(c)).
			WithIndex(index).
			WithSelect(strings.Split(selects, ",")...).
			WithExpand(expands).
			WithCache(!nocache).
			Get(m, c.Param(consts.PARAM_ID), &cache); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		if len(cache) > 0 {
			cached = true
		}
		// 3.Perform business logic processing after get resource.
		if err := svc.GetAfter(helper.NewServiceContext(c).WithPhase(consts.PHASE_GET_AFTER), m); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// It will returns a empty types.Model if found nothing from database,
		// we should response status code "CodeNotFound".
		if !cached {
			if len(m.GetID()) == 0 || (m.GetCreatedAt().Equal(time.Time{})) {
				log.Error(CodeNotFound)
				ResponseJSON(c, CodeNotFound)
				return
			}
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeGet,
			Model:     typ.Name(),
			Table:     tableName,
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})
		if cached {
			ResponseBytes(c, CodeSuccess, cache)
		} else {
			ResponseJSON(c, CodeSuccess, m)
		}
	}
}

// BatchCreate
// Example:
/*
Request Method: POST
Request URL: /api/v1/users/batch

Request Data:
{
  "items": [
    {
      "username": "johndoe",
      "email": "john.doe@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "department": "Engineering"
    },
    {
      "username": "janedoe",
      "email": "jane.doe@example.com",
      "firstName": "Jane",
      "lastName": "Doe",
      "department": "Marketing"
    },
    {
      "username": "bobsmith",
      "email": "bob.smith@example.com",
      "firstName": "Bob",
      "lastName": "Smith",
      "department": "Finance"
    }
  ],
  "options": {
    "continueOnError": true
  }
}

Response Data:
{
  "items": [
    {
      "id": "user-123",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "department": "Engineering",
      "createdAt": "2025-03-06T10:15:30Z"
    },
    {
      "status": "error",
      "error": {
        "code": 400,
        "message": "Email already in use"
      },
      "request": {
        "username": "janedoe",
        "email": "jane.doe@example.com",
        "firstName": "Jane",
        "lastName": "Doe",
        "department": "Marketing"
      }
    },
    {
      "id": "user-125",
      "username": "bobsmith",
      "email": "bob.smith@example.com",
      "firstName": "Bob",
      "lastName": "Smith",
      "department": "Finance",
      "createdAt": "2025-03-06T10:15:30Z"
    }
  ],
  "summary": {
    "total": 3,
    "succeeded": 2,
    "failed": 1
  }
}
}
*/

type requestData[M types.Model] struct {
	// Ids is the id list that should be batch delete.
	Ids []string `json:"ids,omitempty"`
	// Items is the resource list that should be batch create/update/partial update.
	Items []M `json:"items,omitempty"`
	// Options is the batch operation options.
	Options *options `json:"options,omitempty"`
	// Summary is the batch operation result summary.
	Summary *summary `json:"summary,omitempty"`
}

type options struct {
	Atomic bool `json:"atomic,omitempty"`
	Purge  bool `json:"purge,omitempty"`
}

type summary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

func CreateMany[M types.Model](c *gin.Context) {
	CreateManyFactory[M]()(c)
}

// CreateManyFactory
func CreateManyFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		var err error
		var reqErr error
		var req requestData[M]
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_CREATE_MANY)
		typ := reflect.TypeOf(*new(M)).Elem()
		val := reflect.New(typ).Interface().(M)
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request.
			req := model.NewRequest[M]()
			if reqErr = c.ShouldBindJSON(req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				// Has custom response.
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}

			if req.Options == nil {
				req.Options = new(options)
			}
			for _, m := range req.Items {
				m.SetCreatedBy(c.GetString(consts.CTX_USERNAME))
				m.SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
				log.Infoz("create_many", zap.Bool("atomic", req.Options.Atomic), zap.Object(typ.Name(), m))
			}
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before batch create resource.
		if err = svc.CreateManyBefore(ctx.WithPhase(consts.PHASE_CREATE_MANY_BEFORE), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 2.Batch create resource in database.
		if reqErr != io.EOF && !hasCustomReq {
			if err = handler(helper.NewDatabaseContext(c)).WithExpand(val.Expands()).Create(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after batch create resource
		if err = svc.CreateManyAfter(ctx.WithPhase(consts.PHASE_CREATE_MANY_AFTER), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetResponse())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeCreateMany,
			Model:     typ.Name(),
			Table:     tableName,
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated), ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated))
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			// FIXME: 如果某些字段增加了 gorm unique tag, 则更新成功后的资源 ID 时随机生成的，并不是数据库中的
			if reqErr != io.EOF {
				req.Summary = &summary{
					Total:     len(req.Items),
					Succeeded: len(req.Items),
					Failed:    0,
				}
			}
			ResponseJSON(c, CodeSuccess.WithStatus(http.StatusCreated), req)
		}
	}
}

// DeleteMany
func DeleteMany[M types.Model](c *gin.Context) {
	DeleteManyFactory[M]()(c)
}

// DeleteManyFactory
func DeleteManyFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_DELETE_MANY)
		log.Info("delete_many")

		var err error
		var reqErr error
		var req requestData[M]
		if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
			log.Error(reqErr)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		if reqErr == io.EOF {
			log.Warn("empty request body")
		}

		// 1.Perform business logic processing before batch delete resources.
		svc := service.Factory[M]().Service()
		typ := reflect.TypeOf(*new(M)).Elem()
		req.Items = make([]M, 0, len(req.Ids))
		for _, id := range req.Ids {
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)
			req.Items = append(req.Items, m)
		}
		if err = svc.DeleteManyBefore(helper.NewServiceContext(c).WithPhase(consts.PHASE_DELETE_MANY_BEFORE), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		if req.Options == nil {
			req.Options = new(options)
		}
		// 2.Batch delete resources in database.
		if reqErr != io.EOF {
			if err = handler(helper.NewDatabaseContext(c)).WithPurge(req.Options.Purge).Delete(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after batch delete resources.
		if err = svc.DeleteManyAfter(helper.NewServiceContext(c).WithPhase(consts.PHASE_DELETE_MANY_AFTER), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		record, _ := json.Marshal(req)
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeDeleteMany,
			Model:     typ.Name(),
			Table:     tableName,
			Record:    util.BytesToString(record),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if reqErr != io.EOF {
			req.Summary = &summary{
				Total:     len(req.Items),
				Succeeded: len(req.Items),
				Failed:    0,
			}
		}
		req.Ids = nil
		req.Items = nil
		req.Options = nil
		// not response req.
		ResponseJSON(c, CodeSuccess.WithStatus(http.StatusNoContent))
	}
}

// UpdateMany
func UpdateMany[M types.Model](c *gin.Context) {
	UpdateManyFactory[M]()(c)
}

// UpdateManyFactory
func UpdateManyFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		var err error
		var reqErr error
		var req requestData[M]
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_UPDATE_MANY)
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request
			req := model.NewRequest[M]()
			if reqErr = c.ShouldBindJSON(req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeInvalidParam.WithErr(reqErr))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				// Has custom response
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request, the request type is the same as the resource type.
			if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before batch update resource.
		if err = svc.UpdateManyBefore(ctx.WithPhase(consts.PHASE_UPDATE_MANY_BEFORE), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Batch update resource in database.
		if reqErr != io.EOF && !hasCustomReq {
			if err = handler(helper.NewDatabaseContext(c)).Update(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after batch update resource.
		if err = svc.UpdateManyAfter(ctx.WithPhase(consts.PHASE_UPDATE_MANY_AFTER), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		typ := reflect.TypeOf(*new(M)).Elem()
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetResponse())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypeUpdateMany,
			Model:     typ.Name(),
			Table:     tableName,
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess, ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess)
		} else {
			// Does not have custom response, the response type is the same as the request type.
			if reqErr != io.EOF {
				req.Summary = &summary{
					Total:     len(req.Items),
					Succeeded: len(req.Items),
					Failed:    0,
				}
			}
			ResponseJSON(c, CodeSuccess, req)
		}
	}
}

// PatchMany
func PatchMany[M types.Model](c *gin.Context) {
	PatchManyFactory[M]()(c)
}

// PatchManyFactory
func PatchManyFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		var err error
		var reqErr error
		var req requestData[M]
		var shouldUpdates []M
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_PATCH_MANY)
		typ := reflect.TypeOf(*new(M)).Elem()
		ctx := helper.NewServiceContext(c)
		hasCustomReq := model.HasRequest[M]()
		hasCustomResp := model.HasResponse[M]()

		if hasCustomReq {
			// Has custom request
			req := model.NewRequest[M]()
			if reqErr = c.ShouldBindJSON(req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeInvalidParam.WithErr(reqErr))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
			ctx.SetRequest(req)
			if hasCustomResp {
				// Has custom response
				ctx.SetResponse(model.NewResponse[M]())
			}
		} else {
			// Does not have custom request
			if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
				log.Error(reqErr)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
			if reqErr == io.EOF {
				log.Warn("empty request body")
			}
			for _, m := range req.Items {
				var results []M
				v := reflect.New(typ).Interface().(M)
				v.SetID(m.GetID())
				if err = handler(helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(v).List(&results); err != nil {
					log.Error(err)
					continue
				}
				if len(results) != 1 {
					log.Warnf(fmt.Sprintf("partial update resource not found, expect 1 but got: %d", len(results)))
					continue
				}
				if len(results[0].GetID()) == 0 {
					log.Warnf("partial update resource not found, id is empty")
					continue
				}
				oldVal, newVal := reflect.ValueOf(results[0]).Elem(), reflect.ValueOf(m).Elem()
				patchValue(log, typ, oldVal, newVal)
				shouldUpdates = append(shouldUpdates, oldVal.Addr().Interface().(M))
			}
		}

		svc := service.Factory[M]().Service()
		// 1.Perform business logic processing before batch patch resource.
		if err = svc.PatchManyBefore(ctx.WithPhase(consts.PHASE_PATCH_MANY_BEFORE), shouldUpdates...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 2.Batch partial update resource in database.
		if reqErr != io.EOF && !hasCustomReq {
			if err = handler(helper.NewDatabaseContext(c)).Update(shouldUpdates...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure.WithErr(err))
				return
			}
		}
		// 3.Perform business logic processing after batch patch resource.
		if err := svc.PatchManyAfter(ctx.WithPhase(consts.PHASE_PATCH_MANY_AFTER), shouldUpdates...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// 4.record operation log to database.
		var tableName string
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
		// NOTE: We should record the `req` instead of `oldVal`, the req is `newVal`.
		record, _ := json.Marshal(req)
		reqData, _ := json.Marshal(ctx.GetRequest())
		respData, _ := json.Marshal(ctx.GetResponse())
		cb.Enqueue(&model_log.OperationLog{
			Op:        model_log.OperationTypePatchMany,
			Model:     typ.Name(),
			Table:     tableName,
			Record:    util.BytesToString(record),
			Request:   util.BytesToString(reqData),
			Response:  util.BytesToString(respData),
			IP:        c.ClientIP(),
			User:      c.GetString(consts.CTX_USERNAME),
			RequestId: c.GetString(consts.REQUEST_ID),
			URI:       c.Request.RequestURI,
			Method:    c.Request.Method,
			UserAgent: c.Request.UserAgent(),
		})

		if hasCustomResp {
			// Has custom response.
			ResponseJSON(c, CodeSuccess, ctx.GetResponse())
		} else if hasCustomReq {
			// Has custom request but not custom response.
			ResponseJSON(c, CodeSuccess)
		} else {
			// Does not have custom response, the response type is the same as the request type.
			if reqErr != io.EOF {
				req.Summary = &summary{
					Total:     len(req.Items),
					Succeeded: len(req.Items),
					Failed:    0,
				}
			}
			ResponseJSON(c, CodeSuccess, req)
		}
	}
}

// Export is a generic function to product gin handler to export resources to frontend.
// The resource type deponds on the type of interface types.Model.
//
// If you want make a structure field as query parameter, you should add a "schema"
// tag for it. for example: schema:"name"
//
// TODO:combine query parameter 'page' and 'size' into decoded types.Model
// FIX: retrieve records recursive (current not support in gorm.)
// https://stackoverflow.com/questions/69395891/get-recursive-field-values-in-gorm
// DB.Preload("Category.Category.Category").Find(&Category)
// its works for me.
//
// Query parameters:
//   - All feilds of types.Model's underlying structure but excluding some special fields,
//     such as "password", field value too large, json tag is "-", etc.
//   - `_expand`: strings (multiple items separated by ",").
//     The responsed data to frontend will expanded(retrieve data from external table accoding to foreign key)
//     For examples:
//     /department/myid?_expand=children
//     /department/myid?_expand=children,parent
//   - `_depth`: strings or interger.
//     How depth to retrieve records from datab recursivly, default to 1, value scope is [1,99].
//     For examples:
//     /department/myid?_expand=children&_depth=3
//     /department/myid?_expand=children,parent&_depth=10
//   - `_fuzzy`: bool
//     fuzzy match records in database, default to fase.
//     For examples:
//     /department/myid?_fuzzy=true
func Export[M types.Model](c *gin.Context) {
	ExportFactory[M]()(c)
}

// ExportFactory is a factory function to export resources to frontend.
func ExportFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_EXPORT)
		var page, size, limit int
		var startTime, endTime time.Time
		if pageStr, ok := c.GetQuery(consts.QUERY_PAGE); ok {
			page, _ = strconv.Atoi(pageStr)
		}
		if sizeStr, ok := c.GetQuery(consts.QUERY_SIZE); ok {
			size, _ = strconv.Atoi(sizeStr)
		}
		if limitStr, ok := c.GetQuery(consts.QUERY_LIMIT); ok {
			limit, _ = strconv.Atoi(limitStr)
		}
		columnName, _ := c.GetQuery(consts.QUERY_COLUMN_NAME)
		index, _ := c.GetQuery(consts.QUERY_INDEX)
		selects, _ := c.GetQuery(consts.QUERY_SELECT)
		if startTimeStr, ok := c.GetQuery(consts.QUERY_START_TIME); ok {
			startTime, _ = time.ParseInLocation(consts.DATE_TIME_LAYOUT, startTimeStr, time.Local)
		}
		if endTimeStr, ok := c.GetQuery(consts.QUERY_END_TIME); ok {
			endTime, _ = time.ParseInLocation(consts.DATE_TIME_LAYOUT, endTimeStr, time.Local)
		}

		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		// 'm' is the structure value, such as: &model.User{ID: myid, Name: myname}.
		typ := reflect.TypeOf(*new(M)).Elem() // the real underlying structure type
		m := reflect.New(typ).Interface().(M)

		if err := schema.NewDecoder().Decode(m, c.Request.URL.Query()); err != nil {
			log.Warn("failed to parse uri query parameter into model: ", err)
		}
		log.Info("query parameter: ", m)

		var err error
		var or bool
		var fuzzy bool
		var depth int = 1
		var expands []string
		data := make([]M, 0)
		if orStr, ok := c.GetQuery(consts.QUERY_OR); ok {
			or, _ = strconv.ParseBool(orStr)
		}
		if fuzzyStr, ok := c.GetQuery(consts.QUERY_FUZZY); ok {
			fuzzy, _ = strconv.ParseBool(fuzzyStr)
		}
		if depthStr, ok := c.GetQuery(consts.QUERY_DEPTH); ok {
			depth, _ = strconv.Atoi(depthStr)
			if depth < 1 || depth > 99 {
				depth = 1
			}
		}
		if expandStr, ok := c.GetQuery(consts.QUERY_EXPAND); ok {
			var _expands []string
			items := strings.Split(expandStr, ",")
			if len(items) > 0 {
				if items[0] == consts.VALUE_ALL { // expand all feilds
					items = m.Expands()
				}
			}
			for _, e := range m.Expands() {
				for _, item := range items {
					if strings.EqualFold(item, e) {
						_expands = append(_expands, e)
					}
				}
			}
			// fmt.Println("_expends: ", _expands)
			fieldsMap := make(map[string]reflect.Kind)
			for i := range typ.NumField() {
				fieldsMap[typ.Field(i).Name] = typ.Field(i).Type.Kind()
			}
			for _, e := range _expands {
				// If the expanding field not exists in the structure fiedls, skip depth expand.
				kind, found := fieldsMap[e]
				if !found {
					expands = append(expands, e)
					continue
				}
				// If the expanding field exists in the structure but the kind is not slice, skip depth expand.
				if kind != reflect.Slice {
					expands = append(expands, e)
					continue
				}
				t := make([]string, depth)
				for i := range depth {
					t[i] = e
				}
				// fmt.Println("t: ", t)
				// If expand="Children" and depth=3, the depth expanded is "Children.Children.Children"
				expands = append(expands, strings.Join(t, "."))
			}
			// fmt.Println("expands: ", expands)
		}

		svc := service.Factory[M]().Service()
		svcCtx := helper.NewServiceContext(c)
		// 1.Perform business logic processing before list resources.
		if err = svc.ListBefore(svcCtx.WithPhase(consts.PHASE_EXPORT), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		sortBy, _ := c.GetQuery(consts.QUERY_SORTBY)
		_, _ = page, size
		// 2.List resources from database.
		if err = handler(helper.NewDatabaseContext(c)).
			// WithScope(page, size). // 不要使用 WithScope, 否则 WithLimit 不生效
			WithLimit(limit).
			WithOr(or).
			WithIndex(index).
			WithSelect(strings.Split(selects, ",")...).
			WithQuery(svc.Filter(svcCtx, m), fuzzy).
			WithQueryRaw(svc.FilterRaw(svcCtx)).
			WithExclude(m.Excludes()).
			WithExpand(expands, sortBy).
			WithOrder(sortBy).
			WithTimeRange(columnName, startTime, endTime).
			List(&data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// 3.Perform business logic processing after list resources.
		if err = svc.ListAfter(svcCtx.WithPhase(consts.PHASE_EXPORT), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		log.Info("export data length: ", len(data))
		// 4.Export
		exported, err := svc.Export(svcCtx.WithPhase(consts.PHASE_EXPORT), data...)
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// // 5.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// record, _ := json.Marshal(data)
		// if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 	Op:        model.OperationTypeExport,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	Record:    util.BytesToString(record),
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error("failed to write operation log to database: ", err.Error())
		// }
		ResponseDATA(c, exported, map[string]string{
			"Content-Disposition": "attachment; filename=exported.xlsx",
		})
	}
}

// Import
func Import[M types.Model](c *gin.Context) {
	ImportFactory[M]()(c)
}

// ImportFactory
func ImportFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_IMPORT)
		// NOTE:字段为 file 必须和前端协商好.
		file, err := c.FormFile("file")
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// check file size.
		if file.Size > int64(MAX_IMPORT_SIZE) {
			log.Error(CodeTooLargeFile)
			ResponseJSON(c, CodeTooLargeFile)
			return
		}
		fd, err := file.Open()
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		defer fd.Close()

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, fd); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// filetype must be png or jpg.
		filetype, mime := filetype.DetectBytes(buf.Bytes())
		_, _ = filetype, mime

		// check filetype

		ml, err := service.Factory[M]().Service().Import(helper.NewServiceContext(c).WithPhase(consts.PHASE_IMPORT), buf)
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}

		// service layer already create/update the records in database, just update fields "created_by", "updated_by".
		for i := range ml {
			ml[i].SetCreatedBy(c.GetString(consts.CTX_USERNAME))
			ml[i].SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
		}
		if err := handler(helper.NewDatabaseContext(c)).Update(ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure.WithErr(err))
			return
		}
		// // record operation log to database.
		// typ := reflect.TypeOf(*new(M)).Elem()
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// record, _ := json.Marshal(ml)
		// if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 	Op:        model.OperationTypeImport,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	Record:    util.BytesToString(record),
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error("failed to write operation log to database: ", err.Error())
		// }
		ResponseJSON(c, CodeSuccess)
	}
}

func extractConfig[M types.Model](cfg ...*types.ControllerConfig[M]) (handler func(ctx *types.DatabaseContext) types.Database[M], db any) {
	if len(cfg) > 0 {
		if cfg[0] != nil {
			db = cfg[0].DB
		}
	}
	handler = func(ctx *types.DatabaseContext) types.Database[M] {
		fn := database.Database[M](ctx)
		if len(cfg) > 0 {
			if cfg[0] != nil {
				if len(cfg[0].TableName) > 0 {
					fn = database.Database[M](ctx).WithDB(cfg[0].DB).WithTable(cfg[0].TableName)
				} else {
					fn = database.Database[M](ctx).WithDB(cfg[0].DB)
				}
			}
		}
		return fn
	}
	return
}
