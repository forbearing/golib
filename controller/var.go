package controller

import "github.com/forbearing/golib/types"

const (
	MAX_AVATAR_SIZE = 1024 * 1024 * 2   // 2M
	MAX_IMPORT_SIZE = 5 * 1024 * 1024   // 5M
	MAX_UPLOAD_SIZE = 1024 * 1024 * 100 // 100M
)

const (
	QUERY_ID       = types.QUERY_ID
	QUERY_PAGE     = types.QUERY_PAGE
	QUERY_SIZE     = types.QUERY_SIZE
	QUERY_EXPAND   = types.QUERY_EXPAND
	QUERY_DEPTH    = types.QUERY_DEPTH
	QUERY_FUZZY    = types.QUERY_FUZZY
	QUERY_SORTBY   = types.QUERY_SORTBY
	QUERY_NOCACHE  = types.QUERY_NOCACHE
	QUERY_TYPE     = types.QUERY_TYPE
	QUERY_FILENAME = types.QUERY_FILENAME

	PARAM_ID   = types.PARAM_ID
	PARAM_FILE = types.PARAM_FILE

	CTX_USERNAME = types.CTX_USERNAME

	VALUE_ALL = types.VALUE_ALL
)
