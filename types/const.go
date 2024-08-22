package types

const (
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

	PARAM_ID   = "id"
	PARAM_FILE = "file"

	VALUE_ALL = "all"

	REQUEST_ID = "request_id"

	PHASE = "phase"

	USER_SYSTEM = "system"
	USER_ROOT   = "root"
)

type Phase string

const (
	PHASE_CREATE_BEFORE         Phase = "create_before"
	PHASE_CREATE_AFTER          Phase = "create_after"
	PHASE_UPDATE_BEFORE         Phase = "update_before"
	PHASE_UPDATE_AFTER          Phase = "update_after"
	PHASE_UPDATE_PARTIAL_BEFORE Phase = "update_partial_before"
	PHASE_UPDATE_PARTIAL_AFTER  Phase = "update_partial_after"
	PHASE_DELETE_BEFORE         Phase = "delete_before"
	PHASE_DELETE_AFTER          Phase = "delete_after"
	PHASE_LIST_BEFORE           Phase = "list_before"
	PHASE_LIST_AFTER            Phase = "list_after"
	PHASE_GET_BEFORE            Phase = "get_before"
	PHASE_GET_AFTER             Phase = "get_after"
	PHASE_FILTER                Phase = "filter"
	PHASE_IMPORT                Phase = "import"
	PHASE_EXPORT                Phase = "export"
)
