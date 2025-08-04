package consts

import (
	"strings"

	"github.com/stoewer/go-strcase"
)

type AppSide string

const (
	Server AppSide = "server"
	Client AppSide = "client"
)

const (
	CTX_ROUTE      = "route"
	CTX_USERNAME   = "username"
	CTX_USER_ID    = "user_id"
	CTX_SESSION_ID = "session_id"

	DATE_TIME_LAYOUT = "2006-01-02 15:04:05"
	DATE_ID_LAYOUT   = "20060102"

	JOB_SCHEDULE_TIME_LAYOUT = "2006-01-02 15:04:05"
)

const (
	QUERY_ID            = "id"
	QUERY_PAGE          = "page"
	QUERY_SIZE          = "size"
	QUERY_LIMIT         = "limit"
	QUERY_EXPAND        = "_expand"
	QUERY_DEPTH         = "_depth"
	QUERY_OR            = "_or"
	QUERY_FUZZY         = "_fuzzy"
	QUERY_SORTBY        = "_sortby"
	QUERY_COLUMN_NAME   = "_column_name"
	QUERY_START_TIME    = "_start_time"
	QUERY_END_TIME      = "_end_time"
	QUERY_NOCACHE       = "_nocache"
	QUERY_TYPE          = "type"
	QUERY_FILENAME      = "filename"
	QUERY_NOTOTAL       = "_nototal"
	QUERY_INDEX         = "_index"
	QUERY_SELECT        = "_select"
	QUERY_CURSOR_VALUE  = "_cursor_value"
	QUERY_CURSOR_FIELDS = "_cursor_fields"
	QUERY_CURSOR_NEXT   = "_cursor_next"

	PARAM_ID   = "id"
	PARAM_FILE = "file"

	VALUE_ALL = "all"

	REQUEST_ID = "request_id"
	TRACE_ID   = "trace_id"
	SPAN_ID    = "span_id"
	PSPAN_ID   = "pspan_id"
	SEQ        = "seq"

	PREFIX_REQUEST_ID = "rq_"
	PREFIX_TRACE_ID   = "tr_"
	PREFIX_SPAN_ID    = "sp"
	PREFIX_PSPAN_ID   = "psp_"

	HEADER_REQUEST_ID = "X-Request-ID"
	HEADER_TRACE_ID   = "X-Trace-ID"
	HEADER_SPAN_ID    = "X-Span-ID"
	HEADER_PSPAN_ID   = "X-Pspan-ID"

	FIELD_ID = "ID"

	PHASE  = "phase"
	PARAMS = "params"
	QUERY  = "query"

	USER_SYSTEM = "system"
	USER_ROOT   = "root"

	TAG_JSON   = "json"
	TAG_SCHEMA = "schema"
	TAG_QUERY  = "query"
)

type Phase string

func (p Phase) MethodName() string {
	return strcase.UpperCamelCase(string(p))
}

const (
	PHASE_CREATE Phase = "create"
	PHASE_DELETE Phase = "delete"
	PHASE_UPDATE Phase = "update"
	PHASE_PATCH  Phase = "patch"
	PHASE_LIST   Phase = "list"
	PHASE_GET    Phase = "get"

	PHASE_CREATE_MANY Phase = "create_many"
	PHASE_DELETE_MANY Phase = "delete_many"
	PHASE_UPDATE_MANY Phase = "update_many"
	PHASE_PATCH_MANY  Phase = "patch_many"

	PHASE_CREATE_BEFORE Phase = "create_before"
	PHASE_CREATE_AFTER  Phase = "create_after"
	PHASE_DELETE_BEFORE Phase = "delete_before"
	PHASE_DELETE_AFTER  Phase = "delete_after"
	PHASE_UPDATE_BEFORE Phase = "update_before"
	PHASE_UPDATE_AFTER  Phase = "update_after"
	PHASE_PATCH_BEFORE  Phase = "patch_before"
	PHASE_PATCH_AFTER   Phase = "patch_after"
	PHASE_LIST_BEFORE   Phase = "list_before"
	PHASE_LIST_AFTER    Phase = "list_after"
	PHASE_GET_BEFORE    Phase = "get_before"
	PHASE_GET_AFTER     Phase = "get_after"

	PHASE_CREATE_MANY_BEFORE Phase = "create_many_before"
	PHASE_CREATE_MANY_AFTER  Phase = "create_many_after"
	PHASE_DELETE_MANY_BEFORE Phase = "delete_many_before"
	PHASE_DELETE_MANY_AFTER  Phase = "delete_many_after"
	PHASE_UPDATE_MANY_BEFORE Phase = "update_many_before"
	PHASE_UPDATE_MANY_AFTER  Phase = "update_many_after"
	PHASE_PATCH_MANY_BEFORE  Phase = "patch_many_before"
	PHASE_PATCH_MANY_AFTER   Phase = "patch_many_after"

	PHASE_FILTER Phase = "filter"
	PHASE_IMPORT Phase = "import"
	PHASE_EXPORT Phase = "export"
)

// HTTPVerb represents the supported HTTP operations for a resource
type HTTPVerb string

const (
	// Basic operations
	Create HTTPVerb = "create" // POST /resource
	Delete HTTPVerb = "delete" // DELETE /resource, DELETE /resource/:id
	Update HTTPVerb = "update" // PUT /resource, PUT /resource/:id
	Patch  HTTPVerb = "patch"  // PATCH /resource, PATCH /resource/:id
	List   HTTPVerb = "list"   // GET /resource
	Get    HTTPVerb = "get"    // GET /resource/:id

	CreateMany HTTPVerb = "create_many" // POST /resource/batch
	DeleteMany HTTPVerb = "delete_many" // DELETE /resource/batch
	UpdateMany HTTPVerb = "update_many" // PUT /resource/batch
	PatchMany  HTTPVerb = "patch_many"  // PATCH /resource/batch

	Export HTTPVerb = "export" // GET /resource/export
	Import HTTPVerb = "import" // POST /resource/import

	// Verb groups
	Most      HTTPVerb = "most"       // Basic CRUD operations (Create, Delete, Update, UpdatePartial, List, Get)
	MostBatch HTTPVerb = "most_batch" // Basic batch operations (CreateMany, DeleteMany, UpdateMany, PatchMany)
	All       HTTPVerb = "all"        // All operations including Most, MostBatch, Import and Export
)

func (v HTTPVerb) String() string {
	return strings.ReplaceAll(string(v), "_", " ")
}

const (
	ROOT  = "root"
	ADMIN = "admin"
)

const LayoutTimeEncoder = "2006-01-02|15:04:05"

const IMPORT_PATH_MODEL = `"github.com/forbearing/golib/model"`
