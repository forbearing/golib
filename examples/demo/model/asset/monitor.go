package model

import (
	"fmt"

	pkgmodel "github.com/forbearing/golib/model"
)

func init() {
	pkgmodel.Register[*Monitor]()
}

type Monitor struct {
	pkgmodel.Base

	Name string
}

// GetTableName returns a custom table name for the current model.
// This method is optional and can be omitted if using default table naming.
func (*Monitor) GetTableName() string {
	return fmt.Sprintf("%s_%s", PREFIX, pkgmodel.GetTableName[*Monitor]())
}
