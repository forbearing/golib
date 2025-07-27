package codegen

import (
	"github.com/forbearing/golib/types/consts"
	"github.com/stoewer/go-strcase"
)

var methods = []string{
	strcase.UpperCamelCase(string(consts.PHASE_CREATE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_CREATE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_DELETE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_DELETE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_UPDATE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_UPDATE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_UPDATE_PARTIAL_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_UPDATE_PARTIAL_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_LIST_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_LIST_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_GET_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_GET_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_CREATE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_CREATE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_DELETE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_DELETE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_UPDATE_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_UPDATE_AFTER)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_UPDATE_PARTIAL_BEFORE)),
	strcase.UpperCamelCase(string(consts.PHASE_BATCH_UPDATE_PARTIAL_AFTER)),
}
