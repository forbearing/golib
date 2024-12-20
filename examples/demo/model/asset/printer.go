package model

import (
	"fmt"

	pkgmodel "github.com/forbearing/golib/model"
)

type Printer struct {
	pkgmodel.Base

	Name string
}

// GetTableName returns a custom table name for the current model.
// This method is optional and can be omitted if using default table naming.
func (*Printer) GetTableName() string {
	return fmt.Sprintf("%s_%s", PREFIX, pkgmodel.GetTableName[*Printer]())
}
