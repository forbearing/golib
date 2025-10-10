package model

import (
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/util"
	"go.uber.org/zap/zapcore"
)

type Fixed string

const (
	FIXED_RIGHT Fixed = "right" //nolint:staticcheck
	FIXED_LEFT  Fixed = "left"  //nolint:staticcheck
)

func init() {
	Register[*TableColumn]()
}

// TableColumn 表格的列
type TableColumn struct {
	UserID    string `json:"user_id,omitempty" schema:"user_id"`       // 属于哪一个用户的
	TableName string `json:"table_name,omitempty" schema:"table_name"` // 属于哪一张表的
	Name      string `json:"name,omitempty" schema:"name"`             // 列名
	Key       string `json:"key,omitempty" schema:"key"`               // 列名对应的id

	Width    *uint   `json:"width,omitempty"`    // 列宽度
	Sequence *uint   `json:"sequence,omitempty"` // 列顺序
	Visiable *bool   `json:"visiable,omitempty"` // 是否显示
	Fixed    *string `json:"fixed,omitempty"`    // 固定在哪里 left,right, 必须加上 omitempty

	Base
}

func (t *TableColumn) CreateBefore(*types.ModelContext) error {
	if (t.Visiable) == nil {
		t.Visiable = util.ValueOf(true)
	}
	// id cannot be hidden
	if t.Key == "id" {
		t.Visiable = util.ValueOf(true)
	}
	return nil
}

func (t *TableColumn) UpdateBefore(*types.ModelContext) error {
	// id cannot be hidden
	if t.Key == "id" {
		t.Visiable = util.Pointer(true)
	}
	return nil
}

func (t *TableColumn) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("user_id", t.UserID)
	enc.AddString("table_name", t.TableName)
	enc.AddString("name", t.Name)
	enc.AddString("key", t.Key)
	if t.Width != nil {
		enc.AddUint("width", *t.Width)
	}
	if t.Sequence != nil {
		enc.AddUint("sequence", *t.Sequence)
	}
	if t.Visiable != nil {
		enc.AddBool("visiable", *t.Visiable)
	}
	if t.Fixed != nil {
		enc.AddString("fixed", *t.Fixed)
	}
	return nil
}
