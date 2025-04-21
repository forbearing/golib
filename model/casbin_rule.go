package model

import (
	"hash/fnv"
	"strings"
)

// CasbinRule
// RBAC 包中会通过 gormadapter.NewAdapterByDBWithCustomTable(database.DB, &model.CasbinRule{})
// 或 gormadapter.NewAdapterByDB(database.DB) 来创建;
// NOTE: ID 类型必须是整型
type CasbinRule struct {
	// ID    uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ID    uint64 `json:"id" gorm:"primaryKey"`
	Ptype string `json:"ptype" gorm:"size:100" schema:"ptype"`
	V0    string `json:"v0,omitempty" gorm:"size:100" schema:"v0"`
	V1    string `json:"v1,omitempty" gorm:"size:100" schema:"v1"`
	V2    string `json:"v2,omitempty" gorm:"size:100" schema:"v2"`
	V3    string `json:"v3,omitempty" gorm:"size:100" schema:"v3"`
	V4    string `json:"v4,omitempty" gorm:"size:100" schema:"v4"`
	V5    string `json:"v5,omitempty" gorm:"size:100" schema:"v5"`

	Base
}

func (cr *CasbinRule) initDefaultValue() error {
	tokens := []string{cr.Ptype, cr.V0, cr.V1, cr.V2, cr.V3, cr.V4, cr.V5}
	// hash := sha256.Sum256([]byte(strings.Join(tokens, ":")))
	// id := uint64(binary.BigEndian.Uint64(hash[:8]))
	// return id

	h := fnv.New64a()
	h.Write([]byte(strings.Join(tokens, ":")))
	cr.ID = h.Sum64()
	return nil
}

func (cr *CasbinRule) UpdateBefore() error { return cr.initDefaultValue() }
func (cr *CasbinRule) CreateBefore() error { return cr.initDefaultValue() }
func (cr CasbinRule) GetTableName() string { return "casbin_rule" }

// type CasbinRule struct {
// 	ID    uint64 `gorm:"primaryKey;autoIncrement"`
// 	PType string `gorm:"size:100"`
// 	V0    string `gorm:"size:100"`
// 	V1    string `gorm:"size:100"`
// 	V2    string `gorm:"size:100"`
// 	V3    string `gorm:"size:100"`
// 	V4    string `gorm:"size:100"`
// 	V5    string `gorm:"size:100"`
// }
