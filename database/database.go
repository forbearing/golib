package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types"
	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"
)

// references:

var (
	ErrInvalidDB           = errors.New("invalid database, maybe not initialized")
	ErrUnknowDBType        = errors.New("unknow database type, only support mysql or sqlite")
	ErrNotPtrStruct        = errors.New("model is not pointer to structure")
	ErrNotPtrSlice         = errors.New("not pointer to slice")
	ErrNotPtrInt64         = errors.New("not pointer to int64")
	ErrNotAddressableModel = errors.New("model is not addressable")
	ErrNotAddressableSlice = errors.New("slice is not addressable")
	ErrNotSetSlice         = errors.New("slice cannot set")
)

var DB *gorm.DB
var defaultLimit = 1000

// func Init() (err error) {
// 	if err = mysql.Init(); err != nil {
// 		return
// 	}
// 	if err = redis.Init(); err != nil {
// 		return
// 	}
// 	DB = mysql.DB
//
// 	// create the table records that must be pre-exists before database curds.
// 	for _, r := range model.Records {
// 		// FIXME: 如何 preload, 来递归创建表数据
// 		// for i := range r.Expands {
// 		// 	DB = DB.Preload(r.Expands[i])
// 		// }
// 		if err = DB.Model(r.Table).Save(r.Rows).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return
// }

// database inplement types.Database[T types.Model] interface.
type database[M types.Model] struct {
	mu sync.Mutex
	db *gorm.DB

	// options
	enablePurge bool   // delete resource permanently, not only update deleted_at field, only works on 'Delete' method.
	enableCache bool   // using cache or not, only works 'List', 'Get', 'Count' method.
	tableName   string // support multiple custom table name, always used with the `WithDB` method.
	batchSize   int    // batch size for bulk operations. affects Create, Update, Delete.
}

// reset will reset the database interface to default value.
// Dont forget to call this method in all functions except option functions that prefixed with 'With*'.
func (db *database[M]) reset() {
	// default not delete resource permanently.
	// call option method 'WithPurge' to set enablePurge to true.
	db.enablePurge = false
	db.enableCache = false
	db.tableName = ""
	db.batchSize = 0
}

// WithDB returns a new database manipulator, only support *gorm.DB.
func (db *database[M]) WithDB(x any) types.Database[M] {
	// if x is nil, return the default database.
	if x == nil {
		return db
	}
	// if x type is not *gorm.DB, return the default database manipulator.
	_db, ok := x.(*gorm.DB)
	if !ok {
		return db
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if strings.ToLower(config.App.LogLevel) == "debug" {
		db.db = _db.WithContext(context.TODO()).Debug().Limit(defaultLimit)
	} else {
		db.db = _db.WithContext(context.TODO()).Limit(defaultLimit)
	}
	return db
}

// WithTable multiple custom table, always used with the method `WithDB`.
func (db *database[M]) WithTable(name string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.tableName = name
	return db
}

// WithBatchSize set batch size for bulk operations. affects Create, Update, Delete.
func (db *database[M]) WithBatchSize(size int) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.batchSize = size
	return db
}

// WithDebug setting debug mode, the priority is higher than config.Server.LogLevel and default value(false).
func (db *database[M]) WithDebug() types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Debug()
	return db
}

// WithQuery
// Examples:
// - WithQuery(&model.JobHistory{JobID: req.ID})
// - WithQuery(&model.CronJobHistory{CronJobID: req.ID})
// It will using mysql fuzzy matching if fuzzyMatch[0] is ture.
//
// NOTE: The underlying type msut be pointer to struct, otherwise panic will occur.
//
// NOTE: empty query conditions will casee list all resources from database.
func (db *database[M]) WithQuery(query M, fuzzyMatch ...bool) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	var _fuzzyMatch bool
	if len(fuzzyMatch) > 0 {
		_fuzzyMatch = fuzzyMatch[0]
	}

	typ := reflect.TypeOf(query).Elem()
	val := reflect.ValueOf(query).Elem()
	q := make(map[string]string)
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).IsZero() {
			continue
		}
		switch typ.Field(i).Type.Kind() {
		case reflect.Chan, reflect.Map, reflect.Func:
			continue
		case reflect.Struct:
			// All `model.XXX` extends the basic model named `Base`,
			if !val.Field(i).FieldByName("ID").IsZero() {
				// Not overwrite the "ID" value set in types.Model.
				// The "ID" value set in types.Model has higher priority than base model.
				if _, loaded := q["id"]; !loaded {
					q["id"] = val.Field(i).FieldByName("ID").Interface().(string)
				}
			}
			continue
		}
		// not typ.Field(i).Name
		jsonTagStr := typ.Field(i).Tag.Get("json")
		jsonTagItems := strings.Split(jsonTagStr, ",")
		jsonTag := ""
		if len(jsonTagItems) == 0 {
			continue
		}
		jsonTag = jsonTagItems[0]

		var v = val.Field(i).Interface()
		var _v string
		switch val.Field(i).Kind() {
		case reflect.Bool:
			// 由于 WHERE IN 语句会自动加上单引号,比如 WHERE `default` IN ('true')
			// 但是我们想要的是 WHERE `default` IN (true),
			// 所以没办法就只能直接转成 int 了.
			_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_v = fmt.Sprintf("%d", v)
		case reflect.Float32, reflect.Float64:
			_v = fmt.Sprintf("%g", v)
		case reflect.String:
			_v = fmt.Sprintf("%s", v)
		case reflect.Pointer:
			v = val.Field(i).Elem().Interface()
			// switch typ.Elem().Kind() {
			switch val.Field(i).Elem().Kind() {
			case reflect.Bool:
				_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				_v = fmt.Sprintf("%d", v)
			case reflect.Float32, reflect.Float64:
				_v = fmt.Sprintf("%g", v)
			case reflect.String:
				_v = fmt.Sprintf("%s", v)
			default:
				_v = fmt.Sprintf("%v", v)
			}
		case reflect.Slice:
			_len := val.Field(i).Len()
			if _len == 0 {
				logger.Database.Warn("reflect.Slice length is 0")
				_len = 1
			}
			slice := reflect.MakeSlice(val.Field(i).Type(), _len, _len)
			// fmt.Println("--------------- slice element", slice.Index(0), slice.Index(0).Kind(), slice.Index(0).Type())
			switch slice.Index(0).Kind() {
			case reflect.String: // handle string slice.
				// WARN: val.Field(i).Type() is model.GormStrings not []string,
				// execute statement `slice.Interface().([]string)` directly will case panic.
				// _v = strings.Join(slice.Interface().([]string), ",") // the slice type is GormStrings not []string.
				// We should make the slice of []string again.
				slice = reflect.MakeSlice(reflect.TypeOf([]string{}), _len, _len)
				reflect.Copy(slice, val.Field(i))
				_v = strings.Join(slice.Interface().([]string), ",")
			default:
				_v = fmt.Sprintf("%v", v)
			}
		default:
			_v = fmt.Sprintf("%v", v)
		}

		// json tag name naming format must be same as gorm table columns,
		// both should be "snake case" or "camel case".
		// gorm table columns naming format default to 'snake case',
		// so the json tag name is converted to "snake case here"
		// q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
		q[strcase.SnakeCase(jsonTag)] = _v
	}

	if _fuzzyMatch {
		// // Deprecated!
		// for k, v := range q {
		// 	// WARN: THE SQL STATEMENT MUST CONTAINS backticks ``.
		// 	db.db = db.db.Where(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
		// }

		// TODO: empty query conditions will list all resource from database.
		// We should make sure nothing record will be matched.
		// db.db = db.db.Where(`1 = 0`)

		// If the query strings has multiple value(seperated by ',')
		// construct the 'WHERE' 'REGEXP' SQL statement
		// eg: SELECT * FROM `assets` WHERE `category_level2_id` REGEXP '.*XS.*|.*NU.*'
		//     SELECT count(*) FROM `assets` WHERE `category_level2_id` REGEXP '.*XS.*|.*NU.*'
		for k, v := range q {
			items := strings.Split(v, ",")
			// TODO: should we skip if items length is 0?
			var regexpVal string
			for _, item := range items {
				regexpVal = regexpVal + "|.*" + item + ".*"
			}
			regexpVal = strings.TrimPrefix(regexpVal, "|")
			db.db = db.db.Where(fmt.Sprintf("`%s` REGEXP ?", k), regexpVal)
		}
	} else {
		// // Deprecated!
		// // SELECT * FROM `assets` WHERE `assets`.`category_level2_id` = 'NU
		// // SELECT count(*) FROM `assets` WHERE `assets`.`category_level2_id` = 'NU'
		// db.db = db.db.Where(query)

		// TODO: empty query conditions will list all resource from database.
		// We should make sure nothing record will be matched.
		// db.db = db.db.Where(`1 = 0`)

		// If the query string has multiple value(seperated by ','),
		// construct the 'WHERE' 'IN' SQL statement.
		// eg: SELECT id FROM users WHERE name IN ('user01', 'user02', 'user03', 'user04')
		for k, v := range q {
			items := strings.Split(v, ",")
			// TODO: should we skip if items length is 0?
			db.db = db.db.Where(fmt.Sprintf("`%s` IN (?)", k), items)
		}
	}
	return db
}

// WithQueryRaw
// Examples:
// - WithQueryRaw("name = ?", "hybfkuf")
// - WithQueryRaw("name <> ?", "hybfkuf")
// - WithQueryRaw("name IN (?)", []string{"hybfkuf", "hybfkuf 2"})
// - WithQueryRaw("name LIKE ?", "%hybfkuf%")
// - WithQueryRaw("name = ? AND age >= ?", "hybfkuf", "100")

// - WithQueryRaw("updated_at > ?", lastWeek)
// - WithQueryRaw("created_at BETWEEN ? AND ?", lastWeek, today)
func (db *database[M]) WithQueryRaw(query any, args ...any) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Where(query, args...)
	return db
}

// WithOrder
// reference: https://www.cnblogs.com/Braveliu/p/10654091.html
// For example:
// - WithOrder("name") // default ASC.
// - WithOrder("name desc")
// - WithOrder("created_at")
// - WithOrder("updated_at desc")
// NOTE: you cannot using the mysql keyword, such as: "order", "limit".
func (db *database[M]) WithOrder(order string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	// 可以多多个字段进行排序, 每个字段之间通过逗号分隔,
	// order 的值比如: "field1, field2 desc, field3 asc"
	items := strings.Split(order, ",")
	for _, _order := range items {
		if len(order) != 0 {
			items := strings.Fields(_order)
			for i := range items {
				if strings.EqualFold(items[i], "asc") || strings.EqualFold(items[i], "desc") {
					items[i] = strings.ToUpper(items[i])
				} else {
					// 第一个是排序字段,必须加上反引号,因为排序的字符串可能是 sql 语句关键字
					// 第二个可能是 "desc" 等关键字不需要加反引号
					// items[0] = "`" + items[0] + "`"
					// 如果不是关键字都加上反引号
					items[i] = "`" + items[i] + "`"
				}
			}
			_orders := strings.Join(items, " ")
			db.db = db.db.Order(_orders)
		}
	}
	return db
}

// WithScope
// Examples:
//   - pageStr, _ := c.GetQuery("page")
//     sizeStr, _ := c.GetQuery("size")
//     page, _ := strconv.Atoi(pageStr)
//     size, _ := strconv.Atoi(sizeStr)
//     WithScope(page, size)
func (db *database[M]) WithScope(page, size int) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = defaultLimit
	}
	offset := (page - 1) * size
	db.db = db.db.Scopes(func(d *gorm.DB) *gorm.DB {
		return d.Offset(offset).Limit(size)
	})
	return db
}

// WithLimit specifies the number of record to be retrieved.
func (db *database[M]) WithLimit(limit int) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Limit(limit)
	return db
}

// WithExpand preload associations with given conditions.
// order: preload with order.
// eg: [Children.Children.Children Parent.Parent.Parent]
// eg: [Children Parent]
//
// NOTE: WithExpand only workds on mysql foreign key relationship.
// If you want expand the custom field that without gorm tag about foregin key defination,
// you should define the GetAfter/ListAfter in model layer or service layoer.
func (db *database[M]) WithExpand(expand []string, order ...string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	var _orders string
	if len(order) > 0 {
		if len(order[0]) > 0 {
			items := strings.Fields(order[0])
			// 第一个是排序字段,必须加上反引号,因为排序的字符串可能是 sql 语句关键字
			// 第二个可能是 "desc" 等关键字不需要加反引号
			items[0] = "`" + items[0] + "`"
			_orders = strings.Join(items, " ")
		}
	}
	withOrder := func(db *gorm.DB) *gorm.DB {
		if len(_orders) > 0 {
			return db.Order(_orders)
		} else {
			return db
		}
	}
	// FIXME: 前端加了 _depth 查询参数, 但是层数不匹配就无法递归排序,
	// _depth 的作用:
	// _depth = 2: Children -> Children.Children
	// _depth = 3: Children -> Children.Children.Children
	// 假设一共有3层, 但是 _depth=5, 则无法递归排序
	//
	// 解决办法:
	// 假设: [Children.Children.Children, Parent]
	// 以前:
	//      db.db = db.db.Preload("Children.Children.Children", withOrder)
	//      db.db = db.db.Preload("Parent", withOrder)
	// 现在: (递归 Children)
	//      db.db = db.db.Preload("Children", withOrder)
	//      db.db = db.db.Preload("Children.Children", withOrder)
	//      db.db = db.db.Preload("Children.Children.Children", withOrder)
	//      db.db = db.db.Preload("Parent", withOrder)

	for i := range expand {
		// preload 排序问题
		// https://www.jianshu.com/p/a88fb2d4b2ef
		// https://gorm.io/docs/preload.html#Custom-Preloading-SQL

		items := strings.Split(expand[i], ".")
		switch len(items) {
		case 0:
		case 1:
			db.db = db.db.Preload(expand[i], withOrder)
		default:
			for j := range items {
				// fmt.Println("================== ", strings.Join(items[0:j+1], "."))
				db.db = db.db.Preload(strings.Join(items[0:j+1], "."), withOrder)
			}
		}
	}

	return db
}

// WithExclude excludes records that matchs a condition within a list.
// For example:
//   - If you want exlcude users with specific ids from your query,
//     you can use WithExclude(excludes),
//     excludes: "id" as key, ["myid1", "myid2", "myid3"] as value.
//   - If you want excludes users that id not ["myid1", "myid2"] and not not ["root", "noname"],
//     the `excludes` should be:
//     excludes := make(map[string][]any)
//     excludes["id"] = []any{"myid1", "myid2"}
//     excludes["name"] = []any{"root", "noname"}.
func (db *database[M]) WithExclude(excludes map[string][]any) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	for k, v := range excludes {
		db.db = db.db.Not(k, v)
	}
	return db
}

// WithPurge will delete resource in database permanently.
// It only works on 'Delete' method.
func (db *database[M]) WithPurge(enable ...bool) types.Database[M] {
	_enable := true
	if len(enable) > 0 {
		_enable = enable[0]
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	db.enablePurge = _enable
	return db
}

// WithCache will make query resource count from cache.
// If cache not found or expired. query from database directly.
func (db *database[M]) WithCache(enable ...bool) types.Database[M] {
	_enable := true
	if len(enable) > 0 {
		_enable = enable[0]
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	db.enableCache = _enable
	return db
}

// WithOmit omit specific columns when create/update/query.
func (db *database[M]) WithOmit(columns ...string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Omit(columns...)
	return db
}

// Create one or multiple objects in database.
// It will update the "created_at" and "updated_at" field.
// Examples:
// - database.XXX().Create(&model.XXX{ID: id, Field1: field1, Field2: field2})
func (db *database[M]) Create(objs ...M) error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	if config.App.RedisConfig.Enable {
		defer func() {
			go func() {
				begin := time.Now()
				prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "create")
				defer logger.Database.Infow("remove cache after create", "cost", time.Since(begin).String(), "prefix", prefix)
				if err := redis.RemovePrefix(prefix); err != nil {
					logger.Database.Errorw("failed to remove cache keys", err, "action", "create")
				}
			}()
		}()
	}
	// if len(objs) == 0 {
	// 	return nil
	// }
	// for i := range objs {
	// 	objs[i].SetID() // set id when id is empty.
	// }
	// return db.db.Create(objs).Error
	if len(objs) == 0 {
		return nil
	}
	// Invoke model hook: CreateBefore.
	for i := range objs {
		if err := objs[i].CreateBefore(); err != nil {
			return err
		}
	}
	for i := range objs {
		objs[i].SetID() // set id when id is empty.
	}

	// TODO: batch size mode
	// var shouldSplit bool
	// if db.batchSize > 0 {
	// 	if len(objs) > int(db.batchSize) {
	// 		shouldSplit = true
	// 	}
	// }
	// if shouldSplit {
	// 	for i := 0; i < len(objs); i += db.batchSize {
	// 		if i+db.batchSize > len(objs) {
	// 			if err := db.db.Table(db.tableName).Save(objs[i:]).Error; err != nil {
	// 				return err
	// 			}
	// 		} else {
	// 			if err := db.db.Table(db.tableName).Save(objs[i : i+db.batchSize]).Error; err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// } else {
	// 	if err := db.db.Table(db.tableName).Save(objs).Error; err != nil {
	// 		return err
	// 	}
	// }

	// because db.db.Delete method just update field "delete_at" to current time,
	// not really delete it(soft delete).
	// If record already exists, Update method update all fields but exlcude "created_at" by
	// mysql "ON DUPLICATE KEY UPDATE" mechanism. so we should call UpdateById to
	// specially update the "created_at" field.
	// if err := db.db.Save(objs).Error; err != nil {
	if err := db.db.Table(db.tableName).Save(objs).Error; err != nil {
		return err
	}
	for i := range objs {
		// 这里要重新创建一个 gorm.DB 实例, 否则会出现这种语句, id 出现多次了.
		// UPDATE `assets` SET `created_at`='2023-11-12 14:35:42.604',`updated_at`='2023-11-12 14:35:42.604' WHERE id = '010103NU000020' AND `assets`.`deleted_at` IS NULL AND id = '010103NU000021' AND id = '010103NU000022' LIMIT 1000
		var _db *gorm.DB
		if strings.ToLower(config.App.LogLevel) == "debug" {
			_db = DB.Debug()
		} else {
			_db = DB
		}
		// if err := _db.Model(*new(M)).Where("id = ?", objs[i].GetID()).Update("created_at", time.Now()).Error; err != nil {
		if err := _db.Table(db.tableName).Model(*new(M)).Where("id = ?", objs[i].GetID()).Update("created_at", time.Now()).Error; err != nil {
			return err
		}
	}
	// Invoke model hook: CreateAfter.
	for i := range objs {
		if err := objs[i].CreateAfter(); err != nil {
			return err
		}
	}
	return nil
}

// Delete one or multiple objects in database.
// Examples:
// - database.XXX().Delete(&model.XXX{ID: id}) // delete record with primary key
// - database.XXX().WithQuery(req).Delete(req) // delete record with where condiions.
func (db *database[M]) Delete(objs ...M) error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	if config.App.RedisConfig.Enable {
		defer func() {
			// TODO:only delete cache of all list statement and cache for current get statements.
			go func() {
				begin := time.Now()
				prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "delete")
				defer logger.Database.Infow("remove cache after delete", "cost", time.Since(begin).String(), "prefix", prefix)
				if err := redis.RemovePrefix(prefix); err != nil {
					logger.Database.Errorw("failed to remove cache keys", err, "action", "delete")
				}
			}()
		}()
	}

	if len(objs) == 0 {
		return nil
	}
	// Invoke model hook: DeleteBefore.
	for i := range objs {
		if err := objs[i].DeleteBefore(); err != nil {
			return nil
		}
	}
	if db.enablePurge {
		// delete permanently.
		// if err := db.db.Unscoped().Delete(objs).Error; err != nil {
		if err := db.db.Table(db.tableName).Unscoped().Delete(objs).Error; err != nil {
			return err
		}
	} else {
		// Delete() method just update field "delete_at" to currrent time.
		// DO NOT FORGET update the "created_at" field when create/update if record already exists.
		// if err := db.db.Delete(objs).Error; err != nil {
		if err := db.db.Table(db.tableName).Delete(objs).Error; err != nil {
			return err
		}
	}
	// Invoke model hook: DeleteAfter.
	for i := range objs {
		if err := objs[i].DeleteAfter(); err != nil {
			return err
		}
	}
	return nil
}

// Update one or multiple objects in database.
// It will just update the "updated_at" field.
// Examples:
//   - obj := new(model.XXX)
//     objs := make([]*model.XXX, 0)
//     // do something on objs and objs.
//     doSomething(obj)
//     doSomething(objs)
//     database.XXX().Update(obj)
//     database.XXX().Update(objs)
func (db *database[M]) Update(objs ...M) error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	if config.App.RedisConfig.Enable {
		defer func() {
			go func() {
				begin := time.Now()
				prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "update")
				defer logger.Database.Infow("remove cache after update", "cost", time.Since(begin).String(), "prefix", prefix)
				if err := redis.RemovePrefix(prefix); err != nil {
					logger.Database.Errorw("failed to remove cache keys", err, "action", "update")
				}
			}()
		}()
	}
	if len(objs) == 0 {
		return nil
	}
	// Invoke model hook: UpdateBefore.
	for i := range objs {
		if err := objs[i].UpdateBefore(); err != nil {
			return err
		}
	}
	for i := range objs {
		objs[i].SetID() // set id when id is empty
	}
	// if err := db.db.Save(objs).Error; err != nil {
	if err := db.db.Table(db.tableName).Save(objs).Error; err != nil {
		return err
	}
	// Invoke model hook: UpdateAfter.
	for i := range objs {
		if err := objs[i].UpdateAfter(); err != nil {
			return err
		}
	}
	return nil
}
func (db *database[M]) UpdateById(id any, key string, val any) error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	if config.App.RedisConfig.Enable {
		defer func() {
			go func() {
				begin := time.Now()
				prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "update_by_id")
				defer logger.Database.Infow("remove cache after update_by_id", "cost", time.Since(begin).String(), "prefix", prefix)
				if err := redis.RemovePrefix(prefix); err != nil {
					logger.Database.Errorw("failed to remove cache keys", err, "action", "update")
				}
			}()
		}()
	}
	// return db.db.Model(*new(M)).Where("id = ?", id).Update(key, val).Error
	return db.db.Table(db.tableName).Model(*new(M)).Where("id = ?", id).Update(key, val).Error
}

// List find all record if not run WithQuery(query).
// List find record with `where` conditioan if run WithQuery(query).
// Examples:
//   - data := make([]*model.XXX, 0)
//     database.XXX().WithScope(page, size).WithQuery(&model.XXX{Field1: field1, Field2: field2}).List(&data)
//   - data := make([]*model.XXX, 0)
//     database.XXX().List(&data)
func (db *database[M]) List(dest *[]M, cache ...*[]byte) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	if dest == nil {
		return nil
	}
	begin := time.Now()
	var key string
	var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	var _cache []byte
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Find(dest).Statement, "list")
	if len(cache) == 0 {
		shouldDecode = true
	} else {
		if cache[0] == nil {
			shouldDecode = true
		}
	}
	if shouldDecode {
		var _dest []M
		if _dest, err = redis.GetML[M](key); err == nil {
			val := reflect.ValueOf(dest)
			if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			if !val.Elem().CanSet() {
				return ErrNotSetSlice
			}
			val.Elem().Set(reflect.ValueOf(_dest))
			logger.Database.Infow("list and decode from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	} else {
		if _cache, err = redis.Get(key); err == nil {
			val := reflect.ValueOf(cache[0])
			if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			if !val.Elem().CanSet() {
				return ErrNotSetSlice
			}
			val.Elem().Set(reflect.ValueOf(_cache))
			logger.Database.Infow("list from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	}
	if !errors.Is(err, redis.ErrKeyNotExists) {
		logger.Database.Error(err)
		return err
	}
	// Not Found cache and continue.

QUERY:
	// Invoke model hook: ListBefore.
	for i := range *dest {
		if err = (*dest)[i].ListBefore(); err != nil {
			return err
		}
	}
	// if err = db.db.Find(dest).Error; err != nil {
	if err = db.db.Table(db.tableName).Find(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: ListAfter()
	for i := range *dest {
		if err = (*dest)[i].ListAfter(); err != nil {
			return err
		}
	}
	// Cache the result.
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("list from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.SetML[M](key, *dest); err != nil {
				logger.Database.Error(err)
			}
		}()
	}

	return nil
}

// // Find equal to WithQuery(condition).List()
// // More detail see `List` document.
// func (db *database[T]) Find(dest *[]T, query T) error {
// 	return db.db.Where(query).Find(dest).Error
// }

// Get find one record accoding to `id`.
func (db *database[M]) Get(dest M, id string, cache ...*[]byte) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	var key string
	var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	begin := time.Now()
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	// _, key = BuildKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Find(dest).Statement, "get")
	// 我发现这个 id 的值怎么都无法填充到 sql 语句中, 所以直接用 id 作为 key 的一部分,而不是用 sql 语句
	// 如果不用 id 作为 redis key, 那么多个 get 的语句相同则 key 相同,获取到的都是同一个缓存.
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Find(dest).Statement, "get", id)
	if len(cache) == 0 {
		shouldDecode = true
	} else {
		if cache[0] == nil {
			shouldDecode = true
		}
	}
	if shouldDecode { // query cache from redis and decoded into 'dest'.
		var _dest M
		if _dest, err = redis.GetM[M](key); err == nil {
			val := reflect.ValueOf(dest)
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrStruct
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableModel
			}
			val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
			logger.Database.Infow("get and decode from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	} else {
		var _cache []byte
		if _cache, err = redis.Get(key); err == nil {
			val := reflect.ValueOf(cache[0])
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			val.Elem().Set(reflect.ValueOf(_cache))
			logger.Database.Infow("get from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	}
	if err != redis.ErrKeyNotExists {
		logger.Database.Error(err)
		return err
	}
	// Not Found cache, continue query database.

QUERY:
	// Invoke model hook: GetBefore.
	if err = dest.GetBefore(); err != nil {
		return err
	}
	// if err = db.db.Where("id = ?", id).Find(dest).Error; err != nil {
	if err = db.db.Table(db.tableName).Where("id = ?", id).Find(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter.
	if err = dest.GetAfter(); err != nil {
		return err
	}
	// Cache the result.
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("get from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.SetM[M](key, dest); err != nil {
				logger.Database.Error(err)
			}
		}()
	}
	return nil
}

// Count returns the total number of records with the given query condition.
// NOTE: The underlying type msut be pointer to struct, otherwise panic will occur.
func (db *database[M]) Count(count *int64) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	var key string
	var _cache int64
	begin := time.Now()
	if count == nil {
		return nil
	}
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Model(*new(M)).Count(count).Statement, "count")
	if _cache, err = redis.GetInt(key); err == nil {
		*count = _cache
		logger.Database.Infow("count from cache", "cost", time.Since(begin).String(), "key", key)
		return
	}
	if !errors.Is(err, redis.ErrKeyNotExists) {
		logger.Database.Error(err)
		return err
	}
	// NOT FOUND cache, continue query.

QUERY:
	// if err = db.db.Model(*new(M)).Count(count).Error; err != nil {
	if err = db.db.Table(db.tableName).Model(*new(M)).Count(count).Error; err != nil {
		logger.Database.Error(err)
		return err
	}
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("count from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.Set(key, *count); err != nil {
				logger.Database.Error(err)
			}
		}()

	}
	return nil
}

// First finds the first record ordered by primary key.
func (db *database[M]) First(dest M, cache ...*[]byte) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	var key string
	var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	begin := time.Now()
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "first")
	if len(cache) == 0 {
		shouldDecode = true
	} else {
		if cache[0] == nil {
			shouldDecode = true
		}
	}

	if shouldDecode { // query cache from redis and decode into 'dest'.
		var _dest M
		if _dest, err = redis.GetM[M](key); err == nil {
			val := reflect.ValueOf(dest)
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrStruct
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableModel
			}
			val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
			logger.Database.Infow("first and decode from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	} else {
		var _cache []byte
		if _cache, err = redis.Get(key); err == nil {
			val := reflect.ValueOf(cache[0])
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			val.Elem().Set(reflect.ValueOf(_cache))
			logger.Database.Infow("first from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
		if err != redis.ErrKeyNotExists {
			logger.Database.Error(err)
			return err
		}
	}
	// Not Found cache, continue query database.

QUERY:
	// Invoke model hook: GetBefore
	if err = dest.GetBefore(); err != nil {
		return err
	}
	// if err = db.db.First(dest).Error; err != nil {
	if err = db.db.Table(db.tableName).First(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter
	if err = dest.GetAfter(); err != nil {
		return err
	}
	// Cache the result.
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("first from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.SetM[M](key, dest); err != nil {
				logger.Database.Error(err)
			}
		}()
	}
	return nil
}

// Last finds the last record ordered by primary key.
func (db *database[M]) Last(dest M, cache ...*[]byte) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	var key string
	var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	begin := time.Now()
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "last")
	if len(cache) == 0 {
		shouldDecode = true
	} else {
		if cache[0] == nil {
			shouldDecode = true
		}
	}
	if shouldDecode { // query cache from redis and decode into 'dest'.
		var _dest M
		if _dest, err = redis.GetM[M](key); err == nil {
			val := reflect.ValueOf(dest)
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrStruct
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableModel
			}
			val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
			logger.Database.Infow("last and decode from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	} else {
		var _cache []byte
		if _cache, err = redis.Get(key); err == nil {
			val := reflect.ValueOf(cache[0])
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			val.Elem().Set(reflect.ValueOf(_cache))
			logger.Database.Infow("last from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	}
	if err != redis.ErrKeyNotExists {
		logger.Database.Error(err)
		return err
	}
	// Not Found cache, continue query database.

QUERY:
	// Invoke model hook: GetBefore.
	if err = dest.GetBefore(); err != nil {
		return err
	}
	// if err = db.db.Last(dest).Error; err != nil {
	if err = db.db.Table(db.tableName).Last(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter
	if err = dest.GetAfter(); err != nil {
		return err
	}
	// Cache the result.
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("last from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.SetM[M](key, dest); err != nil {
				logger.Database.Error(err)
			}
		}()
	}
	return nil
}

// Take finds the first record returned by the database in no specified order.
func (db *database[M]) Take(dest M, cache ...*[]byte) (err error) {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	var key string
	var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	begin := time.Now()
	if !db.enableCache {
		goto QUERY
	}
	if !config.App.RedisConfig.Enable {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "take")
	if len(cache) == 0 {
		shouldDecode = true
	} else {
		if cache[0] == nil {
			shouldDecode = true
		}
	}
	if shouldDecode { // query cache from redis and decode into 'dest'.
		var _dest M
		if _dest, err = redis.GetM[M](key); err == nil {
			val := reflect.ValueOf(dest)
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrStruct
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableModel
			}
			val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
			logger.Database.Infow("get and decode from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	} else {
		var _cache []byte
		if _cache, err = redis.Get(key); err == nil {
			val := reflect.ValueOf(cache[0])
			if val.Kind() != reflect.Ptr {
				return ErrNotPtrSlice
			}
			if !val.Elem().CanAddr() {
				return ErrNotAddressableSlice
			}
			val.Elem().Set(reflect.ValueOf(_cache))
			logger.Database.Infow("take from cache", "cost", time.Since(begin).String(), "key", key)
			return nil // Found cache and return.
		}
	}
	if err != redis.ErrKeyNotExists {
		logger.Database.Error(err)
		return err
	}
	// Not Found cache, continue query database.

QUERY:
	// Invoke model hook: GetBefore.
	if err = dest.GetBefore(); err != nil {
		return err
	}
	// if err = db.db.Take(dest).Error; err != nil {
	if err = db.db.Table(db.tableName).Take(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter.
	if err = dest.GetAfter(); err != nil {
		return err
	}
	// Cache the result.
	if db.enableCache && config.App.RedisConfig.Enable {
		logger.Database.Infow("take from database", "cost", time.Since(begin).String(), "key", key)
		go func() {
			if err = redis.SetM[M](key, dest); err != nil {
				logger.Database.Error(err)
			}
		}()

	}
	return nil
}

// Cleanup delete all records that column 'deleted_at' is not null.
func (db *database[M]) Cleanup() error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	defer db.reset()
	// return db.db.Limit(-1).Where("deleted_at IS NOT NULL").Model(*new(M)).Unscoped().Delete(make([]M, 0)).Error
	return db.db.Table(db.tableName).Limit(-1).Where("deleted_at IS NOT NULL").Model(*new(M)).Unscoped().Delete(make([]M, 0)).Error
}

// Database implement interface types.Database, its default database manipulator.
func Database[M types.Model](ctx ...context.Context) types.Database[M] {
	if DB == nil || DB == new(gorm.DB) {
		panic("database is not initialized")
	}
	c := context.TODO()
	if len(ctx) > 0 {
		if ctx[0] != nil {
			c = ctx[0]
		}
	}
	if strings.ToLower(config.App.LogLevel) == "debug" {
		return &database[M]{db: DB.WithContext(c).Debug().Limit(defaultLimit)}
	}
	return &database[M]{db: DB.WithContext(c).Limit(defaultLimit)}
}

// // User returns a generic database manipulator with the `curd` capabilities
// // for *model.User to create/delete/update/list/get in database.
// // The database type deponds on the value of config.Server.DBType.
// func User(ctx ...context.Context) types.Database[*model.User] {
// 	c := context.TODO()
// 	if len(ctx) > 0 {
// 		if ctx[0] != nil {
// 			c = ctx[0]
// 		}
// 	}
// 	if strings.ToLower(config.App.LogLevel) == "debug" {
// 		return &database[*model.User]{db: DB.WithContext(c).Debug().Limit(defaultLimit)}
// 	}
// 	return &database[*model.User]{db: DB.WithContext(c).Limit(defaultLimit)}
// }

// // Transaction start a transaction as a block, return error will rollback, otherwise to commit.
// // Transaction executes an arbitrary number of commands in fc within a transaction.
// // On success the changes are committed; if an error occurs they are rolled back.
// func Transaction(fn func(tx *gorm.DB) error) error {
// 	return mysql.DB.Transaction(fn)
// }
//
// // Exec executes raw sql
// func Exec(sql string, values any) error {
// 	return mysql.DB.Exec(sql, values).Error
// }

// buildCacheKey build redis prefix and redis key.
// The Prefix consist of "redis key prefix" + "table name".
// The KEY consist of "redis key prefix" + "table name" + "sql statement".
//
// More details see reference: https://gorm.io/docs/sql_builder.html
func buildCacheKey(stmt *gorm.Statement, action string, id ...string) (prefix, key string) {
	prefix = strings.Join([]string{config.App.RedisConfig.Namespace, stmt.Table}, ":")
	switch strings.ToLower(action) {
	case "get":
		if len(id) > 0 {
			key = strings.Join([]string{config.App.RedisConfig.Namespace, stmt.Table, action, id[0]}, ":")
		} else {
			key = strings.Join([]string{config.App.RedisConfig.Namespace, stmt.Table, action, stmt.SQL.String()}, ":")
		}
	default:
		key = strings.Join([]string{config.App.RedisConfig.Namespace, stmt.Table, action, stmt.SQL.String()}, ":")
	}
	return
}

// // Get find one record accoding to `id`.
//
//	func (db *database[T]) Get(dest T, id any, cache ...*[]byte) (err error) {
//		defer db.reset()
//		var key string
//		if config.App.RedisConfig.Enabled {
//			_, key = BuildKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?").Find(dest).Statement, "get")
//			// Found cache and return.
//			var _dest T
//			if _dest, err = redis.Get[T](key); err == nil {
//				val := reflect.ValueOf(dest)
//				if val.Kind() != reflect.Ptr {
//					return ErrNotPtrStruct
//				}
//				if !val.Elem().CanAddr() {
//					return ErrNotAddressableModel
//				}
//				val.Elem().Set(reflect.ValueOf(_dest).Elem())
//				fmt.Println("[Get] query from cache")
//				return nil
//			}
//			if err != redis.ErrKeyNotExists {
//				zap.S().Error(err)
//				return err
//			}
//			// Not Found cache, continue query database.
//		}
//
//		// Invoke model hook: GetBefore.
//		if err = dest.GetBefore(); err != nil {
//			return err
//		}
//		if err = db.db.Where("id = ?", id).Find(dest).Error; err != nil {
//			return err
//		}
//		if config.App.RedisConfig.Enabled {
//			fmt.Println(key)
//			// Cache the query result.
//			if err = redis.Set[T](key, dest); err != nil {
//				zap.S().Error(err)
//				return err
//			}
//			fmt.Println("[Get] query from database: ", key)
//		}
//		// Invoke model hook: GetAfter.
//		if err = dest.GetAfter(); err != nil {
//			return err
//		}
//		return nil
//	}
//

// // List find all record if not run WithQuery(query).
// // List find record with `where` conditioan if run WithQuery(query).
// // Examples:
// //   - data := make([]*model.XXX, 0)
// //     database.XXX().WithScope(page, size).WithQuery(&model.XXX{Field1: field1, Field2: field2}).List(&data)
// //   - data := make([]*model.XXX, 0)
// //     database.XXX().List(&data)
// func (db *database[T]) List(dest *[]T, cache ...*[]byte) (err error) {
// 	defer db.reset()
// 	if dest == nil {
// 		return nil
// 	}
// 	var key string
// 	if config.App.RedisConfig.Enabled {
// 		_, key = BuildKey(db.db.Session(&gorm.Session{DryRun: true}).Find(dest).Statement, "list")
// 		var _dest []T
// 		begin := time.Now()
// 		// Found cache and return.
// 		if _dest, err = redis.GetM[T](key); err == nil {
// 			val := reflect.ValueOf(dest)
// 			if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
// 				return ErrNotPtrSlice
// 			}
// 			if !val.Elem().CanAddr() {
// 				return ErrNotAddressableSlice
// 			}
// 			if !val.Elem().CanSet() {
// 				return ErrNotSetSlice
// 			}
// 			val.Elem().Set(reflect.ValueOf(_dest))
// 			fmt.Println("[List] query from cache", key, time.Since(begin))
// 			return nil
// 		}
// 		if !errors.Is(err, redis.ErrKeyNotExists) {
// 			zap.S().Error(err)
// 			return err
// 		}
// 		// Not Found cache and continue.
// 	}
//
// 	// Invoke model hook: ListBefore.
// 	for i := range *dest {
// 		if err = (*dest)[i].ListBefore(); err != nil {
// 			return err
// 		}
// 	}
// 	if err = db.db.Find(dest).Error; err != nil {
// 		return err
// 	}
// 	// Invoke model hook: ListAfter()
// 	for i := range *dest {
// 		if err = (*dest)[i].ListAfter(); err != nil {
// 			return err
// 		}
// 	}
//
// 	if config.App.RedisConfig.Enabled {
// 		begin := time.Now()
// 		if err = redis.SetM[T](key, *dest); err != nil {
// 			zap.S().Error(err)
// 			return err
// 		}
// 		fmt.Println("[List] query from database, costed: ", time.Since(begin), key)
// 	}
//
// 	return nil
// }

// // Count returns the total number of records with the given query condition.
// // NOTE: The underlying type msut be pointer to struct, otherwise panic will occur.
// func (db *database[T]) Count(count *int64, query T, fuzzyMatch ...bool) error {
// 	defer db.reset()
// 	var _fuzzyMatch bool
// 	if len(fuzzyMatch) > 0 {
// 		_fuzzyMatch = fuzzyMatch[0]
// 	}
// 	if _fuzzyMatch {
// 		typ := reflect.TypeOf(query).Elem()
// 		val := reflect.ValueOf(query).Elem()
// 		q := make(map[string]any)
// 		for i := 0; i < typ.NumField(); i++ {
// 			if val.Field(i).IsZero() {
// 				continue
// 			}
// 			switch typ.Field(i).Type.Kind() {
// 			case reflect.Chan, reflect.Map, reflect.Func:
// 				continue
// 			case reflect.Struct:
// 				// All `model.XXX` extends the basic model named `Base`,
// 				if !val.Field(i).FieldByName("ID").IsZero() {
// 					// Not overwrite the "ID" value set in types.Model.
// 					// The "ID" value set in types.Model has higher priority than base model.
// 					if _, loaded := q["id"]; !loaded {
// 						q["id"] = val.Field(i).FieldByName("ID").Interface()
// 					}
// 				}
// 				continue
// 			}
// 			// not typ.Field(i).Name
// 			jsonTagStr := typ.Field(i).Tag.Get("json")
// 			jsonTagItems := strings.Split(jsonTagStr, ",")
// 			jsonTag := ""
// 			if len(jsonTagItems) == 0 {
// 				continue
// 			}
// 			jsonTag = jsonTagItems[0]
// 			// json tag name naming format must be same as gorm table columns,
// 			// both should be "snake case" or "camel case".
// 			// gorm table columns naming format default to 'snake case',
// 			// so the json tag name is converted to "snake case here"
// 			q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
// 		}
// 		db.db = db.db.Model(*new(T))
// 		for k, v := range q {
// 			// db.db = db.db.Where(fmt.Sprintf("%s LIKE ?", k), fmt.Sprintf("%%%v%%", v))
// 			// WARN: THE SQL STATEMENT MUST CONTAINS backticks ``.
// 			db.db = db.db.Where(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
// 		}
// 	} else {
// 		db.db = db.db.Model(*new(T)).Where(query)
// 	}
// 	return db.db.Count(count).Error
// }

// 2023 12-09 14:18
// // WithQuery
// // Examples:
// // - WithQuery(&model.JobHistory{JobID: req.ID})
// // - WithQuery(&model.CronJobHistory{CronJobID: req.ID})
// // It will using mysql fuzzy matching if fuzzyMatch[0] is ture.
// // NOTE: The underlying type msut be pointer to struct, otherwise panic will occur.
// func (db *database[M]) WithQuery(query M, fuzzyMatch ...bool) types.Database[M] {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()
// 	var _fuzzyMatch bool
// 	if len(fuzzyMatch) > 0 {
// 		_fuzzyMatch = fuzzyMatch[0]
// 	}
// 	if _fuzzyMatch {
// 		typ := reflect.TypeOf(query).Elem()
// 		val := reflect.ValueOf(query).Elem()
// 		q := make(map[string]any)
// 		for i := 0; i < typ.NumField(); i++ {
// 			if val.Field(i).IsZero() {
// 				continue
// 			}
// 			switch typ.Field(i).Type.Kind() {
// 			case reflect.Chan, reflect.Map, reflect.Func:
// 				continue
// 			case reflect.Struct:
// 				// All `model.XXX` extends the basic model named `Base`,
// 				if !val.Field(i).FieldByName("ID").IsZero() {
// 					// Not overwrite the "ID" value set in types.Model.
// 					// The "ID" value set in types.Model has higher priority than base model.
// 					if _, loaded := q["id"]; !loaded {
// 						q["id"] = val.Field(i).FieldByName("ID").Interface()
// 					}
// 				}
// 				continue
// 			}
// 			// not typ.Field(i).Name
// 			jsonTagStr := typ.Field(i).Tag.Get("json")
// 			jsonTagItems := strings.Split(jsonTagStr, ",")
// 			jsonTag := ""
// 			if len(jsonTagItems) == 0 {
// 				continue
// 			}
// 			jsonTag = jsonTagItems[0]
// 			// json tag name naming format must be same as gorm table columns,
// 			// both should be "snake case" or "camel case".
// 			// gorm table columns naming format default to 'snake case',
// 			// so the json tag name is converted to "snake case here"
// 			q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
// 		}
// 		for k, v := range q {
// 			// db.db = db.db.Where(fmt.Sprintf("%s LIKE ?", k), fmt.Sprintf("%%%v%%", v))
// 			// WARN: THE SQL STATEMENT MUST CONTAINS backticks ``.
// 			db.db = db.db.Where(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
// 		}
// 	} else {
// 		// If the query string has multiple value(seperated by ','),
// 		// construct the 'WHERE' 'IN' SQL statement, eg: `SELECT id FROM users WHERE name IN ('user01', 'user02', 'user03', 'user04')`
// 		db.db = db.db.Where(query)
// 	}
// 	return db
// }

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
