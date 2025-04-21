package model

// CasbinRule
// RBAC 包中会通过 gormadapter.NewAdapterByDBWithCustomTable(database.DB, &model.CasbinRule{})
// 或 gormadapter.NewAdapterByDB(database.DB) 来创建;
// NOTE: ID 类型必须是整型
type CasbinRule struct {
	// ID uint64 `json:"id" gorm:"primaryKey"`
	ID    uint64 `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Ptype string `json:"ptype" gorm:"size:100" schema:"ptype"`
	V0    string `json:"v0,omitempty" gorm:"size:100" schema:"v0"`
	V1    string `json:"v1,omitempty" gorm:"size:100" schema:"v1"`
	V2    string `json:"v2,omitempty" gorm:"size:100" schema:"v2"`
	V3    string `json:"v3,omitempty" gorm:"size:100" schema:"v3"`
	V4    string `json:"v4,omitempty" gorm:"size:100" schema:"v4"`
	V5    string `json:"v5,omitempty" gorm:"size:100" schema:"v5"`

	Base
}

// SetID 为一个空的函数,不允许自动设置ID, 因为 gormadapter.NewAdapterByDBWithCustomTable 创建的表的ID总是为 autoIncrement.
// 如果设置了自定义ID则会报错.
func (cr *CasbinRule) SetID(id ...string) {}

// GetTableName 用来指定 CasbinRule 在数据库中的表名为 casbin_rule.
// gormadapter.NewAdapterByDBWithCustomTable 创建的表名总是 casbin_rule, 但是 gorm 默认创建的表名为 casbin_rules,
// 为了统一,直接就用 casbin_rule 了.
func (cr CasbinRule) GetTableName() string { return "casbin_rule" }
