package model

import (
	"bytes"
	"database/sql/driver"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

var ErrMobileLength = errors.New("mobile number length must be 11")

var (
	// Records is the table records must be pr-eexists before any database curd,
	// its register by register function.
	// The underlying type of map value must be model pointer to structure, eg: *model.User
	//
	// Records is the table records that should created automatically when app bootstraping.
	Records []*Record = make([]*Record, 0)

	// Tables is the database table that should created automatically when app bootstraping.
	Tables []types.Model

	TablesWithDB []struct {
		Table  types.Model
		DBName string
	}

	// Routes is a map slice that element is map[string][]Verb
	// The map key is the route path, eg: '/asset' or 'asset'
	//   or '/api/asset' (the prefiex /api will be remove automatically).
	// The map value is the Verb slice.
	// - VerbCreate is equivalent to http method "POST".
	// - VerbDelete is equivalent to http method "DELETE".
	// - VerbUpdate is equivalent to http method "PUT".
	// - VerbUpdatePartial is equivalent to http method "PATCH".
	// - VerbList is equivalent to http method "GET /xxx".
	// - VerbGET is equivalent to http method "GET /xxx/:id".
	// - VerbExport is equivalent to http method "GET" but specifially used to export resources.
	// - VerbImport is equivalent to http method "POST" but specifically used to import resources.
	Routes []route = make([]route, 0)

	mu sync.Mutex
)

// Record is table record
type Record struct {
	Table   types.Model
	Rows    any
	Expands []string
	DBName  string
}
type route struct {
	Model types.Model
	Path  string
	Verbs []consts.HTTPVerb
}

// Register associates the model with database table and will created automatically.
// If records provided, they will be inserted when application bootstrapping.
//
// Parameters:
//   - records: Optional initial records to be seeded into the table. Can be single or multiple records.
//
// Examples:
//
//	// Create table 'users' only
//	Register[*model.User]()
//
//	// Create table 'users' and insert one record
//	Register[*model.User](&model.User{ID: 1, Name: "admin"})
//
//	// Create table 'users' and insert a single user record
//	Register[*model.User](user)
//
//	// Create table 'users' and insert multiple records
//	Register[*model.User](users...)  // where users is []*model.User
//
// NOTE:
//  1. Always call this function in init().
//  2. Ensure the model pacakge is imported in main.go.
//     The init() function will only executed if the file is imported directly or indirectly by main.go.
func Register[M types.Model](records ...M) {
	mu.Lock()
	defer mu.Unlock()
	// table := *new(M)
	table := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	Tables = append(Tables, table)
	// NOTE: it's necessary to set id before insert.
	for i := range records {
		if len(records[i].GetID()) == 0 {
			records[i].SetID()
		}
	}
	if len(records) != 0 {
		Records = append(Records, &Record{Table: table, Rows: records, Expands: table.Expands()})
	}
}

// RegisterTo works identically to Register(), but registers the model on the specified database instance.
// more details see: Register().
func RegisterTo[M types.Model](dbname string, records ...M) {
	mu.Lock()
	defer mu.Unlock()
	dbname = strings.ToLower(dbname)
	table := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	TablesWithDB = append(TablesWithDB, struct {
		Table  types.Model
		DBName string
	}{table, dbname})
	if len(records) != 0 {
		Records = append(Records, &Record{Table: table, Rows: records, Expands: table.Expands(), DBName: dbname})
	}
}

// RegisterRoutes register one route path with multiple api verbs.
// call this function multiple to register multiple route path.
// If route path is same, using latest register route path.
//
// Deprecated: use router.Register() instead. This function is a no-op.
func RegisterRoutes[M types.Model](path string, verbs ...consts.HTTPVerb) {
	// mu.Lock()
	// defer mu.Unlock()
	// if len(path) != 0 && len(verbs) != 0 {
	// 	Routes = append(Routes, route{Path: path, Verbs: verbs, Model: reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(types.Model)})
	// }
}

var _ types.Model = (*Base)(nil)

// Base implement types.Model interface.
// Each model must be expands the Base structure.
// You can implements your custom method to overwrite the defaults methods.
//
// Usually, there are some gorm tags that may be of interest to you.
// gorm:"unique"
// gorm:"foreignKey:ParentID"
// gorm:"foreignKey:ParentID,references:ID"
type Base struct {
	ID string `json:"id" gorm:"primaryKey" schema:"id" url:"-"`

	CreatedBy string     `json:"created_by,omitempty" gorm:"index" schema:"created_by" url:"-"`
	UpdatedBy string     `json:"updated_by,omitempty" gorm:"index" schema:"updated_by" url:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"index" schema:"-" url:"-"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"index" schema:"-" url:"-"`
	Remark    *string    `json:"remark,omitempty" gorm:"size:10240" schema:"-" url:"-"` // 如果需要支持 PATCH 更新,则必须是指针类型
	Order     *uint      `json:"order,omitempty" schema:"-" url:"-"`

	// Query parameter
	Page       uint    `json:"-" gorm:"-" schema:"page" url:"page,omitempty"`                 // Query parameter, eg: "page=2"
	Size       uint    `json:"-" gorm:"-" schema:"size" url:"size,omitempty"`                 // Query parameter, eg: "size=10"
	Expand     *string `json:"-" gorm:"-" schema:"_expand" url:"_expand,omitempty"`           // Query parameter, eg: "_expand=children,parent".
	Depth      *uint   `json:"-" gorm:"-" schema:"_depth" url:"_depth,omitempty"`             // Query parameter, eg: "_depth=3".
	Fuzzy      *bool   `json:"-" gorm:"-" schema:"_fuzzy" url:"_fuzzy,omitempty"`             // Query parameter, eg: "_fuzzy=true"
	SortBy     string  `json:"-" gorm:"-" schema:"_sortby" url:"_sortby,omitempty"`           // Query parameter, eg: "_sortby=name"
	NoCache    bool    `json:"-" gorm:"-" schema:"_nocache" url:"_nocache,omitempty"`         // Query parameter: eg: "_nocache=false"
	ColumnName string  `json:"-" gorm:"-" schema:"_column_name" url:"_column_name,omitempty"` // Query parameter: eg: "_column_name=created_at"
	StartTime  string  `json:"-" gorm:"-" schema:"_start_time" url:"_start_time,omitempty"`   // Query parameter: eg: "_start_time=2024-04-29+23:59:59"
	EndTime    string  `json:"-" gorm:"-" schema:"_end_time" url:"_end_time,omitempty"`       // Query parameter: eg: "_end_time=2024-04-29+23:59:59"
	Or         *bool   `json:"-" gorm:"-" schema:"_or" url:"_or,omitempty"`                   // query parameter: eg: "_or=true"
	Index      string  `json:"-" gorm:"-" schema:"_index" url:"_index,omitempty"`             // Query parameter: eg: "_index=name"
	Select     string  `json:"-" gorm:"-" schema:"_select" url:"_select,omitempty"`           // Query parameter: eg: "_select=field1,field2"

	gorm.Model `json:"-" schema:"-" url:"-"`
}

func (b *Base) GetTableName() string       { return "" }
func (b *Base) GetCreatedBy() string       { return b.CreatedBy }
func (b *Base) GetUpdatedBy() string       { return b.UpdatedBy }
func (b *Base) GetCreatedAt() time.Time    { return util.Deref(b.CreatedAt) }
func (b *Base) GetUpdatedAt() time.Time    { return util.Deref(b.UpdatedAt) }
func (b *Base) SetCreatedBy(s string)      { b.CreatedBy = s }
func (b *Base) SetUpdatedBy(s string)      { b.UpdatedBy = s }
func (b *Base) SetCreatedAt(t time.Time)   { b.CreatedAt = &t }
func (b *Base) SetUpdatedAt(t time.Time)   { b.UpdatedAt = &t }
func (b *Base) GetID() string              { return b.ID }
func (b *Base) SetID(id ...string)         { setID(b, id...) }
func (b *Base) Expands() []string          { return nil }
func (b *Base) Excludes() map[string][]any { return nil }
func (b *Base) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", b.ID)
	enc.AddString("created_by", b.CreatedBy)
	enc.AddString("updated_by", b.UpdatedBy)
	enc.AddUint("page", b.Page)
	enc.AddUint("size", b.Size)
	return nil
}

// These methods implement types.Hooker interface.
// model should create custom hook to overwrite default hooks.
func (*Base) CreateBefore() error { return nil }
func (*Base) CreateAfter() error  { return nil }
func (*Base) DeleteBefore() error { return nil }
func (*Base) DeleteAfter() error  { return nil }
func (*Base) UpdateBefore() error { return nil }
func (*Base) UpdateAfter() error  { return nil }
func (*Base) ListBefore() error   { return nil }
func (*Base) ListAfter() error    { return nil }
func (*Base) GetBefore() error    { return nil }
func (*Base) GetAfter() error     { return nil }

func setID(m types.Model, id ...string) {
	val := reflect.ValueOf(m).Elem()
	idField := val.FieldByName("ID")
	if len(idField.String()) != 0 {
		return
	}
	if len(id) == 0 {
		idField.SetString(util.UUID())
		return
	}

	zap.S().Debug("setting id: " + id[0])
	if len(id[0]) == 0 {
		idField.SetString(util.UUID())
	} else {
		idField.SetString(id[0])
	}
}

type GormTime time.Time

func (t *GormTime) Scan(value any) error {
	localTime, err := time.Parse(consts.DATE_TIME_LAYOUT, string(value.([]byte)))
	if err != nil {
		return err
	}
	*t = GormTime(localTime)
	return nil
}

func (t GormTime) Value() (driver.Value, error) {
	return time.Time(t).Format(consts.DATE_TIME_LAYOUT), nil
}

func (t *GormTime) UnmarshalJSON(b []byte) error {
	// Trim quotes from the stringified JSON value
	s := strings.Trim(string(b), "\"")
	// Parse the time using the custom format
	parsedTime, err := time.Parse(consts.DATE_TIME_LAYOUT, s)
	if err != nil {
		return err
	}

	*t = GormTime(parsedTime)
	return nil
}

func (ct GormTime) MarshalJSON() ([]byte, error) {
	// Convert the time to the custom format and stringify it
	return []byte("\"" + time.Time(ct).Format(consts.DATE_TIME_LAYOUT) + "\""), nil
}

type GormStrings []string

func (gs *GormStrings) Scan(value any) error {
	if value == nil {
		*gs = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		_v := bytes.TrimSpace(v)
		_v = bytes.Trim(_v, ",")
		*gs = strings.Split(string(_v), ",")
	case string:
		_v := strings.TrimSpace(v)
		_v = strings.Trim(_v, ",")
		*gs = strings.Split(_v, ",")
	default:
		return errors.New("unsupported type for GormStrings, expected []byte or string")
	}
	return nil
}
