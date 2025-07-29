package service_cmdb_test

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"nebula/model"
	model_cmdb "nebula/model/cmdb"

	"github.com/araddon/dateparse"
	model_internal "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"github.com/xuri/excelize/v2"
)

func TestImportDNS(t *testing.T) {
	filename := "../../testdata.nosync/6.各项目域名汇总.xlsx"

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if f.WorkBook == nil {
		t.Fatal(err)
	}
	if len(f.WorkBook.Sheets.Sheet) == 0 {
		t.Fatal(err)
	}
	sheetName := f.WorkBook.Sheets.Sheet[0].Name
	rows, err := f.GetRows(sheetName)
	if err != nil {
		t.Fatal(err)
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
			_ = cell
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
		dnsList = append(dnsList, dns)
	}
	for i, m := range dnsList {
		switch i {
		case 0, 1, 2:
			continue
		}
		fmt.Println(m.ProjectId, m.SubProjectId)
	}
	_ = dnsList
}
