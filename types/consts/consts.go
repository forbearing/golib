package consts

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
	FileRbacConf = "rbac.conf"
)

const (
	QUERY_ID          = "id"
	QUERY_PAGE        = "page"
	QUERY_SIZE        = "size"
	QUERY_LIMIT       = "limit"
	QUERY_EXPAND      = "_expand"
	QUERY_DEPTH       = "_depth"
	QUERY_OR          = "_or"
	QUERY_FUZZY       = "_fuzzy"
	QUERY_SORTBY      = "_sortby"
	QUERY_COLUMN_NAME = "_column_name"
	QUERY_START_TIME  = "_start_time"
	QUERY_END_TIME    = "_end_time"
	QUERY_NOCACHE     = "_nocache"
	QUERY_TYPE        = "type"
	QUERY_FILENAME    = "filename"
	QUERY_NOTOTAL     = "_nototal"
	QUERY_INDEX       = "_index"
	QUERY_SELECT      = "_select"

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

	PHASE  = "phase"
	PARAMS = "params"
	QUERY  = "query"

	USER_SYSTEM = "system"
	USER_ROOT   = "root"
)

type Phase string

const (
	PHASE_CREATE                Phase = "create"
	PHASE_CREATE_BEFORE         Phase = "create_before"
	PHASE_CREATE_AFTER          Phase = "create_after"
	PHASE_UPDATE                Phase = "update"
	PHASE_UPDATE_BEFORE         Phase = "update_before"
	PHASE_UPDATE_AFTER          Phase = "update_after"
	PHASE_UPDATE_PARTIAL        Phase = "update_partial"
	PHASE_UPDATE_PARTIAL_BEFORE Phase = "update_partial_before"
	PHASE_UPDATE_PARTIAL_AFTER  Phase = "update_partial_after"
	PHASE_DELETE                Phase = "delete"
	PHASE_DELETE_BEFORE         Phase = "delete_before"
	PHASE_DELETE_AFTER          Phase = "delete_after"
	PHASE_LIST                  Phase = "list"
	PHASE_LIST_BEFORE           Phase = "list_before"
	PHASE_LIST_AFTER            Phase = "list_after"
	PHASE_GET                   Phase = "get"
	PHASE_GET_BEFORE            Phase = "get_before"
	PHASE_GET_AFTER             Phase = "get_after"

	PHASE_BATCH_CREATE                Phase = "batch_create"
	PHASE_BATCH_CREATE_BEFORE         Phase = "batch_create_before"
	PHASE_BATCH_CREATE_AFTER          Phase = "batch_create_after"
	PHASE_BATCH_DELETE                Phase = "batch_delete"
	PHASE_BATCH_DELETE_BEFORE         Phase = "batch_delete_before"
	PHASE_BATCH_DELETE_AFTER          Phase = "batch_delete_after"
	PHASE_BATCH_UPDATE                Phase = "batch_update"
	PHASE_BATCH_UPDATE_BEFORE         Phase = "batch_update_before"
	PHASE_BATCH_UPDATE_AFTER          Phase = "batch_update_after"
	PHASE_BATCH_UPDATE_PARTIAL        Phase = "batch_update_partial"
	PHASE_BATCH_UPDATE_PARTIAL_BEFORE Phase = "batch_update_partial_before"
	PHASE_BATCH_UPDATE_PARTIAL_AFTER  Phase = "batch_update_partial_after"

	PHASE_FILTER Phase = "filter"
	PHASE_IMPORT Phase = "import"
	PHASE_EXPORT Phase = "export"
)

// HTTPVerb represents the supported HTTP operations for a resource
type HTTPVerb string

const (
	// Basic operations
	Create        HTTPVerb = "create"         // POST /resource
	Delete        HTTPVerb = "delete"         // DELETE /resource, DELETE /resource/:id
	Update        HTTPVerb = "update"         // PUT /resource, PUT /resource/:id
	UpdatePartial HTTPVerb = "update_partial" // PATCH /resource, PATCH /resource/:id
	List          HTTPVerb = "list"           // GET /resource
	Get           HTTPVerb = "get"            // GET /resource/:id

	BatchCreate        HTTPVerb = "batch_create"         // POST /resource/batch
	BatchDelete        HTTPVerb = "batch_delete"         // DELETE /resource/batch
	BatchUpdate        HTTPVerb = "batch_update"         // PUT /resource/batch
	BatchUpdatePartial HTTPVerb = "batch_update_partial" // PATCH /resource/batch

	Export HTTPVerb = "export" // GET /resource/export
	Import HTTPVerb = "import" // POST /resource/import

	// Verb groups
	Most HTTPVerb = "most" // Basic CRUD operations (Create, Delete, Update, UpdatePartial, List, Get)
	All  HTTPVerb = "all"  // All operations including Import and Export
)

const (
	ROOT  = "root"
	ADMIN = "admin"
)

const LayoutTimeEncoder = "2006-01-02|15:04:05"
