package model_log

import "github.com/forbearing/golib/model"

type OperationType string

const (
	OperationTypeCreate        OperationType = "create"
	OperationTypeDelete        OperationType = "delete"
	OperationTypeUpdate        OperationType = "update"
	OperationTypeUpdatePartial OperationType = "update_partial"
	OperationTypeList          OperationType = "list"
	OperationTypeGet           OperationType = "get"
	OperationTypeExport        OperationType = "export"
	OperationTypeImport        OperationType = "import"

	OperationTypeBatchCreate        OperationType = "batch_create"
	OperationTypeBatchDelete        OperationType = "batch_delete"
	OperationTypeBatchUpdate        OperationType = "batch_update"
	OperationTypeBatchUpdatePartial OperationType = "batch_update_partial"
)

type OperationLog struct {
	User       string        `json:"user,omitempty" schema:"user"`   // 操作者, 本地账号该字段为空,例如 root
	IP         string        `json:"ip,omitempty" schema:"ip"`       // 操作者的 ip
	Op         OperationType `json:"op,omitempty" schema:"op"`       // 动作: 增删改查
	Table      string        `json:"table,omitempty" schema:"table"` // 操作了哪张表
	Model      string        `json:"model,omitempty" schema:"model"`
	RecordId   string        `json:"record_id,omitempty" schema:"record_id"`     // 表记录的 id
	RecordName string        `json:"record_name,omitempty" schema:"record_name"` // 表记录的 name
	Record     string        `json:"record,omitempty" schema:"record"`           // 记录全部内容
	OldRecord  string        `json:"old_record,omitempty"`                       // 更新前的内容
	NewRecord  string        `json:"new_record,omitempty"`                       // 更新后的内容
	Method     string        `json:"method,omitempty" schema:"method"`
	URI        string        `json:"uri,omitempty" schema:"uri"` // request uri
	UserAgent  string        `json:"user_agent,omitempty" schema:"user_agent"`
	RequestId  string        `json:"request_id,omitempty" schema:"request_id"`

	model.Base
}
