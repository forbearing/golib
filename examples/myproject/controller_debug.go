package main

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/examples/myproject/internal/model"
	. "github.com/forbearing/golib/response"
	"github.com/gin-gonic/gin"
)

type debug struct{}

var Debug = new(debug)

func (debug) Debug(c *gin.Context) {
	// storages1 := []model.StorageDevice{
	// 	{Name: "storage-device-1"},
	// 	{Name: "storage-device-2"},
	// }
	// storageDevices2 := []model.StorageDevice{
	// 	{Name: "storage-device-1"},
	// 	{Name: "storage-device-2"},
	// }
	// devices1 := []model.NetworkDevice{
	// 	{Name: "netework-device-1"},
	// 	{Name: "netework-device-2"},
	// }
	// networkDevices2 := []model.NetworkDevice{
	// 	{Name: "netework-device-3"},
	// 	{Name: "netework-device-4"},
	// }
	//
	// if err := database.Database[*model.SysInfo]().Update([]*model.SysInfo{
	// 	{Node: model.Node{MachineID: "machine-id-1"}, Storages: storages1, Networks: devices1},
	// 	{Node: model.Node{MachineID: "machine-id-2"}, Storages: storageDevices2, Networks: networkDevices2},
	// }...); err != nil {
	// 	zap.S().Error(err)
	// 	ResponseJSON(c, CodeFailure, err)
	// 	return
	// }
	// infos := make([]*model.SysInfo, 0)
	// if err := database.Database[*model.SysInfo]().WithLimit(-1).List(&infos); err != nil {
	// 	zap.S().Error(err)
	// 	ResponseJSON(c, CodeFailure, err)
	// 	return
	// }
	// ResponseJSON(c, CodeSuccess, infos)

	groups := make([]*model.Group, 0)
	if err := database.Database[*model.Group]().
		WithCursor("0196a0b3-c9d1-713c-870e-adc76af9f857", true).
		WithLimit(2).
		List(&groups); err != nil {
		ResponseJSON(c, CodeFailure.WithErr(err))
	}
	ResponseJSON(c, CodeSuccess, groups)
}
