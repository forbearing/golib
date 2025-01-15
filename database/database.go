package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/lru"
	"github.com/forbearing/golib/types"

	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
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

var (
	DB *gorm.DB

	defaultLimit           = 1000
	defaultBatchSize       = 1000
	defaultDeleteBatchSize = 10000
	defaultsColumns        = []string{
		"id",
		"created_by",
		"updated_by",
		"created_at",
		"updated_at",
		"deleted_at",
	}
)

// database inplement types.Database[T types.Model] interface.
type database[M types.Model] struct {
	mu  sync.Mutex
	db  *gorm.DB
	ctx context.Context

	// options
	enablePurge   bool   // delete resource permanently, not only update deleted_at field, only works on 'Delete' method.
	enableCache   bool   // using cache or not, only works 'List', 'Get', 'Count' method.
	tableName     string // support multiple custom table name, always used with the `WithDB` method.
	batchSize     int    // batch size for bulk operations. affects Create, Update, Delete.
	noHook        bool   // disable model hook.
	orQuery       bool   // or query
	inTransaction bool   // in transaction
	tryRun        bool   // try run

	shouldAutoMigrate bool
}

// reset will reset the database interface to default value.
// Dont forget to call this method in all functions except option functions that prefixed with 'With*'.
func (db *database[M]) reset() {
	db.mu.Lock()
	defer db.mu.Unlock()

	// default not delete resource permanently.
	// call option method 'WithPurge' to set enablePurge to true.
	db.enablePurge = false
	db.enableCache = false
	db.tableName = ""
	db.batchSize = 0
	db.noHook = false
	db.orQuery = false
	db.inTransaction = false
	db.shouldAutoMigrate = false
	db.tryRun = false
	db.db = DB.WithContext(context.Background())
}

func (db *database[M]) prepare() error {
	if db.db == nil || db.db == new(gorm.DB) {
		return ErrInvalidDB
	}
	if db.shouldAutoMigrate {
		if err := db.db.AutoMigrate(new(M)); err != nil {
			return err
		}
	}
	return nil
}

// WithDB returns a new database manipulator, only support *gorm.DB.
// eg: database.Database[*model.MeetingRoom]().WithDB(mysql.Software).WithTable("meeting_rooms").List(&rooms)
func (db *database[M]) WithDB(x any) types.Database[M] {
	var empty *gorm.DB
	if x == nil || x == new(gorm.DB) || x == empty {
		return db
	}
	// v := reflect.ValueOf(x)
	// if v.Kind() != reflect.Ptr {
	// 	return db
	// }
	// if v.IsNil() {
	// 	return db
	// }
	_db, ok := x.(*gorm.DB)
	if !ok {
		logger.Database.Warnw("invalid database type, expect *gorm.DB")
		return db
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	// db.shouldAutoMigrate = true
	if strings.ToLower(config.App.LoggerConfig.Level) == "debug" {
		db.db = _db.WithContext(context.TODO()).Debug().Limit(defaultLimit)
	} else {
		db.db = _db.WithContext(context.TODO()).Limit(defaultLimit)
	}
	return db
}

// WithTable multiple custom table, always used with the method `WithDB`.
// eg: database.Database[*model.MeetingRoom]().WithDB(mysql.Software).WithTable("meeting_rooms").List(&rooms)
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
	// db.db = db.db.Session(&gorm.Session{CreateBatchSize: db.batchSize})
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

// WithAnd with AND query condition(default).
// It must be called before WithQuery.
func (db *database[M]) WithAnd(flag ...bool) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.orQuery = false
	if len(flag) > 0 {
		db.orQuery = flag[0]
	}
	return db
}

// WithAnd with OR query condition.
// It must be called before WithQuery.
func (db *database[M]) WithOr(flag ...bool) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.orQuery = true
	if len(flag) > 0 {
		db.orQuery = flag[0]
	}
	return db
}

// WithIndex use specific index to query.
func (db *database[M]) WithIndex(index string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(strings.TrimSpace(index)) == 0 {
		return db
	}
	db.db = db.db.Clauses(hints.UseIndex(strings.TrimSpace(index)))
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

	// for i := 0; i < typ.NumField(); i++ {
	// 	// fmt.Println("---------------- for type ", typ.Field(i).Type, typ.Field(i).Type.Kind(), typ.Field(i).Name, val.Field(i).IsZero())
	// 	if val.Field(i).IsZero() {
	// 		continue
	// 	}
	//
	// 	switch typ.Field(i).Type.Kind() {
	// 	case reflect.Chan, reflect.Map, reflect.Func:
	// 		continue
	// 	case reflect.Struct:
	// 		// All `model.XXX` extends the basic model named `Base`,
	// 		if typ.Field(i).Name == "Base" {
	// 			if !val.Field(i).FieldByName("ID").IsZero() {
	// 				// Not overwrite the "ID" value set in types.Model.
	// 				// The "ID" value set in types.Model has higher priority than base model.
	// 				if _, loaded := q["id"]; !loaded {
	// 					q["id"] = val.Field(i).FieldByName("ID").Interface().(string)
	// 				}
	// 			}
	// 		} else {
	// 			structFieldToMap(typ.Field(i).Type, val.Field(i), q)
	// 		}
	// 		continue
	// 	}
	// 	// "json" tag priority is higher than typ.Field(i).Name
	// 	jsonTagStr := typ.Field(i).Tag.Get("json")
	// 	jsonTagItems := strings.Split(jsonTagStr, ",")
	// 	// NOTE: strings.Split always returns at least one element(empty string)
	// 	// We should not use len(jsonTagItems) to check the json tags exists.
	// 	jsonTag := ""
	// 	if len(jsonTagItems) == 0 {
	// 		// the structure lowercase field name as the query condition.
	// 		jsonTagItems[0] = typ.Field(i).Name
	// 	}
	// 	jsonTag = jsonTagItems[0]
	// 	if len(jsonTag) == 0 {
	// 		// the structure lowercase field name as the query condition.
	// 		jsonTag = typ.Field(i).Name
	// 	}
	// 	// "schema" tag have higher priority than "json" tag
	// 	schemaTagStr := strings.TrimSpace(typ.Field(i).Tag.Get("schema"))
	// 	if len(schemaTagStr) > 0 && schemaTagStr != jsonTagStr {
	// 		fmt.Println("-------------  schema tag overwrite json tag")
	// 		jsonTag = strings.TrimSpace(typ.Field(i).Tag.Get("schema"))
	// 	}
	//
	// 	v := val.Field(i).Interface()
	// 	var _v string
	// 	switch val.Field(i).Kind() {
	// 	case reflect.Bool:
	// 		// 由于 WHERE IN 语句会自动加上单引号,比如 WHERE `default` IN ('true')
	// 		// 但是我们想要的是 WHERE `default` IN (true),
	// 		// 所以没办法就只能直接转成 int 了.
	// 		_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
	// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 		_v = fmt.Sprintf("%d", v)
	// 	case reflect.Float32, reflect.Float64:
	// 		_v = fmt.Sprintf("%g", v)
	// 	case reflect.String:
	// 		_v = fmt.Sprintf("%s", v)
	// 	case reflect.Pointer:
	// 		v = val.Field(i).Elem().Interface()
	// 		// switch typ.Elem().Kind() {
	// 		switch val.Field(i).Elem().Kind() {
	// 		case reflect.Bool:
	// 			_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
	// 		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 			_v = fmt.Sprintf("%d", v)
	// 		case reflect.Float32, reflect.Float64:
	// 			_v = fmt.Sprintf("%g", v)
	// 		case reflect.String:
	// 			_v = fmt.Sprintf("%s", v)
	// 		case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
	// 		default:
	// 			_v = fmt.Sprintf("%v", v)
	// 		}
	// 	case reflect.Slice:
	// 		_len := val.Field(i).Len()
	// 		if _len == 0 {
	// 			logger.Database.Warn("reflect.Slice length is 0")
	// 			_len = 1
	// 		}
	// 		slice := reflect.MakeSlice(val.Field(i).Type(), _len, _len)
	// 		// fmt.Println("--------------- slice element", slice.Index(0), slice.Index(0).Kind(), slice.Index(0).Type())
	// 		switch slice.Index(0).Kind() {
	// 		case reflect.String: // handle string slice.
	// 			// WARN: val.Field(i).Type() is model.GormStrings not []string,
	// 			// execute statement `slice.Interface().([]string)` directly will case panic.
	// 			// _v = strings.Join(slice.Interface().([]string), ",") // the slice type is GormStrings not []string.
	// 			// We should make the slice of []string again.
	// 			slice = reflect.MakeSlice(reflect.TypeOf([]string{}), _len, _len)
	// 			reflect.Copy(slice, val.Field(i))
	// 			_v = strings.Join(slice.Interface().([]string), ",")
	// 		default:
	// 			_v = fmt.Sprintf("%v", v)
	// 		}
	// 	case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
	// 	default:
	// 		_v = fmt.Sprintf("%v", v)
	// 	}
	//
	// 	// json tag name naming format must be same as gorm table columns,
	// 	// both should be "snake case" or "camel case".
	// 	// gorm table columns naming format default to 'snake case',
	// 	// so the json tag name is converted to "snake case here"
	// 	// q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
	// 	q[strcase.SnakeCase(jsonTag)] = _v
	// }

	structFieldToMap(typ, val, q)
	// fmt.Println("------------- WithQuery", q)

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
			// skip the string slice which all element is empty.
			if len(strings.Join(items, "")) == 0 {
				continue
			}
			if len(items) > 1 { // If the query string has multiple value(seperated by ','), using regexp
				var regexpVal string
				for _, item := range items {
					// WARN: not forget to escape the regexp value using regexp.QuoteMeta.
					// eg: localhost\hello.world -> localhost\\hello\.world
					regexpVal = regexpVal + "|.*" + regexp.QuoteMeta(item) + ".*"
				}
				regexpVal = strings.TrimPrefix(regexpVal, "|")
				// db.db = db.db.Where(fmt.Sprintf("`%s` REGEXP ?", k), regexpVal)
				if db.orQuery {
					db.db = db.db.Or(fmt.Sprintf("`%s` REGEXP ?", k), regexpVal)
				} else {
					db.db = db.db.Where(fmt.Sprintf("`%s` REGEXP ?", k), regexpVal)
				}
			} else { // If the query string has only one value, using LIKE
				// db.db = db.db.Where(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
				if db.orQuery {
					db.db = db.db.Or(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
				} else {
					db.db = db.db.Where(fmt.Sprintf("`%s` LIKE ?", k), fmt.Sprintf("%%%v%%", v))
				}
			}
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
			// TODO: should we skip if items total length is 0?
			if len(strings.Join(items, "")) == 0 {
				continue
			}
			// db.db = db.db.Where(fmt.Sprintf("`%s` IN (?)", k), items)
			if db.orQuery {
				db.db = db.db.Or(fmt.Sprintf("`%s` IN (?)", k), items)
			} else {
				db.db = db.db.Where(fmt.Sprintf("`%s` IN (?)", k), items)
			}
		}
	}
	return db
}

// StructFieldToMap extracts the field tags from a struct and writes them into a map.
// This map can then be used to build SQL query conditions.
// FIXME: if the field type is boolean or ineger, disable the fuzzy matching.
func structFieldToMap(typ reflect.Type, val reflect.Value, q map[string]string) {
	if q == nil {
		q = make(map[string]string)
	}
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).IsZero() {
			continue
		}

		switch typ.Field(i).Type.Kind() {
		case reflect.Chan, reflect.Map, reflect.Func:
			continue
		case reflect.Struct:
			// All `model.XXX` extends the basic model named `Base`,
			if typ.Field(i).Name == "Base" {
				if !val.Field(i).FieldByName("CreatedBy").IsZero() {
					// Not overwrite the "CreatedBy" value set in types.Model.
					// The "CreatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["created_by"]; !loaded {
						q["created_by"] = val.Field(i).FieldByName("CreatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("UpdatedBy").IsZero() {
					// Not overwrite the "UpdatedBy" value set in types.Model.
					// The "UpdatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["updated_by"]; !loaded {
						q["updated_by"] = val.Field(i).FieldByName("UpdatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("ID").IsZero() {
					// Not overwrite the "ID" value set in types.Model.
					// The "ID" value set in types.Model has higher priority than base model.
					if _, loaded := q["id"]; !loaded {
						q["id"] = val.Field(i).FieldByName("ID").Interface().(string)
					}
				}
			} else {
				structFieldToMap(typ.Field(i).Type, val.Field(i), q)
			}
			continue
		}
		// "json" tag priority is higher than typ.Field(i).Name
		jsonTagStr := strings.TrimSpace(typ.Field(i).Tag.Get("json"))
		jsonTagItems := strings.Split(jsonTagStr, ",")
		// NOTE: strings.Split always returns at least one element(empty string)
		// We should not use len(jsonTagItems) to check the json tags exists.
		jsonTag := ""
		if len(jsonTagItems) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTagItems[0] = typ.Field(i).Name
		}
		jsonTag = jsonTagItems[0]
		if len(jsonTag) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTag = typ.Field(i).Name
		}
		// "schema" tag have higher priority than "json" tag
		schemaTagStr := strings.TrimSpace(typ.Field(i).Tag.Get("schema"))
		schemaTagItems := strings.Split(schemaTagStr, ",")
		schemaTag := ""
		if len(schemaTagItems) > 0 {
			schemaTag = schemaTagItems[0]
		}
		if len(schemaTag) > 0 && schemaTag != jsonTag {
			fmt.Printf("------------------ json tag replace by schema tag: %s -> %s\n", jsonTag, schemaTag)
			jsonTag = schemaTag
		}

		v := val.Field(i).Interface()
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
			case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
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
		case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
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

// WithTimeRange
func (db *database[M]) WithTimeRange(columnName string, startTime time.Time, endTime time.Time) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(columnName) == 0 || startTime.IsZero() || endTime.IsZero() {
		return db
	}
	db.db = db.db.Where(fmt.Sprintf("`%s` BETWEEN ? AND ?", columnName), startTime, endTime)
	return db
}

// WithSelect specify fields that you want when querying, creating, updating
// default select all fields.
// WARNING: Using WithSelect may result in the removal of certain fields from table records
// if there are multiple hooks in the service and model layers. Use with caution.
func (db *database[M]) WithSelect(columns ...string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	if len(columns) == 0 {
		return db
	}
	_columns := make([]string, 0)
	for i := range columns {
		col := strings.TrimSpace(columns[i])
		if len(col) > 0 && !contains(defaultsColumns, col) {
			_columns = append(_columns, col)
		}
	}
	if len(_columns) == 0 {
		return db
	}
	db.db = db.db.Select(append(_columns, defaultsColumns...))
	return db
}

// WithSelectRaw
func (db *database[M]) WithSelectRaw(query any, args ...any) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Select(query, args...)
	return db
}

// WithTransaction executes operations within a transaction.
// It's typically used with DB.Transaction():
// NOTE:
// 1. If tx is provided, disable GORM's default transaction
// 2. If tx is nil, use default behavior
//
// Example:
//
//	err := DB.Transaction(func(tx *gorm.DB) error {
//	    return Database[*Order]().
//	        WithTransaction(tx).
//	        WithLock().
//	        Get(&order, orderID)
//	})
func (db *database[M]) WithTransaction(tx any) types.Database[M] {
	var empty *gorm.DB
	if tx == nil || tx == new(gorm.DB) || tx == empty {
		return db
	}
	// v := reflect.ValueOf(x)
	// if v.Kind() != reflect.Ptr {
	// 	return db
	// }
	// if v.IsNil() {
	// 	return db
	// }
	_tx, ok := tx.(*gorm.DB)
	if !ok {
		logger.Database.Warnw("invalid database type, expect *gorm.DB")
		return db
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.inTransaction = true
	db.db = _tx.Session(&gorm.Session{SkipDefaultTransaction: true})
	return db
}

// WithLock adds locking clause to SELECT statement.
// It must be used within a transaction (WithTransaction).
//
// Lock modes:
//   - "UPDATE" (default): SELECT ... FOR UPDATE
//   - "SHARE": SELECT ... FOR SHARE
//   - "UPDATE_NOWAIT": SELECT ... FOR UPDATE NOWAIT
//   - "SHARE_NOWAIT": SELECT ... FOR SHARE NOWAIT
//   - "UPDATE_SKIP_LOCKED": SELECT ... FOR UPDATE SKIP LOCKED
//   - "SHARE_SKIP_LOCKED": SELECT ... FOR SHARE SKIP LOCKED
//
// Example:
//
//	DB.Transaction(func(tx *gorm.DB) error {
//	    // Default FOR UPDATE lock
//	    err := Database[*Order]().
//	        WithTx(tx).
//	        WithLock().
//	        Get(&order, orderID)
//
//	    // FOR UPDATE NOWAIT
//	    err = Database[*Order]().
//	        WithTx(tx).
//	        WithLock("UPDATE_NOWAIT").
//	        Get(&order, orderID)
//	})
func (db *database[M]) WithLock(mode ...string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.inTransaction {
		logger.Database.Warnw("WithLock must be used within a transaction")
		return db
	}

	strength := "UPDATE"
	options := ""
	if len(mode) > 0 {
		switch strings.ToUpper(mode[0]) {
		case "SHARE":
			strength = "SHARE"
		case "UPDATE_NOWAIT":
			strength = "UPDATE"
			options = "NOWAIT"
		case "SHARE_NOWAIT":
			strength = "SHARE"
			options = "NOWAIT"
		case "UPDATE_SKIP_LOCKED":
			strength = "UPDATE"
			options = "SKIP LOCKED"
		case "SHARE_SKIP_LOCKED":
			strength = "SHARE"
			options = "SKIP LOCKED"
		}
	}

	db.db = db.db.Clauses(clause.Locking{
		Strength: strength,
		Options:  options,
	})
	return db
}

// WithJoinRaw adds JOIN clause to query.
//
// Basic Join:
//
//	db.WithJoinRaw("JOIN users ON users.id = orders.user_id")
//
// Left Join with conditions:
//
//	db.WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id AND users.active = ?", 1)
//
// Multiple Joins:
//
//	db.WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
//	    WithJoinRaw("LEFT JOIN products ON products.id = orders.product_id")
//
// Join with Select:
//
//	db.WithSelectRaw("orders.*, users.name").
//	    WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id")
//
// Complex Examples:
//
// 1. Query order with user info:
//
//	type Order struct {
//	    ID     string `gorm:"primarykey"`
//	    UserID string
//	    Amount float64
//	    User   User   `gorm:"foreignKey:UserID"`
//	}
//
//	type User struct {
//	    ID    string `gorm:"primarykey"`
//	    Name  string
//	    Email string
//	}
//
//	var orders []Order
//	err := Database[*Order]().
//	    WithSelectRaw("orders.*, users.name as user_name").
//	    WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
//	    List(&orders)
//
// 2. Multi-table join query:
//
//	var details []OrderDetail
//	err := Database[*OrderDetail]().
//	    WithSelectRaw("order_details.*, orders.amount, products.name as product_name").
//	    WithJoinRaw("LEFT JOIN orders ON orders.id = order_details.order_id").
//	    WithJoinRaw("LEFT JOIN products ON products.id = order_details.product_id").
//	    List(&details)
//
// 3. Query orders with active users:
//
//	var orders []Order
//	err := Database[*Order]().
//	    WithSelectRaw("orders.*, users.name").
//	    WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id AND users.active = ?", 1).
//	    List(&orders)
//
// 4. Complex query with multiple conditions:
//
//	var orders []Order
//	err := Database[*Order]().
//	    WithSelectRaw("orders.*, users.name, products.name as product_name").
//	    WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
//	    WithJoinRaw("LEFT JOIN order_details ON order_details.order_id = orders.id").
//	    WithJoinRaw("LEFT JOIN products ON products.id = order_details.product_id").
//	    WithTimeRange("orders.created_at", startTime, endTime).
//	    WithOrder("orders.created_at DESC").
//	    WithScope(page, size).
//	    List(&orders)
func (db *database[M]) WithJoinRaw(query string, args ...any) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()

	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return db
	}

	upperQuery := strings.ToUpper(query)
	if !strings.Contains(upperQuery, "JOIN") || !strings.Contains(upperQuery, "ON") {
		logger.Database.Warn("invalid join clause, must contain JOIN and ON",
			"query", query,
			"table", reflect.TypeOf(*new(M)).Elem().Name(),
		)
		return db
	}

	db.db = db.db.Joins(query, args...)
	return db
}

// WithGroup adds GROUP BY clause to SELECT statement.
// For example:
//
//	// Basic group by
//	db.WithGroup("user_id")
//
//	// Group by multiple columns
//	db.WithGroup("user_id, order_status")
//
//	// Common usage with aggregate functions
//	db.WithSelectRaw("user_id, COUNT(*) as order_count, SUM(amount) as total_amount").
//	   WithGroup("user_id")
//
// Note: WithGroup is typically used with aggregate functions (COUNT, SUM, AVG, etc.)
// and should be combined with WithSelectRaw to specify the grouped fields.
func (db *database[M]) WithGroup(name string) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	name = strings.TrimSpace(name)
	if len(name) > 0 {
		db.db = db.db.Group(name)
	}
	return db
}

// WithHaving adds HAVING clause to filter grouped records.
// HAVING clause is used to filter groups, similar to WHERE but operates on grouped records.
// For example:
//
//	// Basic having clause
//	db.WithHaving("COUNT(*) > ?", 5)
//
//	// With aggregate functions
//	db.WithSelectRaw("user_id, COUNT(*) as order_count, SUM(amount) as total_amount").
//	   WithGroup("user_id").
//	   WithHaving("SUM(amount) > ?", 1000)
//
//	// Multiple conditions
//	db.WithHaving("COUNT(*) > ? AND SUM(amount) > ?", 5, 1000)
//
// Note: WithHaving must be used with GROUP BY clause and aggregate functions.
func (db *database[M]) WithHaving(query any, args ...any) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db = db.db.Having(query, args...)
	return db
}

// WithOrder
// reference: https://www.cnblogs.com/Braveliu/p/10654091.html
// For example:
// - WithOrder("name") // default ASC.
// - WithOrder("name desc")
// - WithOrder("created_at")
// - WithOrder("updated_at desc")
// multiple keyw order, eg:
// - "field1, field2 desc, field3 asc"
// - "created_at desc, id desc"
// NOTE: you can using the mysql keyword, such as: "order", "limit".
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

// WithTryRun only executes model hooks without performing actual database operations.
// Also logs the SQL statements that would have been executed.
func (db *database[M]) WithTryRun(enable ...bool) types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.tryRun = true
	if len(enable) > 0 {
		db.tryRun = enable[0]
	}
	return db
}

// WithoutHook will disable all model hooks.
func (db *database[M]) WithoutHook() types.Database[M] {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.noHook = true
	return db
}

// Create one or multiple objects in database.
// It will update the "created_at" and "updated_at" field.
// Examples:
// - database.XXX().Create(&model.XXX{ID: id, Field1: field1, Field2: field2})
func (db *database[M]) Create(objs ...M) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Create")
	defer done(err)
	if len(objs) == 0 {
		return nil
	}

	if db.enableCache {
		defer lru.Cache[[]M]().Flush()
	}
	// if config.App.RedisConfig.Enable {
	// 	defer func() {
	// 		go func() {
	// 			begin := time.Now()
	// 			prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "create")
	// 			defer logger.Cache.Infow("remove cache after create", "cost", time.Since(begin).String(), "prefix", prefix)
	// 			if err = redis.RemovePrefix(prefix); err != nil {
	// 				logger.Cache.Errorw("failed to remove cache keys", err, "action", "create")
	// 			}
	// 		}()
	// 	}()
	// }

	var empty M // call nil value M will cause panic.
	// Invoke model hook: CreateBefore.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].CreateBefore(); err != nil {
				return err
			}
		}
	}
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) {
			objs[i].SetID() // set id when id is empty.
		}
	}

	// if err = db.db.Save(objs).Error; err != nil {
	// if err = db.db.Table(db.tableName).Save(objs).Error; err != nil {
	// 	return err
	// }
	//
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	batchSize := defaultBatchSize
	if db.batchSize > 0 {
		batchSize = db.batchSize
	}
	for i := 0; i < len(objs); i += batchSize {
		end := i + batchSize
		if end > len(objs) {
			end = len(objs)
		}
		if err = db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Save(objs[i:end]).Error; err != nil {
			return err
		}
	}

	// because db.db.Delete method just update field "delete_at" to current time,
	// not really delete it(soft delete).
	// If record already exists, Update method update all fields but exlcude "created_at" by
	// mysql "ON DUPLICATE KEY UPDATE" mechanism. so we should update the "created_at" field manually.
	for i := range objs {
		// 这里要重新创建一个 gorm.DB 实例, 否则会出现这种语句, id 出现多次了.
		// UPDATE `assets` SET `created_at`='2023-11-12 14:35:42.604',`updated_at`='2023-11-12 14:35:42.604' WHERE id = '010103NU000020' AND `assets`.`deleted_at` IS NULL AND id = '010103NU000021' AND id = '010103NU000022' LIMIT 1000
		var _db *gorm.DB
		if strings.ToLower(config.App.LoggerConfig.Level) == "debug" {
			_db = DB.Debug()
		} else {
			_db = DB
		}
		createdAt := time.Now()
		// if err = _db.Model(*new(M)).Where("id = ?", objs[i].GetID()).Update("created_at", time.Now()).Error; err != nil {
		if err = _db.Table(tableName).Model(*new(M)).Where("id = ?", objs[i].GetID()).Update("created_at", createdAt).Error; err != nil {
			return err
		}
		objs[i].SetCreatedAt(createdAt)
		if db.enableCache {
			lru.Cache[M]().Remove(objs[i].GetID())
		}

	}
	// Invoke model hook: CreateAfter.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].CreateAfter(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete one or multiple objects in database.
// Examples:
// - database.XXX().Delete(&model.XXX{ID: id}) // delete record with primary key
// - database.XXX().WithQuery(req).Delete(req) // delete record with where condiions.
// FIXME: is delete should use defaultLimit or defaultBatchSize?
func (db *database[M]) Delete(objs ...M) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Delete")
	defer done(err)
	if len(objs) == 0 {
		return nil
	}

	if db.enableCache {
		defer lru.Cache[[]M]().Flush()
	}
	// if config.App.RedisConfig.Enable {
	// 	defer func() {
	// 		// TODO:only delete cache of all list statement and cache for current get statements.
	// 		go func() {
	// 			begin := time.Now()
	// 			prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "delete")
	// 			defer logger.Cache.Infow("remove cache after delete", "cost", time.Since(begin).String(), "prefix", prefix)
	// 			if err = redis.RemovePrefix(prefix); err != nil {
	// 				logger.Cache.Errorw("failed to remove cache keys", err, "action", "delete")
	// 			}
	// 		}()
	// 	}()
	// }

	var empty M // call nil value M will cause panic.
	// Invoke model hook: DeleteBefore.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].DeleteBefore(); err != nil {
				return err
			}
		}
	}
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if db.enablePurge {
		// delete permanently.
		// if err = db.db.Unscoped().Delete(objs).Error; err != nil {
		// if err = db.db.Table(db.tableName).Unscoped().Delete(objs).Error; err != nil {
		// 	return err
		// }
		//
		batchSize := defaultDeleteBatchSize
		if db.batchSize > 0 {
			batchSize = db.batchSize
		}
		for i := 0; i < len(objs); i += batchSize {
			end := i + batchSize
			if end > len(objs) {
				end = len(objs)
			}
			if err = db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Unscoped().Delete(objs[i:end]).Error; err != nil {
				return err
			}
			if db.enableCache {
				lru.Cache[M]().Remove(objs[i].GetID())
			}
		}
	} else {
		// Delete() method just update field "delete_at" to currrent time.
		// DO NOT FORGET update the "created_at" field when create/update if record already exists.
		// if err = db.db.Delete(objs).Error; err != nil {
		// if err = db.db.Table(db.tableName).Delete(objs).Error; err != nil {
		// 	return err
		// }
		//
		batchSize := defaultDeleteBatchSize
		if db.batchSize > 0 {
			batchSize = db.batchSize
		}
		for i := 0; i < len(objs); i += batchSize {
			end := i + batchSize
			if end > len(objs) {
				end = len(objs)
			}
			if err = db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Delete(objs[i:end]).Error; err != nil {
				return err
			}
			if db.enableCache {
				lru.Cache[M]().Remove(objs[i].GetID())
			}
		}
	}
	// Invoke model hook: DeleteAfter.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].DeleteAfter(); err != nil {
				return err
			}
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
func (db *database[M]) Update(objs ...M) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Update")
	defer done(err)
	if len(objs) == 0 {
		return nil
	}

	if db.enableCache {
		defer lru.Cache[[]M]().Flush()
	}
	// if config.App.RedisConfig.Enable {
	// 	defer func() {
	// 		go func() {
	// 			begin := time.Now()
	// 			prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "update")
	// 			defer logger.Cache.Infow("remove cache after update", "cost", time.Since(begin).String(), "prefix", prefix)
	// 			if err = redis.RemovePrefix(prefix); err != nil {
	// 				logger.Cache.Errorw("failed to remove cache keys", err, "action", "update")
	// 			}
	// 		}()
	// 	}()
	// }

	var empty M // call nil value M will cause panic.
	// Invoke model hook: UpdateBefore.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].UpdateBefore(); err != nil {
				return err
			}
		}
	}
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) {
			objs[i].SetID() // set id when id is empty
		}
	}
	// if err = db.db.Save(objs).Error; err != nil {
	// if err = db.db.Table(db.tableName).Save(objs).Error; err != nil {
	// 	return err
	// }
	//
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	batchSize := defaultBatchSize
	if db.batchSize > 0 {
		batchSize = db.batchSize
	}
	for i := 0; i < len(objs); i += batchSize {
		end := i + batchSize
		if end > len(objs) {
			end = len(objs)
		}
		if err = db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Save(objs[i:end]).Error; err != nil {
			zap.S().Error(err)
			return err
		}
		if db.enableCache {
			for j := range objs[i:end] {
				lru.Cache[M]().Remove(objs[j].GetID())
			}
		}
	}
	// Invoke model hook: UpdateAfter.
	for i := range objs {
		if !reflect.DeepEqual(empty, objs[i]) && !db.noHook {
			if err = objs[i].UpdateAfter(); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateById only update one record with specific id.
// its not invoke model hook.
func (db *database[M]) UpdateById(id string, key string, val any) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("UpdateById")
	defer done(err)

	if db.enableCache {
		defer lru.Cache[[]M]().Flush()
	}
	// if config.App.RedisConfig.Enable {
	// 	defer func() {
	// 		go func() {
	// 			begin := time.Now()
	// 			prefix, _ := buildCacheKey(db.db.Model(*new(M)).Session(&gorm.Session{DryRun: true}).Statement, "update_by_id")
	// 			defer logger.Cache.Infow("remove cache after update_by_id", "cost", time.Since(begin).String(), "prefix", prefix)
	// 			if err = redis.RemovePrefix(prefix); err != nil {
	// 				logger.Cache.Errorw("failed to remove cache keys", err, "action", "update")
	// 			}
	// 		}()
	// 	}()
	// }

	// return db.db.Model(*new(M)).Where("id = ?", id).Update(key, val).Error
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Model(*new(M)).Where("id = ?", id).Update(key, val).Error; err != nil {
		return err
	}
	if db.enableCache {
		lru.Cache[M]().Remove(id)
	}
	return nil
}

// List find all record if not run WithQuery(query).
// List find record with `where` conditioan if run WithQuery(query).
// Examples:
//   - data := make([]*model.XXX, 0)
//     database.XXX().WithScope(page, size).WithQuery(&model.XXX{Field1: field1, Field2: field2}).List(&data)
//   - data := make([]*model.XXX, 0)
//     database.XXX().List(&data)
func (db *database[M]) List(dest *[]M, _cache ...*[]byte) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("List")
	defer done(err)
	if dest == nil {
		return nil
	}

	begin := time.Now()
	var key string
	var _dest []M
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Find(dest).Statement, "list")
	if _dest, exists = lru.Cache[[]M]().Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("list", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("list", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		*dest = _dest
		logger.Cache.Infow("list from cache", "cost", time.Since(begin).String(), "key", key)
		return nil
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	// var _cache []byte
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Find(dest).Statement, "list")
	// if len(cache) == 0 {
	// 	shouldDecode = true
	// } else {
	// 	if cache[0] == nil {
	// 		shouldDecode = true
	// 	}
	// }
	// if shouldDecode {
	// 	var _dest []M
	// 	if _dest, err = redis.GetML[M](key); err == nil {
	// 		val := reflect.ValueOf(dest)
	// 		if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		if !val.Elem().CanSet() {
	// 			return ErrNotSetSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_dest))
	// 		logger.Cache.Infow("list and decode from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// } else {
	// 	if _cache, err = redis.Get(key); err == nil {
	// 		val := reflect.ValueOf(cache[0])
	// 		if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		if !val.Elem().CanSet() {
	// 			return ErrNotSetSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_cache))
	// 		logger.Cache.Infow("list from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// }
	// if !errors.Is(err, redis.ErrKeyNotExists) {
	// 	logger.Cache.Error(err)
	// 	return err
	// }
	// // Not Found cache and continue.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	var empty M // call nil value M will cause panic.
	// Invoke model hook: ListBefore.
	for i := range *dest {
		if !reflect.DeepEqual(empty, (*dest)[i]) && !db.noHook {
			if err = (*dest)[i].ListBefore(); err != nil {
				return err
			}
		}
	}
	// if err = db.db.Find(dest).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).Find(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: ListAfter()
	for i := range *dest {
		if !reflect.DeepEqual(empty, (*dest)[i]) && !db.noHook {
			if err = (*dest)[i].ListAfter(); err != nil {
				return err
			}
		}
	}
	// Cache the result.
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("list from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.SetML[M](key, *dest); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()
	// }
	if db.enableCache {
		logger.Cache.Infow("list from database", "cost", time.Since(begin).String(), "key", key)
		lru.Cache[[]M]().Set(key, *dest)
	}

	return nil
}

// // Find equal to WithQuery(condition).List()
// // More detail see `List` document.
// func (db *database[T]) Find(dest *[]T, query T) error {
// 	return db.db.Where(query).Find(dest).Error
// }

// Get find one record accoding to `id`.
func (db *database[M]) Get(dest M, id string, _cache ...*[]byte) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Get")
	defer done(err)

	begin := time.Now()
	var key string
	var _dest M
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Find(dest).Statement, "get", id)
	if _dest, exists = lru.Cache[M]().Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("get", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("get", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		val := reflect.ValueOf(dest)
		if val.Kind() != reflect.Ptr {
			return ErrNotPtrStruct
		}
		if !val.Elem().CanAddr() {
			return ErrNotAddressableModel
		}
		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
		logger.Cache.Infow("get from cache", "cost", time.Since(begin).String(), "key", key)
		return nil
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// // _, key = BuildKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Find(dest).Statement, "get")
	// // 我发现这个 id 的值怎么都无法填充到 sql 语句中, 所以直接用 id 作为 key 的一部分,而不是用 sql 语句
	// // 如果不用 id 作为 redis key, 那么多个 get 的语句相同则 key 相同,获取到的都是同一个缓存.
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Find(dest).Statement, "get", id)
	// if len(cache) == 0 {
	// 	shouldDecode = true
	// } else {
	// 	if cache[0] == nil {
	// 		shouldDecode = true
	// 	}
	// }
	// if shouldDecode { // query cache from redis and decoded into 'dest'.
	// 	var _dest M
	// 	if _dest, err = redis.GetM[M](key); err == nil {
	// 		val := reflect.ValueOf(dest)
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrStruct
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableModel
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
	// 		logger.Cache.Infow("get and decode from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// } else {
	// 	var _cache []byte
	// 	if _cache, err = redis.Get(key); err == nil {
	// 		val := reflect.ValueOf(cache[0])
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_cache))
	// 		logger.Cache.Infow("get from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// }
	// if err != redis.ErrKeyNotExists {
	// 	logger.Cache.Error(err)
	// 	return err
	// }
	// // Not Found cache, continue query database.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	var empty M // call nil value M will cause panic.
	// Invoke model hook: GetBefore.
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetBefore(); err != nil {
			return err
		}
	}
	// if err = db.db.Where("id = ?", id).Find(dest).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).Where("id = ?", id).Find(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter.
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetAfter(); err != nil {
			return err
		}
	}
	// // Cache the result.
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("get from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.SetM[M](key, dest); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()
	// }
	if db.enableCache {
		logger.Cache.Infow("get from database", "cost", time.Since(begin).String(), "key", key)
		lru.Cache[M]().Set(key, dest)
	}
	return nil
}

// Count returns the total number of records with the given query condition.
// NOTE: The underlying type msut be pointer to struct, otherwise panic will occur.
func (db *database[M]) Count(count *int64) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Count")
	defer done(err)

	begin := time.Now()
	var key string
	var _cache int64
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Model(*new(M)).Count(count).Statement, "count")
	if _cache, exists = lru.Int64.Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("count", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("count", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		*count = _cache
		logger.Cache.Infow("count from cache", "cost", time.Since(begin).String(), "key", key)
		return
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var _cache int64
	// if count == nil {
	// 	return nil
	// }
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// fmt.Println("---- buildCacheKey count")
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).Model(*new(M)).Count(count).Statement, "count")
	// if _cache, err = redis.GetInt(key); err == nil {
	// 	*count = _cache
	// 	logger.Cache.Infow("count from cache", "cost", time.Since(begin).String(), "key", key)
	// 	return
	// }
	// if !errors.Is(err, redis.ErrKeyNotExists) {
	// 	logger.Cache.Error(err)
	// 	return err
	// }
	// // NOT FOUND cache, continue query.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	// if err = db.db.Model(*new(M)).Count(count).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).Model(*new(M)).Limit(-1).Count(count).Error; err != nil {
		logger.Cache.Error(err)
		return err
	}
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("count from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.Set(key, *count); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()
	//
	// }
	if db.enableCache {
		logger.Cache.Infow("count from database", "cost", time.Since(begin).String(), "key", key)
		lru.Int64.Set(key, *count)

	}
	return nil
}

// First finds the first record ordered by primary key.
func (db *database[M]) First(dest M, _cache ...*[]byte) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("First")
	defer done(err)

	begin := time.Now()
	var key string
	var _dest M
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "first")
	if _dest, exists = lru.Cache[M]().Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("first", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("first", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		val := reflect.ValueOf(dest)
		if val.Kind() != reflect.Ptr {
			return ErrNotPtrStruct
		}
		if !val.Elem().CanAddr() {
			return ErrNotAddressableModel
		}
		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
		logger.Cache.Infow("first from cache", "cost", time.Since(begin).String(), "key", key)
		return nil // Found cache and return.
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "first")
	// if len(cache) == 0 {
	// 	shouldDecode = true
	// } else {
	// 	if cache[0] == nil {
	// 		shouldDecode = true
	// 	}
	// }
	// if shouldDecode { // query cache from redis and decode into 'dest'.
	// 	var _dest M
	// 	if _dest, err = redis.GetM[M](key); err == nil {
	// 		val := reflect.ValueOf(dest)
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrStruct
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableModel
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
	// 		logger.Cache.Infow("first and decode from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// } else {
	// 	var _cache []byte
	// 	if _cache, err = redis.Get(key); err == nil {
	// 		val := reflect.ValueOf(cache[0])
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_cache))
	// 		logger.Cache.Infow("first from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// 	if err != redis.ErrKeyNotExists {
	// 		logger.Cache.Error(err)
	// 		return err
	// 	}
	// }
	// Not Found cache, continue query database.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	var empty M // call nil value M will cause panic.
	// Invoke model hook: GetBefore
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetBefore(); err != nil {
			return err
		}
	}
	// if err = db.db.First(dest).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).First(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetAfter(); err != nil {
			return err
		}
	}
	// // Cache the result.
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("first from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.SetM[M](key, dest); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()
	// }
	if db.enableCache {
		logger.Cache.Infow("first from database", "cost", time.Since(begin).String(), "key", key)
		lru.Cache[M]().Set(key, dest)
	}
	return nil
}

// Last finds the last record ordered by primary key.
func (db *database[M]) Last(dest M, _cache ...*[]byte) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Last")
	defer done(err)

	begin := time.Now()
	var key string
	var _dest M
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "last")
	if _dest, exists = lru.Cache[M]().Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("last", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("last", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		val := reflect.ValueOf(dest)
		if val.Kind() != reflect.Ptr {
			return ErrNotPtrStruct
		}
		if !val.Elem().CanAddr() {
			return ErrNotAddressableModel
		}
		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
		logger.Cache.Infow("last from cache", "cost", time.Since(begin).String(), "key", key)
		return nil // Found cache and return.
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "last")
	// if len(cache) == 0 {
	// 	shouldDecode = true
	// } else {
	// 	if cache[0] == nil {
	// 		shouldDecode = true
	// 	}
	// }
	// if shouldDecode { // query cache from redis and decode into 'dest'.
	// 	var _dest M
	// 	if _dest, err = redis.GetM[M](key); err == nil {
	// 		val := reflect.ValueOf(dest)
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrStruct
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableModel
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
	// 		logger.Cache.Infow("last and decode from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// } else {
	// 	var _cache []byte
	// 	if _cache, err = redis.Get(key); err == nil {
	// 		val := reflect.ValueOf(cache[0])
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_cache))
	// 		logger.Cache.Infow("last from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// }
	// if err != redis.ErrKeyNotExists {
	// 	logger.Cache.Error(err)
	// 	return err
	// }
	// // Not Found cache, continue query database.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	var empty M // call nil value M will cause panic.
	// Invoke model hook: GetBefore.
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetBefore(); err != nil {
			return err
		}
	}
	// if err = db.db.Last(dest).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).Last(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetAfter(); err != nil {
			return err
		}
	}
	// // Cache the result.
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("last from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.SetM[M](key, dest); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()
	// }
	if db.enableCache {
		logger.Cache.Infow("last from database", "cost", time.Since(begin).String(), "key", key)
		lru.Cache[M]().Set(key, dest)
	}
	return nil
}

// Take finds the first record returned by the database in no specified order.
func (db *database[M]) Take(dest M, _cache ...*[]byte) (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Take")
	defer done(err)

	begin := time.Now()
	var key string
	var _dest M
	var exists bool
	if !db.enableCache {
		goto QUERY
	}
	_, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "take")
	if _dest, exists = lru.Cache[M]().Get(key); !exists {
		// metrics.CacheMiss.WithLabelValues("take", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		goto QUERY
	} else {
		// metrics.CacheHit.WithLabelValues("take", reflect.TypeOf(*new(M)).Elem().Name()).Inc()
		val := reflect.ValueOf(dest)
		if val.Kind() != reflect.Ptr {
			return ErrNotPtrStruct
		}
		if !val.Elem().CanAddr() {
			return ErrNotAddressableModel
		}
		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
		logger.Cache.Infow("take from cache", "cost", time.Since(begin).String(), "key", key)
		return nil // Found cache and return.
	}

	// =============================
	// ===== BEGIN redis cache =====
	// =============================
	// begin := time.Now()
	// var key string
	// var shouldDecode bool // if cache is nil or cache[0] is nil, we should decod the queryed cache in to "dest".
	// if !db.enableCache {
	// 	goto QUERY
	// }
	// if !config.App.RedisConfig.Enable {
	// 	goto QUERY
	// }
	// _, key = buildCacheKey(db.db.Session(&gorm.Session{DryRun: true}).First(dest).Statement, "take")
	// if len(cache) == 0 {
	// 	shouldDecode = true
	// } else {
	// 	if cache[0] == nil {
	// 		shouldDecode = true
	// 	}
	// }
	// if shouldDecode { // query cache from redis and decode into 'dest'.
	// 	var _dest M
	// 	if _dest, err = redis.GetM[M](key); err == nil {
	// 		val := reflect.ValueOf(dest)
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrStruct
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableModel
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_dest).Elem()) // the type of M is pointer to struct.
	// 		logger.Cache.Infow("get and decode from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// } else {
	// 	var _cache []byte
	// 	if _cache, err = redis.Get(key); err == nil {
	// 		val := reflect.ValueOf(cache[0])
	// 		if val.Kind() != reflect.Ptr {
	// 			return ErrNotPtrSlice
	// 		}
	// 		if !val.Elem().CanAddr() {
	// 			return ErrNotAddressableSlice
	// 		}
	// 		val.Elem().Set(reflect.ValueOf(_cache))
	// 		logger.Cache.Infow("take from cache", "cost", time.Since(begin).String(), "key", key)
	// 		return nil // Found cache and return.
	// 	}
	// }
	// if err != redis.ErrKeyNotExists {
	// 	logger.Cache.Error(err)
	// 	return err
	// }
	// // Not Found cache, continue query database.
	// ===========================
	// ===== END redis cache =====
	// ===========================

QUERY:
	var empty M // call nil value M will cause panic.
	// Invoke model hook: GetBefore.
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetBefore(); err != nil {
			return err
		}
	}
	// if err = db.db.Take(dest).Error; err != nil {
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	if err = db.db.Table(tableName).Take(dest).Error; err != nil {
		return err
	}
	// Invoke model hook: GetAfter.
	if !reflect.DeepEqual(empty, dest) && !db.noHook {
		if err = dest.GetAfter(); err != nil {
			return err
		}
	}
	// // Cache the result.
	// if db.enableCache && config.App.RedisConfig.Enable {
	// 	logger.Cache.Infow("take from database", "cost", time.Since(begin).String(), "key", key)
	// 	go func() {
	// 		if err = redis.SetM[M](key, dest); err != nil {
	// 			logger.Cache.Error(err)
	// 		}
	// 	}()

	//
	// }
	if db.enableCache {
		logger.Cache.Infow("take from database", "cost", time.Since(begin).String(), "key", key)
		lru.Cache[M]().Set(key, dest)
	}
	return nil
}

// Cleanup delete all records that column 'deleted_at' is not null.
func (db *database[M]) Cleanup() (err error) {
	if err = db.prepare(); err != nil {
		return err
	}
	defer db.reset()
	done := db.trace("Cleanup")
	defer done(err)

	// return db.db.Limit(-1).Where("deleted_at IS NOT NULL").Model(*new(M)).Unscoped().Delete(make([]M, 0)).Error
	typ := reflect.TypeOf(*new(M)).Elem()
	tableName := reflect.New(typ).Interface().(M).GetTableName()
	if len(db.tableName) > 0 {
		tableName = db.tableName
	}
	return db.db.Session(&gorm.Session{DryRun: db.tryRun}).Table(tableName).Limit(-1).Where("deleted_at IS NOT NULL").Model(*new(M)).Unscoped().Delete(make([]M, 0)).Error
}

// Health checks the database connectivity and basic operations.
// It returns nil if the database is healthy, otherwise returns an error.
func (db *database[M]) Health() error {
	if err := db.prepare(); err != nil {
		return err
	}
	defer db.reset()

	begin := time.Now()

	// 1.check database connection
	if err := db.db.Exec("SELECT 1").Error; err != nil {
		logger.Database.Errorw("database connection check failed",
			"error", err,
			"cost", time.Since(begin).String(),
		)
		return fmt.Errorf("database connection check failed: %w", err)
	}

	// 2.check database connection pool
	sqlDB, err := db.db.DB()
	if err != nil {
		logger.Database.Errorw("get sql.DB instance failed",
			"error", err,
			"cost", time.Since(begin).String(),
		)
		return fmt.Errorf("get sql.DB instance failed: %w", err)
	}

	// check database connection pool config
	stats := sqlDB.Stats()
	if stats.OpenConnections >= stats.MaxOpenConnections {
		logger.Database.Warnw("database connection pool is full",
			"open", stats.OpenConnections,
			"max", stats.MaxOpenConnections,
			"in_use", stats.InUse,
			"idle", stats.Idle,
			"cost", time.Since(begin).String(),
		)
	}

	// 3.check database response time
	if err := sqlDB.PingContext(db.ctx); err != nil {
		logger.Database.Errorw("database ping failed",
			"error", err,
			"cost", time.Since(begin).String(),
		)
		return fmt.Errorf("database ping failed: %w", err)
	}

	logger.Database.Infow("database health check passed",
		"open_connections", stats.OpenConnections,
		"in_use_connections", stats.InUse,
		"idle_connections", stats.Idle,
		"max_open_connections", stats.MaxOpenConnections,
		"cost", time.Since(begin).String(),
	)

	return nil
}

// Database implement interface types.Database, its default database manipulator.
// Database[M] returns a generic database manipulator with the `curd` capabilities.
func Database[M types.Model](ctx ...context.Context) types.Database[M] {
	if DB == nil || DB == new(gorm.DB) {
		panic("database is not initialized")
	}
	c := context.Background()
	if len(ctx) > 0 {
		if ctx[0] != nil {
			c = ctx[0]
		}
	}
	if strings.ToLower(config.App.LoggerConfig.Level) == "debug" {
		return &database[M]{db: DB.WithContext(c).Debug().Limit(defaultLimit), ctx: c}
	}
	return &database[M]{db: DB.WithContext(c).Limit(defaultLimit), ctx: c}
}

// trace returns timing function for database operations.
// It logs operation name, table name and elapsed time when done function is called.
//
// NOTE: trace function must be called after `defer db.reset()`.
func (db *database[M]) trace(op string) func(error) {
	begin := time.Now()
	return func(err error) {
		if err != nil {
			logger.Database.Errorw(op,
				"table", reflect.TypeOf(*new(M)).Elem().Name(),
				"cost", time.Since(begin).String(),
				"cache_enabled", db.enableCache,
				"try_run", db.tryRun,
				"error", err,
			)
		} else {
			logger.Database.Infow(op,
				"table", reflect.TypeOf(*new(M)).Elem().Name(),
				"cost", time.Since(begin).String(),
				"cache_enabled", db.enableCache,
				"try_run", db.tryRun,
			)
		}
	}
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

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
