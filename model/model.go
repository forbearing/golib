package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/forbearing/golib/types"
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

// Register table with records and it will be created in database before any curd..
func Register[M types.Model](records ...M) {
	mu.Lock()
	defer mu.Unlock()
	// table := *new(M)
	table := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	Tables = append(Tables, table)
	if len(records) != 0 {
		Records = append(Records, &Record{Table: table, Rows: records, Expands: table.Expands()})
	}
}

// RegisterTo same as Register but with custom database instances.
// dbname should always lowercase.
func RegisterTo[M types.Model](dbname string, records ...M) {
	mu.Lock()
	defer mu.Unlock()
	// table := *new(M)
	table := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	TablesWithDB = append(TablesWithDB, struct {
		Table  types.Model
		DBName string
	}{table, dbname})
	if len(records) != 0 {
		Records = append(Records, &Record{Table: table, Rows: records, Expands: table.Expands(), DBName: dbname})
	}
}

// Verb is router Verb
type Verb string

const (
	VerbCreate        Verb = "create"
	VerbDelete        Verb = "delete"
	VerbUpdate        Verb = "update"
	VerbUpdatePartial Verb = "update_partial"
	VerbList          Verb = "list"
	VerbGet           Verb = "get"
	VerbExport        Verb = "export"
	VerbImport        Verb = "import"

	// VerbMost includes verbs: `create`, `delete`, `update`, `update_partial`, `list`, `get`
	VerbMost Verb = "most"
	// VerbAll includes VerbGeneral, VerbImport, VerbExport
	VerbAll Verb = "all"
)

type route struct {
	Model types.Model
	Path  string
	Verbs []Verb
}

// RegisterRoutes register one route path with multiple api verbs.
// call this function multiple to register multiple route path.
// If route path is same, using latest register route path.
func RegisterRoutes[M types.Model](path string, verbs ...Verb) {
	mu.Lock()
	defer mu.Unlock()
	if len(path) != 0 && len(verbs) != 0 {
		Routes = append(Routes, route{Path: path, Verbs: verbs, Model: reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(types.Model)})
	}
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
	ID string `json:"id" gorm:"primaryKey" schema:"id"`

	CreatedBy      string    `json:"created_by,omitempty" schema:"created_by" gorm:"index"`
	UpdatedBy      string    `json:"updated_by,omitempty" schema:"updated_by" gorm:"index"`
	CreatedAt      time.Time `json:"created_at,omitempty" schema:"-" gorm:"index"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" schema:"-" gorm:"index"`
	Remark         *string   `json:"remark,omitempty" gorm:"size:10240" schema:"-"` // 如果需要支持 PATCH 更新,则必须是指针类型
	Order          *uint     `json:"order,omitempty" schema:"-"`
	Error          string    `json:"error,omitempty" schema:"-"`
	InternalRemark string    `json:"internal_remark,omitempty" schema:"-"` // 内部系统的备注

	// Query parameter
	Page       uint    `json:"-" gorm:"-" schema:"page"`         // Query parameter, eg: "page=2"
	Size       uint    `json:"-" gorm:"-" schema:"size"`         // Query parameter, eg: "size=10"
	Expand     *string `json:"-" gorm:"-" schema:"_expand"`      // Query parameter, eg: "_expand=children,parent".
	Depth      *uint   `json:"-" gorm:"-" schema:"_depth"`       // Query parameter, eg: "_depth=3".
	Fuzzy      *bool   `json:"-" gorm:"-" schema:"_fuzzy"`       // Query parameter, eg: "_fuzzy=true"
	SortBy     string  `json:"-" gorm:"-" schema:"_sortby"`      // Query parameter, eg: "_sortby=name"
	NoCache    bool    `json:"-" gorm:"-" schema:"_nocache"`     // Query parameter: eg: "_nocache=false"
	ColumnName string  `json:"-" gorm:"-" schema:"_column_name"` // Query parameter: eg: "_column_name=created_at"
	StartTime  string  `json:"-" gorm:"-" schema:"_start_time"`  // Query parameter: eg: "_start_time=2024-04-29+23:59:59"
	EndTime    string  `json:"-" gorm:"-" schema:"_end_time"`    // Query parameter: eg: "_end_time=2024-04-29+23:59:59"
	Or         *bool   `json:"-" gorm:"-" schema:"_or"`          // query parameter: eg: "_or=true"

	gorm.Model `json:"-"`
}

func (b *Base) GetTableName() string       { return "" }
func (b *Base) GetCreatedBy() string       { return b.CreatedBy }
func (b *Base) GetUpdatedBy() string       { return b.UpdatedBy }
func (b *Base) GetCreatedAt() time.Time    { return b.CreatedAt }
func (b *Base) GetUpdatedAt() time.Time    { return b.UpdatedAt }
func (b *Base) SetCreatedBy(s string)      { b.CreatedBy = s }
func (b *Base) SetUpdatedBy(s string)      { b.UpdatedBy = s }
func (b *Base) SetCreatedAt(t time.Time)   { b.CreatedAt = t }
func (b *Base) SetUpdatedAt(t time.Time)   { b.UpdatedAt = t }
func (b *Base) GetID() string              { return b.ID }
func (b *Base) SetID(id ...string)         { SetID(b, id...) }
func (b *Base) Expands() []string          { return nil }
func (b *Base) Excludes() map[string][]any { return nil }
func (b *Base) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// if b == nil {
	// 	return nil
	// }
	enc.AddString("id", b.ID)
	enc.AddString("created_by", b.CreatedBy)
	enc.AddString("updated_by", b.UpdatedBy)
	enc.AddUint("page", b.Page)
	enc.AddUint("size", b.Size)
	// enc.AddString("remark", util.Depointer(b.Remark))
	// enc.AddString("create_by", b.CreatedBy)
	// enc.AddString("updated_by", b.UpdatedBy)
	// enc.AddTime("created_at", b.CreatedAt)
	// enc.AddTime("updated_at", b.UpdatedAt)
	// enc.AddString("error", b.Error)
	return nil
}

// These methods implement types.Hooker interface.
// model should create custom hook to overwrite default hooks.
func (*Base) CreateBefore() error        { return nil }
func (*Base) CreateAfter() error         { return nil }
func (*Base) DeleteBefore() error        { return nil }
func (*Base) DeleteAfter() error         { return nil }
func (*Base) UpdateBefore() error        { return nil }
func (*Base) UpdateAfter() error         { return nil }
func (*Base) UpdatePartialBefore() error { return nil }
func (*Base) UpdatePartialAfter() error  { return nil }
func (*Base) ListBefore() error          { return nil }
func (*Base) ListAfter() error           { return nil }
func (*Base) GetBefore() error           { return nil }
func (*Base) GetAfter() error            { return nil }

func SetID(m types.Model, id ...string) {
	val := reflect.ValueOf(m).Elem()
	idField := val.FieldByName("ID")
	if len(idField.String()) != 0 {
		// zap.S().Debug("id already exits, skip")
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
	localTime, err := time.Parse(types.DATE_TIME_LAYOUT, string(value.([]byte)))
	if err != nil {
		return err
	}
	*t = GormTime(localTime)
	return nil
}

func (t GormTime) Value() (driver.Value, error) {
	return time.Time(t).Format(types.DATE_TIME_LAYOUT), nil
}

func (t *GormTime) UnmarshalJSON(b []byte) error {
	// Trim quotes from the stringified JSON value
	s := strings.Trim(string(b), "\"")
	// Parse the time using the custom format
	parsedTime, err := time.Parse(types.DATE_TIME_LAYOUT, s)
	if err != nil {
		return err
	}

	*t = GormTime(parsedTime)
	return nil
}

func (ct GormTime) MarshalJSON() ([]byte, error) {
	// Convert the time to the custom format and stringify it
	return []byte("\"" + time.Time(ct).Format(types.DATE_TIME_LAYOUT) + "\""), nil
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
		return errors.New("unsupported data type for GormStrings")
	}
	return nil
}

func (gs GormStrings) Value() (driver.Value, error) {
	// It will return "", if gs is nil or empty string.
	return strings.Trim(strings.Join(gs, ","), ","), nil
}

// GormScannerWrapper converts object to GormScanner that can be used in GORM.
// WARN: you must pass pointer to object.
func GormScannerWrapper(object any) *GormScanner {
	return &GormScanner{Object: object}
}

type GormScanner struct {
	Object any
}

func (g *GormScanner) Scan(value any) (err error) {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		err = json.Unmarshal(util.StringToBytes(v), g.Object)
	case []byte:
		err = json.Unmarshal(v, g.Object)
	default:
		err = errors.New("unsupported type")
	}
	return err
}

func (g *GormScanner) Value() (driver.Value, error) {
	data, err := json.Marshal(g.Object)
	if err != nil {
		return nil, err
	}
	return util.BytesToString(data), nil
}
