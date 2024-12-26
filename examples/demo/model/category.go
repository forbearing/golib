package model

import (
	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	rootId      = "root"
	rootName    = "root"
	unknownId   = "unknown"
	unknownName = "unknown"
	noneId      = "none"
	noneName    = "无"

	keyName = "name"
	keyId   = "id"
)

var (
	categoryRoot = &Category{Name: rootName, Status: util.Pointer(uint(0)), ParentId: rootId, Base: pkgmodel.Base{ID: rootId}}
	// 未知分类, 导入数据时,如果发现 parentId 为空,将将其 parentId 指向这个未知分类
	// 数据库初始化时, 会自动创建这个记录
	categoryUnknown = &Category{Name: unknownName, Status: util.Pointer(uint(0)), ParentId: unknownId, Base: pkgmodel.Base{ID: unknownId}}
	categoryNone    = &Category{Name: noneName, Status: util.Pointer(uint(0)), ParentId: rootId, Base: pkgmodel.Base{ID: noneId}}

	categoryFruit  = &Category{Name: "fruit", Status: util.Pointer(uint(1)), ParentId: rootId, Base: pkgmodel.Base{ID: "fruit"}}
	categoryApple  = &Category{Name: "apple", Status: util.Pointer(uint(1)), ParentId: categoryFruit.ID, Base: pkgmodel.Base{ID: "apple"}}
	categoryBanana = &Category{Name: "banana", Status: util.Pointer(uint(1)), ParentId: categoryFruit.ID, Base: pkgmodel.Base{ID: "banana"}}

	categoryApple1  = &Category{Name: "apple1", Status: util.Pointer(uint(1)), ParentId: categoryApple.ID, Base: pkgmodel.Base{ID: "apple1"}}
	categoryApple2  = &Category{Name: "apple2", Status: util.Pointer(uint(1)), ParentId: categoryApple.ID, Base: pkgmodel.Base{ID: "apple2"}}
	categoryBanana1 = &Category{Name: "banana1", Status: util.Pointer(uint(1)), ParentId: categoryBanana.ID, Base: pkgmodel.Base{ID: "banana1"}}
	categoryBanana2 = &Category{Name: "banana2", Status: util.Pointer(uint(1)), ParentId: categoryBanana.ID, Base: pkgmodel.Base{ID: "banana2"}}
)

func init() {
	pkgmodel.Register(
		categoryRoot, categoryNone, categoryUnknown,
		categoryFruit, categoryApple, categoryBanana,
		categoryApple1, categoryApple2,
		categoryBanana1, categoryBanana2,
	)
}

type Category struct {
	pkgmodel.Base

	Name     string     `json:"name,omitempty" gorm:"unique" schema:"name"`
	Status   *uint      `json:"status,omitempty" gorm:"type:smallint;comment:status(0: disabled, 1: enabled)" schema:"status"`
	ParentId string     `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentId"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"`
}

func (*Category) Expands() []string {
	return []string{"Children", "Parent"}
}

func (*Category) Excludes() map[string][]any {
	return map[string][]any{keyName: {rootName, unknownName, noneName}}
}

func (c *Category) CreateBefore() error {
	if len(c.ParentId) == 0 {
		c.ParentId = rootId
	}
	zap.L().Info("", zap.Object("category", c))
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
	enc.AddUint("status", util.Deref(c.Status))
	enc.AddString("parent_id", c.ParentId)
	enc.AddObject("base", &c.Base)
	return nil
}
