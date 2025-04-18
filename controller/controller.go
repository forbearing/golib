package controller

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/pkg/filetype"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/types/helper"
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

var pluralizeCli = pluralize.NewClient()

// Create is a generic function to product gin handler to create one resource.
// The resource type depends on the type of interface types.Model.
func Create[M types.Model](c *gin.Context) {
	CreateFactory[M]()(c)
}

// CreateFactory is a factory function to product gin handler to create one resource.
func CreateFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_CREATE)
		req := *new(M)
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeInvalidParam)
			return
		}

		req.SetCreatedBy(c.GetString(consts.CTX_USERNAME))
		req.SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
		log.Infoz("create", zap.Object(reflect.TypeOf(*new(M)).Elem().String(), req))

		// 1.Perform business logic processing before create resource.
		if err := new(service.Factory[M]).Service().CreateBefore(helper.NewServiceContext(c), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Create resource in database.
		// database.Database().Delete just set "deleted_at" field to current time, not really delete.
		// We should update it instead of creating it, and update the "created_at" and "updated_at" field.
		// NOTE: WithExpand(req.Expands()...) is not a good choices.
		// if err := database.Database[M]().WithExpand(req.Expands()...).Update(req); err != nil {
		if err := handler(helper.NewDatabaseContext(c)).WithExpand(req.Expands()).Create(req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 3.Perform business logic processing after create resource
		if err := new(service.Factory[M]).Service().CreateAfter(helper.NewServiceContext(c), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// // 4.record operation log to database.
		// typ := reflect.TypeOf(*new(M)).Elem()
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// record, _ := json.Marshal(req)
		// // TODO: should we record the operation in the database of `db` instance.
		// if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 	Op:        model.OperationTypeCreate,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	RecordId:  req.GetID(),
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

		ResponseJSON(c, CodeSuccess, req)
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

		// Delete one record accoding to "query parameter `id`".
		if id, ok := c.GetQuery(consts.QUERY_ID); ok {
			// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)
			ml = append(ml, m)
		}
		// Delete one record accoding to "route parameter `id`".
		if id := c.Param(consts.PARAM_ID); len(id) != 0 {
			// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
			m := reflect.New(typ).Interface().(M)
			m.SetID(id)
			ml = append(ml, m)
		}
		// Delete multiple records accoding to "http body data".
		ids := make([]string, 0)
		if err := c.ShouldBindJSON(&ids); err == nil {
			// remove empty string
			_ids := make([]string, 0)
			for i := range ids {
				if len(ids[i]) == 0 {
					continue
				}
				_ids = append(_ids, ids[i])
			}
			if len(_ids) == 0 {
				log.Warn("id list is empty, skip delete")
				ResponseJSON(c, CodeSuccess)
				return
			}
			for i := range _ids {
				// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
				m := reflect.New(typ).Interface().(M)
				m.SetID(_ids[i])
				ml = append(ml, m)
			}
		} else if err == io.EOF {
			log.Warnf("empty request body")
		} else {
			log.Warn(err)
		}

		log.Info(fmt.Sprintf("%s delete %v", typ.Name(), ids))
		// 1.Perform business logic processing before delete resources.
		if err := new(service.Factory[M]).Service().DeleteBefore(helper.NewServiceContext(c), ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// // find out the records and record to operation log.
		// copied := make([]M, len(ml))
		// for i := range ml {
		// 	m := reflect.New(typ).Interface().(M)
		// 	m.SetID(ml[i].GetID())
		// 	if err := handler(helper.NewDatabaseContext(c)).WithExpand(m.Expands()).Get(m, ml[i].GetID()); err != nil {
		// 		log.Error(err)
		// 	}
		// 	copied[i] = m
		// }

		// 2.Delete resources in database.
		if err := handler(helper.NewDatabaseContext(c)).Delete(ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 3.Perform business logic processing after delete resources.
		if err := new(service.Factory[M]).Service().DeleteAfter(helper.NewServiceContext(c), ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// // 4.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// for i := range ml {
		// 	record, _ := json.Marshal(copied[i])
		// 	if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 		Op:        model.OperationTypeDelete,
		// 		Model:     typ.Name(),
		// 		Table:     tableName,
		// 		RecordId:  ml[i].GetID(),
		// 		Record:    util.BytesToString(record),
		// 		IP:        c.ClientIP(),
		// 		User:      c.GetString(consts.CTX_USERNAME),
		// 		RequestId: c.GetString(consts.REQUEST_ID),
		// 		URI:       c.Request.RequestURI,
		// 		Method:    c.Request.Method,
		// 		UserAgent: c.Request.UserAgent(),
		// 	}); err != nil {
		// 		log.Error("failed to write operation log to database: ", err.Error())
		// 	}
		// }

		ResponseJSON(c, CodeSuccess)
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
		id := c.Param(consts.PARAM_ID)
		req := *new(M)
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		log.Infoz("update from request", zap.Object(reflect.TypeOf(*new(M)).Elem().String(), req))
		if len(id) == 0 {
			id = req.GetID()
		}
		data := make([]M, 0)
		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
		typ := reflect.TypeOf(*new(M)).Elem()
		m := reflect.New(typ).Interface().(M)
		m.SetID(id)

		// Make sure the record must be already exists.
		if err := handler(helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(m).List(&data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		if len(data) != 1 {
			log.Errorz(fmt.Sprintf("the total number of records query from database not equal to 1(%d)", len(data)), zap.String("id", id))
			ResponseJSON(c, CodeNotFound)
			return
		}

		req.SetCreatedAt(data[0].GetCreatedAt())           // keep original "created_at"
		req.SetCreatedBy(data[0].GetCreatedBy())           // keep original "created_by"
		req.SetUpdatedBy(c.GetString(consts.CTX_USERNAME)) // keep original "updated_by"
		// 1.Perform business logic processing before update resource.
		if err := new(service.Factory[M]).Service().UpdateBefore(helper.NewServiceContext(c), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Update resource in database.
		log.Infoz("update in database", zap.Object(typ.Name(), req))
		if err := handler(helper.NewDatabaseContext(c)).Update(req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 3.Perform business logic processing after update resource.
		if err := new(service.Factory[M]).Service().UpdateAfter(helper.NewServiceContext(c), req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// // 4.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// record, _ := json.Marshal(req)
		// if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 	Op:        model.OperationTypeUpdate,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	RecordId:  req.GetID(),
		// 	Record:    util.BytesToString(record),
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error(fmt.Sprintf("failed to write operation log to database: %s", err.Error()))
		// }

		ResponseJSON(c, CodeSuccess, req)
	}
}

// UpdatePartial is a generic function to product gin handler to partial update one resource.
// The resource type depends on the type of interface types.Model.
//
// resource id must be specified.
// - specified in "query parameter `id`".
// - specified in "router parameter `id`".
//
// which one or multiple resources desired modify.
// - specified in "query parameter".
// - specified in "http body data".
func UpdatePartial[M types.Model](c *gin.Context) {
	UpdatePartialFactory[M]()(c)
}

// UpdatePartialFactory is a factory function to product gin handler to partial update one resource.
func UpdatePartialFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_UPDATE_PARTIAL)
		id := c.Param(consts.PARAM_ID)
		req := *new(M)
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		if len(id) == 0 {
			id = req.GetID()
		}
		data := make([]M, 0)
		// The underlying type of interface types.Model must be pointer to structure, such as *model.User.
		// 'typ' is the structure type, such as: model.User.
		// 'm' is the structure value such as: &model.User{ID: myid, Name: myname}.
		typ := reflect.TypeOf(*new(M)).Elem()
		m := reflect.New(typ).Interface().(M)
		m.SetID(id)

		// Make sure the record must be already exists.
		if err := handler(helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(m).List(&data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
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

		newVal := reflect.ValueOf(req).Elem()
		oldVal := reflect.ValueOf(data[0]).Elem()
		for i := range typ.NumField() {
			// fmt.Println(typ.Field(i).Name, typ.Field(i).Type, typ.Field(i).Type.Kind(), newVal.Field(i).IsValid(), newVal.Field(i).CanSet())
			switch typ.Field(i).Type.Kind() {
			case reflect.Struct: // skip update base model.
				switch typ.Field(i).Type.Name() {
				case "GormTime": // The underlying type of model.GormTime(type of time.Time) is struct, we should continue handle.

				case "Asset", "Base":
					// AssetChecking 有匿名继承 Asset, 所以要检查是不是结构体 Asset.
					//
					// 可以自动深度查找,不需要链式查找, 例如
					// newVal.FieldByName("Asset").FieldByName("Remark").IsValid() 可以简化为
					// newVal.FieldByName("Remark").IsValid()

					// Make sure the type of "Remark" is pointer to golang base type.
					fieldRemark := "Remark"
					if newVal.FieldByName(fieldRemark).IsValid() { // WARN: oldVal.FieldByName(fieldRemark) maybe <invalid reflect.Value>
						if !newVal.FieldByName(fieldRemark).IsZero() {
							if newVal.FieldByName(fieldRemark).CanSet() {
								// output log must before set value.
								if newVal.FieldByName(fieldRemark).Kind() == reflect.Pointer {
									log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", fieldRemark, typ.Name(),
										oldVal.FieldByName(fieldRemark).Elem(), newVal.FieldByName(fieldRemark).Elem())) // WARN: you shouldn't call oldVal.FieldByName(fieldRemark).Elem().Interface()
								} else {
									log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", fieldRemark, typ.Name(),
										oldVal.FieldByName(fieldRemark).Interface(), newVal.FieldByName(fieldRemark).Interface()))
								}
								oldVal.FieldByName(fieldRemark).Set(newVal.FieldByName(fieldRemark)) // set old value by new value
							}
						}
					}
					// Make sure the type of "Order" is pointer to golang base type.
					fieldOrder := "Order"
					if newVal.FieldByName(fieldOrder).IsValid() { // WARN: oldVal.FieldByName(fieldOrder) maybe <invalid reflect.Value>
						if !newVal.FieldByName(fieldOrder).IsZero() {
							if newVal.FieldByName(fieldOrder).CanSet() {
								// output log must before set value.
								if newVal.FieldByName(fieldOrder).Kind() == reflect.Pointer {
									log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", fieldOrder, typ.Name(),
										oldVal.FieldByName(fieldOrder).Elem(), newVal.FieldByName(fieldOrder).Elem())) // WARN: you shouldn't call oldVal.FieldByName(fieldOrder).Elem().Interface()
								} else {
									log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", fieldOrder, typ.Name(),
										oldVal.FieldByName(fieldOrder).Interface(), newVal.FieldByName(fieldOrder).Interface()))
								}
								oldVal.FieldByName(fieldOrder).Set(newVal.FieldByName(fieldOrder)) // set old value by new value.
							}
						}
					}
					continue

				default:
					continue
				}
			}
			if !newVal.Field(i).IsValid() {
				// log.Warnf("field %s is invalid, skip", typ.Field(i).Name)
				continue
			}
			// base type such like int and string have default value(zero value).
			// If the struct field(the field type is golang base type) supported by patch update,
			// the field type must be pointer to base type, such like *string, *int.
			if newVal.Field(i).IsZero() {
				// log.Warnf("field %s is zero value, skip", typ.Field(i).Name)
				// log.Warnf("DeepEqual: %v : %v : %v : %v", typ.Field(i).Name, newVal.Field(i).Interface(), oldVal.Field(i).Interface(), reflect.DeepEqual(newVal.Field(i), oldVal.Field(i)))
				continue
			}
			if newVal.Field(i).CanSet() {
				// output log must before set value.
				if newVal.Field(i).Kind() == reflect.Pointer {
					log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", typ.Name(), typ.Field(i).Name, oldVal.Field(i).Elem().Interface(), newVal.Field(i).Elem().Interface()))
				} else {
					log.Info(fmt.Sprintf("[UpdatePartial %s] field: %s: %v --> %v", typ.Name(), typ.Field(i).Name, oldVal.Field(i).Interface(), newVal.Field(i).Interface()))
				}
				oldVal.Field(i).Set(newVal.Field(i)) // set old value by new value
			}
		}
		// zap.L().Info("[UpdatePartial]", zap.Object(typ.String(), oldVal.Addr().Interface().(M)))

		// 1.Perform business logic processing before partial update resource.
		if err := new(service.Factory[M]).Service().UpdatePartialBefore(helper.NewServiceContext(c), oldVal.Addr().Interface().(M)); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Partial update resource in database.
		if err := handler(helper.NewDatabaseContext(c)).Update(oldVal.Addr().Interface().(M)); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 3.Perform business logic processing after partial update resource.
		if err := new(service.Factory[M]).Service().UpdatePartialAfter(helper.NewServiceContext(c), oldVal.Addr().Interface().(M)); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// // 4.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// // NOTE: We should record the `req` instead of `oldVal`,
		// // The req is `newVal`.
		// record, _ := json.Marshal(req)
		// if err := database.Database[*model.OperationLog]().WithDB(db).Create(&model.OperationLog{
		// 	Op:        model.OperationTypeUpdatePartial,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	RecordId:  req.GetID(),
		// 	Record:    util.BytesToString(record),
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error(fmt.Sprintf("failed to write operation log to database: %s", err.Error()))
		// }

		// NOTE: You should response `oldVal` instead of `req`.
		// The req is `newVal`.
		ResponseJSON(c, CodeSuccess, oldVal.Addr().Interface())
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
		var nototal bool // default enable total.
		nocache := true  // default disable cache.
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

		svc := new(service.Factory[M]).Service()
		svcCtx := helper.NewServiceContext(c)
		// 1.Perform business logic processing before list resources.
		if err = svc.ListBefore(svcCtx, &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
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
			WithExclude(m.Excludes()).
			WithExpand(expands, sortBy).
			WithOrder(sortBy).
			WithTimeRange(columnName, startTime, endTime).
			WithCache(!nocache).
			List(&data, &cache); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		if len(cache) > 0 {
			cached = true
		}
		// 3.Perform business logic processing after list resources.
		if err := svc.ListAfter(svcCtx, &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		total := new(int64)
		nototalStr, _ := c.GetQuery(consts.QUERY_NOTOTAL)
		nototal, _ = strconv.ParseBool(nototalStr)
		if !nototal {
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
				ResponseJSON(c, CodeFailure)
				return
			}
		}
		// // 4.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// if err := database.Database[*model.OperationLog]().Create(&model.OperationLog{
		// 	Op:        model.OperationTypeList,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error("failed to write operation log to database: ", err.Error())
		// }
		// types.Sort[M](types.UpdatedTime, data)
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

		// 1.Perform business logic processing before get resource.
		if err = new(service.Factory[M]).Service().GetBefore(helper.NewServiceContext(c), m); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
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
			ResponseJSON(c, CodeFailure)
			return
		}
		if len(cache) > 0 {
			cached = true
		}
		// 3.Perform business logic processing after get resource.
		if err := new(service.Factory[M]).Service().GetAfter(helper.NewServiceContext(c), m); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
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

		// // 4.record operation log to database.
		// var tableName string
		// items := strings.Split(typ.Name(), ".")
		// if len(items) > 0 {
		// 	tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		// }
		// if err := database.Database[*model.OperationLog]().Create(&model.OperationLog{
		// 	Op:        model.OperationTypeGet,
		// 	Model:     typ.Name(),
		// 	Table:     tableName,
		// 	IP:        c.ClientIP(),
		// 	User:      c.GetString(consts.CTX_USERNAME),
		// 	RequestId: c.GetString(consts.REQUEST_ID),
		// 	URI:       c.Request.RequestURI,
		// 	Method:    c.Request.Method,
		// 	UserAgent: c.Request.UserAgent(),
		// }); err != nil {
		// 	log.Error("failed to write operation log to database: ", err.Error())
		// }
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
	Items   []M      `json:"items,omitempty"`
	Options *options `json:"options,omitempty"`
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

func BatchCreate[M types.Model](c *gin.Context) {
	BatchCreateFactory[M]()(c)
}

// BatchCreateFactory
func BatchCreateFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	_ = handler
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_BATCH_CREATE)
		var err error
		var reqErr error
		var req requestData[M]
		if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
			log.Error(reqErr)
			ResponseJSON(c, CodeFailure)
			return
		}
		if reqErr == io.EOF {
			log.Warn("empty request body")
		}

		typ := reflect.TypeOf(*new(M)).Elem()
		val := reflect.New(typ).Interface().(M)
		_ = val
		for _, m := range req.Items {
			m.SetCreatedBy(c.GetString(consts.CTX_USERNAME))
			m.SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
			log.Infoz("batch create", zap.Bool("atomic", req.Options.Atomic), zap.Object(typ.Name(), m))
		}

		// 1.Perform business logic processing before batch create resource.
		if err = new(service.Factory[M]).Service().BatchCreateBefore(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// 2.Batch create resource in database.
		if reqErr != io.EOF {
			if err = handler(helper.NewDatabaseContext(c)).WithExpand(val.Expands()).Create(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure)
				return
			}
		}
		// 3.Perform business logic processing after batch create resource
		if err = new(service.Factory[M]).Service().BatchCreateAfter(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// TODO: write operation log

		// FIXME: 如果某些字段增加了 gorm unique tag, 则更新成功后的资源 ID 时随机生成的，并不是数据库中的

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

// BatchDelete
func BatchDelete[M types.Model](c *gin.Context) {
	BatchDeleteFactory[M]()(c)
}

// BatchDeleteFactory
func BatchDeleteFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	_ = handler
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_BATCH_DELETE)
		log.Info("batch delete")

		var err error
		var reqErr error
		var req requestData[M]
		if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
			log.Error(reqErr)
			ResponseJSON(c, CodeFailure)
			return
		}
		if reqErr == io.EOF {
			log.Warn("empty request body")
		}

		// 1.Perform business logic processing before batch delete resources.
		if err = new(service.Factory[M]).Service().BatchDeleteBefore(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Batch delete resources in database.
		if reqErr != io.EOF {
			if err = handler(helper.NewDatabaseContext(c)).WithPurge(req.Options.Purge).Delete(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure)
				return
			}
		}
		// 3.Perform business logic processing after batch delete resources.
		if err = new(service.Factory[M]).Service().BatchDeleteAfter(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

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

// BatchUpdate
func BatchUpdate[M types.Model](c *gin.Context) {
	BatchUpdateFactory[M]()(c)
}

// BatchUpdateFactory
func BatchUpdateFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	_ = handler
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_BATCH_UPDATE)
		log.Info("batch update")

		var err error
		var reqErr error
		var req requestData[M]
		if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
			log.Error(reqErr)
			ResponseJSON(c, CodeFailure)
			return
		}
		if reqErr == io.EOF {
			log.Warn("empty request body")
		}

		// 1.Perform business logic processing before batch update resource.
		if err = new(service.Factory[M]).Service().BatchUpdateBefore(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Batch update resource in database.
		if reqErr != io.EOF {
			if err = handler(helper.NewDatabaseContext(c)).Update(req.Items...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure)
				return
			}
		}
		// 3.Perform business logic processing after batch update resource.
		if err = new(service.Factory[M]).Service().BatchUpdateAfter(helper.NewServiceContext(c), req.Items...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

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

// BatchUpdatePartial
func BatchUpdatePartial[M types.Model](c *gin.Context) {
	BatchUpdatePartialFactory[M]()(c)
}

// BatchUpdatePartialFactory
func BatchUpdatePartialFactory[M types.Model](cfg ...*types.ControllerConfig[M]) gin.HandlerFunc {
	handler, _ := extractConfig(cfg...)
	_ = handler
	return func(c *gin.Context) {
		log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.PHASE_BATCH_UPDATE_PARTIAL)
		log.Info("batch update partial")

		var err error
		var reqErr error
		var req requestData[M]
		if reqErr = c.ShouldBindJSON(&req); reqErr != nil && reqErr != io.EOF {
			log.Error(reqErr)
			ResponseJSON(c, CodeFailure)
			return
		}
		if reqErr == io.EOF {
			log.Warn("empty request body")
		}
		typ := reflect.TypeOf(*new(M)).Elem()
		var shouldUpdates []M
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
			partialUpdateValue(log, typ, oldVal, newVal)
			shouldUpdates = append(shouldUpdates, oldVal.Addr().Interface().(M))
		}

		// 1.Perform business logic processing before batch partial update resource.
		if err = new(service.Factory[M]).Service().BatchUpdatePartialBefore(helper.NewServiceContext(c), shouldUpdates...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// 2.Batch partial update resource in database.
		if reqErr != io.EOF {
			if err = handler(helper.NewDatabaseContext(c)).Update(shouldUpdates...); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure)
				return
			}
		}
		// 3.Perform business logic processing after batch partial update resource.
		if err := new(service.Factory[M]).Service().BatchUpdatePartialAfter(helper.NewServiceContext(c), shouldUpdates...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

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

		// 1.Perform business logic processing before list resources.
		if err = new(service.Factory[M]).Service().ListBefore(helper.NewServiceContext(c), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		sortBy, _ := c.GetQuery(consts.QUERY_SORTBY)
		_, _ = page, size
		svc := new(service.Factory[M]).Service()
		svcCtx := helper.NewServiceContext(c)
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
			ResponseJSON(c, CodeFailure)
			return
		}
		// 3.Perform business logic processing after list resources.
		if err = new(service.Factory[M]).Service().ListAfter(helper.NewServiceContext(c), &data); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		log.Info("export data length: ", len(data))
		// 4.Export
		exported, err := new(service.Factory[M]).Service().Export(helper.NewServiceContext(c), data...)
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
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
			ResponseJSON(c, CodeFailure)
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
			ResponseJSON(c, CodeFailure)
			return
		}
		defer fd.Close()

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, fd); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		// filetype must be png or jpg.
		filetype, mime := filetype.DetectBytes(buf.Bytes())
		_, _ = filetype, mime

		// check filetype

		ml, err := new(service.Factory[M]).Service().Import(helper.NewServiceContext(c), buf)
		if err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}

		// service layer already create/update the records in database, just update fields "created_by", "updated_by".
		for i := range ml {
			ml[i].SetCreatedBy(c.GetString(consts.CTX_USERNAME))
			ml[i].SetUpdatedBy(c.GetString(consts.CTX_USERNAME))
		}
		if err := handler(helper.NewDatabaseContext(c)).Update(ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		if err := new(service.Factory[M]).Service().UpdateAfter(helper.NewServiceContext(c), ml...); err != nil {
			log.Error(err)
			ResponseJSON(c, CodeFailure)
			return
		}
		for i := range ml {
			if err := ml[i].UpdateAfter(); err != nil {
				log.Error(err)
				ResponseJSON(c, CodeFailure)
				return
			}
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
