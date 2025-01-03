package controller

import "github.com/forbearing/golib/types/consts"

const (
	MAX_AVATAR_SIZE = 1024 * 1024 * 2   // 2M
	MAX_IMPORT_SIZE = 5 * 1024 * 1024   // 5M
	MAX_UPLOAD_SIZE = 1024 * 1024 * 100 // 100M
)

const (
	QUERY_ID          = consts.QUERY_ID
	QUERY_PAGE        = consts.QUERY_PAGE
	QUERY_SIZE        = consts.QUERY_SIZE
	QUERY_LIMIT       = consts.QUERY_LIMIT
	QUERY_EXPAND      = consts.QUERY_EXPAND
	QUERY_DEPTH       = consts.QUERY_DEPTH
	QUERY_OR          = consts.QUERY_OR
	QUERY_FUZZY       = consts.QUERY_FUZZY
	QUERY_SORTBY      = consts.QUERY_SORTBY
	QUERY_COLUMN_NAME = consts.QUERY_COLUMN_NAME
	QUERY_START_TIME  = consts.QUERY_START_TIME
	QUERY_END_TIME    = consts.QUERY_END_TIME
	QUERY_NOCACHE     = consts.QUERY_NOCACHE
	QUERY_NOTOTAL     = consts.QUERY_NOTOTAL
	QUERY_TYPE        = consts.QUERY_TYPE
	QUERY_FILENAME    = consts.QUERY_FILENAME
	QUERY_INDEX       = consts.QUERY_INDEX
	QUERY_SELECT      = consts.QUERY_SELECT

	PARAM_ID   = consts.PARAM_ID
	PARAM_FILE = consts.PARAM_FILE

	CTX_USERNAME = consts.CTX_USERNAME

	VALUE_ALL = consts.VALUE_ALL
)
