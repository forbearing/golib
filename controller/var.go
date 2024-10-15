package controller

import "github.com/forbearing/golib/types"

const (
	MAX_AVATAR_SIZE = 1024 * 1024 * 2   // 2M
	MAX_IMPORT_SIZE = 5 * 1024 * 1024   // 5M
	MAX_UPLOAD_SIZE = 1024 * 1024 * 100 // 100M
)

const (
	QUERY_ID          = types.QUERY_ID
	QUERY_PAGE        = types.QUERY_PAGE
	QUERY_SIZE        = types.QUERY_SIZE
	QUERY_LIMIT       = types.QUERY_LIMIT
	QUERY_EXPAND      = types.QUERY_EXPAND
	QUERY_DEPTH       = types.QUERY_DEPTH
	QUERY_OR          = types.QUERY_OR
	QUERY_FUZZY       = types.QUERY_FUZZY
	QUERY_SORTBY      = types.QUERY_SORTBY
	QUERY_COLUMN_NAME = types.QUERY_COLUMN_NAME
	QUERY_START_TIME  = types.QUERY_START_TIME
	QUERY_END_TIME    = types.QUERY_END_TIME
	QUERY_NOCACHE     = types.QUERY_NOCACHE
	QUERY_NOTOTAL     = types.QUERY_NOTOTAL
	QUERY_TYPE        = types.QUERY_TYPE
	QUERY_FILENAME    = types.QUERY_FILENAME
	QUERY_INDEX       = types.QUERY_INDEX
	QUERY_SELECT      = types.QUERY_SELECT

	PARAM_ID   = types.PARAM_ID
	PARAM_FILE = types.PARAM_FILE

	CTX_USERNAME = types.CTX_USERNAME

	VALUE_ALL = types.VALUE_ALL
)
