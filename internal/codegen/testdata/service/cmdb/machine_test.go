package service_cmdb_test

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	model_cmdb "nebula/model/cmdb"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"github.com/xuri/excelize/v2"
)

func TestMachineImport(t *testing.T) {
	filename := "../../testdata.nosync/5.各项目机器汇总.xlsx"

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
			_ = cell
			switch j {
			case 0: // skip id
			case 1:
				machine.ProjectId = cell
			case 2:
				machine.Hostname = cell
			case 3:
				// TODO: IpAddress should slice.
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
		machines = append(machines, machine)
	}
	for i, m := range machines {
		switch i {
		case 0, 1, 2:
			fmt.Println(m.ProjectId, m.Hostname, m.IpAddress)
		}
		// fmt.Printf("%+v\n", m)
	}
	_ = machines
}

func fillRow(row *[]string, expectLen int) {
	for i := len(*row); i < expectLen; i++ {
		*row = append(*row, "")
	}
}
