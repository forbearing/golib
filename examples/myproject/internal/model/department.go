package model

import (
	pkgmodel "github.com/forbearing/golib/model"
)

func init() {
	pkgmodel.Register[*Department]()
}

type Department struct {
	pkgmodel.Base
}

func (*Department) Request(DepartmentReq, DepartmentResp) {}

type DepartmentReq struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type DepartmentResp struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
