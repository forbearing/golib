package model

import (
	"fmt"

	pkgmodel "github.com/forbearing/golib/model"
)

func init() {
	pkgmodel.Register[*Software]()
}

type Software struct {
	pkgmodel.Base

	Name string `json:"name"`
}

// GetTableName returns a custom table name for the current model.
// This method is optional and can be omitted if using default table naming.
func (*Software) GetTableName() string {
	return fmt.Sprintf("%s_%s", PREFIX, pkgmodel.GetTableName[*Software]())
}
