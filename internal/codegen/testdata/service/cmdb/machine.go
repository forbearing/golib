package service_cmdb

import (
	"io"
	"strconv"
	"strings"

	model_cmdb "nebula/model/cmdb"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	"github.com/xuri/excelize/v2"
)

func init() {
	service.Register[*machine]()
}

type machine struct {
	service.Base[*model_cmdb.Machine]
}

// The xls file row format:
// id  project  hostname  ip_address  cpu memory  storage  network  vendor  cost  deployed_apps  account_email  owner  remark
func (m *machine) Import(ctx *types.ServiceContext, reader io.Reader) ([]*model_cmdb.Machine, error) {
	log := m.WithServiceContext(ctx, ctx.GetPhase())

	f, err := excelize.OpenReader(reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer f.Close()
	if f.WorkBook == nil {
		log.Warn("workbook is nil")
		return make([]*model_cmdb.Machine, 0), nil
	}
	if len(f.WorkBook.Sheets.Sheet) == 0 {
		log.Warn("sheet is nil")
		return make([]*model_cmdb.Machine, 0), nil
	}

	sheetName := f.WorkBook.Sheets.Sheet[0].Name
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var (
		expectLen int
		machines  = make([]*model_cmdb.Machine, 0)
	)

	for i, row := range rows {
		machine := new(model_cmdb.Machine)
		if i == 0 {
			expectLen = len(row)
		}
		fillRow(&row, expectLen)
		// fmt.Println(row, len(row))
		for j, cell := range row {
			cell = strings.TrimSpace(cell)
			switch j {
			case 0: // skip id
			case 1:
				machine.ProjectId = cell
			case 2:
				machine.Hostname = cell
			case 3:
				machine.IpAddress = util.ValueOf(model.GormStrings(strings.Split(cell, ",")))
			case 4:
				machine.Cpu = cell
			case 5:
				machine.Memory = cell
			case 6:
				machine.Storage = cell
			case 7:
				machine.Network = cell
			case 8:
				machine.VendorId = cell
			case 9:
				cost, _ := strconv.ParseFloat(cell, 64)
				machine.Cost = cost
			case 10:
				machine.DeployedApps = util.ValueOf(model.GormStrings(strings.Split(cell, ",")))
			case 11:
				machine.Email = cell
			case 12:
				machine.Owner = cell
			case 13:
				machine.Base.Remark = &cell
			}
		}
		switch i {
		case 0, 1: // 第一行和第二行不是数据
		default:
			// 1.必须有 hostname 才是有效的 machine
			if len(machine.Hostname) == 0 {
				continue
			}

			// 2.创建 project, 设置 machine 的 project id
			if len(machine.ProjectId) > 0 {
				name := machine.ProjectId
				projectId, err := getOrCreateProject(name)
				if err != nil {
					log.Error(err)
					continue
				}
				machine.ProjectId = projectId
			}

			// 3.创建 vendor, 设置 machine 的 vendor id
			if len(machine.VendorId) > 0 {
				name := machine.VendorId
				vendorId, err := getOrCreateVendor(name)
				if err != nil {
					log.Error(err)
					continue
				}
				machine.VendorId = vendorId
			}
			machines = append(machines, machine)
		}
	}

	if err := database.Database[*model_cmdb.Machine]().WithLimit(-1).Update(machines...); err != nil {
		log.Error(err)
		return nil, err
	}

	return machines, nil
}
