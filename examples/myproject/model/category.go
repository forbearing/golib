package model

import (
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

var (
	// 根分类, 不会展示出来 所有的一级分类的 parentId 都是这个根分类.
	// 数据库初始化时, 会自动创建这个记录
	categoryRoot = &Category{Name: RootName, Status: util.Pointer(uint(0)), ParentId: RootId, Base: model.Base{ID: RootId}}
	// 未知分类, 导入数据时,如果发现 parentId 为空,将将其 parentId 指向这个未知分类
	// 数据库初始化时, 会自动创建这个记录
	categoryUnknown = &Category{Name: UnknownName, Status: util.Pointer(uint(0)), ParentId: UnknownId, Base: model.Base{ID: UnknownId}}
	categoryNone    = &Category{Name: NoneName, Status: util.Pointer(uint(0)), ParentId: RootId, Base: model.Base{ID: NoneId}}
)

func init() {
	model.Register[*Category](categoryRoot, categoryUnknown, categoryNone)
	model.RegisterRoutes[*Category]("category")
}

type Category struct {
	Name     string     `json:"name,omitempty" gorm:"unique" schema:"name"`
	Status   *uint      `json:"status,omitempty" gorm:"type:tinyint(1);comment:status(0: disabled, 1: enable)" schema:"status"`
	ParentId string     `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentId"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"`

	model.Base
}

func (*Category) Expands() []string {
	return []string{"Children", "Parent"}
}
func (*Category) Excludes() map[string][]any {
	return map[string][]any{KeyName: {RootName, UnknownName, NoneName}}
}
func (c *Category) CreateBefore() error {
	if len(c.ParentId) == 0 {
		c.ParentId = RootId
	}
	return nil
}
func (c *Category) UpdateBefore() error {
	return c.CreateBefore()
}

func (c *Category) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c == nil {
		return nil
	}
	enc.AddString("id", c.ID)
	enc.AddString("name", c.Name)
	enc.AddUint("status", util.Depointer(c.Status))
	enc.AddString("parent_id", c.ParentId)
	enc.AddObject("base", &c.Base)
	return nil
}
