package service_cmdb

import (
	model_setting "nebula/model/setting"

	"github.com/forbearing/golib/database"
)

func getOrCreateProject(projectName string) (string, error) {
	pl := make([]*model_setting.Project, 0)
	if err := database.Database[*model_setting.Project]().WithLimit(1).WithQuery(&model_setting.Project{Name: projectName}).List(&pl); err != nil {
		return "", err
	}
	// 不存在则创建
	if len(pl) == 0 {
		p := &model_setting.Project{Name: projectName}
		if err := database.Database[*model_setting.Project]().Create(p); err != nil {
			return "", err
		}
		return p.ID, nil
	} else {
		return pl[0].ID, nil
	}
}

func getOrCreateVendor(vendorName string) (string, error) {
	vl := make([]*model_setting.Vendor, 0)
	if err := database.Database[*model_setting.Vendor]().WithLimit(1).WithQuery(&model_setting.Vendor{Name: vendorName}).List(&vl); err != nil {
		return "", err
	}
	// 不存在则创建
	if len(vl) == 0 {
		v := &model_setting.Vendor{Name: vendorName}
		if err := database.Database[*model_setting.Vendor]().Create(v); err != nil {
			return "", err
		}
		return v.ID, nil
	} else {
		return vl[0].ID, nil
	}
}

func getOrCreateRegion(regionName string) (string, error) {
	rl := make([]*model_setting.Region, 0)
	if err := database.Database[*model_setting.Region]().WithLimit(1).WithQuery(&model_setting.Region{Name: regionName}).List(&rl); err != nil {
		return "", err
	}
	// 不存在则创建
	if len(rl) == 0 {
		r := &model_setting.Region{Name: regionName}
		if err := database.Database[*model_setting.Region]().Create(r); err != nil {
			return "", err
		}
		return r.ID, nil
	} else {
		return rl[0].ID, nil
	}
}
