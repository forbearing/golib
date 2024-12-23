package model

import (
	"fmt"

	pkgmodel "github.com/forbearing/golib/model"
)

func init() {
	pkgmodel.Register[*Computer]()
}

type Computer struct {
	pkgmodel.Base

	Name string `json:"name"`
}

// GetTableName returns a custom table name for the current model.
// This method is optional and can be omitted if using default table naming.
func (*Computer) GetTableName() string {
	return fmt.Sprintf("%s_%s", PREFIX, pkgmodel.GetTableName[*Computer]())
}