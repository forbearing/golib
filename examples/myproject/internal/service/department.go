package service

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

func init() {
	service.Register[*department]()
}

type department struct {
	service.Base[*model.Department]
}

func (d *department) CreateBefore(ctx *types.ServiceContext, department *model.Department) error {
	req, ok := ctx.GetRequest().(model.DepartmentReq)
	if !ok {
		return nil
	}

	d.WithServiceContext(ctx, ctx.GetPhase()).Error("department create before")

	return database.Database[*model.Group]().Create(&model.Group{Name: req.Name, Desc: req.Desc})
}
