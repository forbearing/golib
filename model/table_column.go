package model

type Fixed string

const (
	FIXED_RIGHT Fixed = "right"
	FIXED_LEFT  Fixed = "left"
)

// TableColumn 表格的列
type TableColumn struct {
	UserId    string `json:"user_id,omitempty" schema:"user_id"`       // 属于哪一个用户的
	TableName string `json:"table_name,omitempty" schema:"table_name"` // 属于哪一张表的
	Name      string `json:"name,omitempty" schema:"name"`             // 列名
	Key       string `json:"key,omitempty" schema:"key"`               // 列名对应的id

	Width    *uint  `json:"width,omitempty"`    // 列宽度
	Sequence *uint  `json:"sequence,omitempty"` // 列顺序
	Visiable *bool  `json:"visiable,omitempty"` // 是否显示
	Fixed    *Fixed `json:"fixed,omitempty"`    // 固定在哪里 left,right, 必须加上 omitempty

	Base
}
