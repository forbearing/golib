package service_cmdb

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/forbearing/golib/database"
	model_internal "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	"github.com/xuri/excelize/v2"

	"nebula/model"
	model_cmdb "nebula/model/cmdb"
)

func init() {
	service.Register[*dns]()
}

type dns struct {
	service.Base[*model_cmdb.DNS]
}

func (d *dns) Import(ctx *types.ServiceContext, reader io.Reader) ([]*model_cmdb.DNS, error) {
	log := d.WithServiceContext(ctx, ctx.GetPhase())

	f, err := excelize.OpenReader(reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer f.Close()
	if f.WorkBook == nil {
		log.Error(err)
		return nil, err
	}
	if len(f.WorkBook.Sheets.Sheet) == 0 {
		log.Error(err)
		return nil, err
	}
	sheetName := f.WorkBook.Sheets.Sheet[0].Name
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var (
		expectLen int
		dnsList   = make([]*model_cmdb.DNS, 0)
	)

	for i, row := range rows {
		dns := new(model_cmdb.DNS)
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
				dns.ProjectId = cell
			case 2:
				dns.SubProjectId = cell
			case 3:
				dns.Env = model.Env(cell)
			case 4:
				dns.TenantId = cell
			case 5:
				dns.Domain = cell
			case 6:
				dns.IcpStatus = cell
			case 7:
				dns.Usage = cell
			case 8:
				dns.RegionId = cell
			case 9:
				val, _ := strconv.ParseBool(cell)
				dns.IsBuiltInLine = val
			case 10:
				val, _ := strconv.ParseBool(cell)
				dns.IsDedicatedLine = val
			case 11:
				dns.Hosts = util.ValueOf(model_internal.GormStrings(strings.Fields(cell)))
			case 12:
				dns.ServerIps = util.ValueOf(model_internal.GormStrings(strings.Fields(cell)))
			case 13:
				dns.SourceIps = util.ValueOf(model_internal.GormStrings(strings.Fields(cell)))
			case 14:
				t, _ := dateparse.ParseAny(cell)
				dns.PurchaseAt = util.ValueOf(model_internal.GormTime(t))
			case 15:
				d, _ := time.ParseDuration(cell)
				dns.PurchaseDuration = d
			case 16:
				t, _ := dateparse.ParseAny(cell)
				dns.ExpireAt = util.ValueOf(model_internal.GormTime(t))
			case 17:
				dns.VendorId = cell
			case 18:
				dns.CdnProvider = cell
			case 19:
				dns.Owner = cell
			}
		}
		switch i {
		case 0, 1: // 第一行、第二行不是数据
			continue
		default:

			// 1.必须 ProjectId 不为空
			if len(dns.ProjectId) == 0 {
				continue
			}
			// 2.创建 project, 设置 dns 的 project id
			if len(dns.ProjectId) > 0 {
				name := dns.ProjectId
				projectId, err := getOrCreateProject(name)
				if err != nil {
					log.Error(err)
					continue
				}
				dns.ProjectId = projectId
			}
			if len(dns.SubProjectId) > 0 {
				name := dns.SubProjectId
				subProjectId, err := getOrCreateProject(name)
				if err != nil {
					log.Error(err)
					continue
				}
				dns.SubProjectId = subProjectId
			}
			// 3.创建 vendor, 设置 machine 的 vendor id
			if len(dns.VendorId) > 0 {
				name := dns.VendorId
				vendorId, err := getOrCreateVendor(name)
				if err != nil {
					log.Error(err)
					continue
				}
				dns.VendorId = vendorId
			}
			// 4.创建 region, 设置 machine 的 region id
			if len(dns.RegionId) > 0 {
				name := dns.RegionId
				regionId, err := getOrCreateRegion(name)
				if err != nil {
					log.Error(err)
					continue
				}
				dns.RegionId = regionId
			}
		}
		dnsList = append(dnsList, dns)
	}

	if err := database.Database[*model_cmdb.DNS]().WithLimit(-1).Update(dnsList...); err != nil {
		log.Error(err)
		return nil, err
	}

	return nil, nil
}
