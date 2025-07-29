package model_cmdb

import (
	model_setting "nebula/model/setting"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
)

type Machine struct {
	Hostname  string             `json:"hostname,omitempty" schema:"hostname"`
	IpAddress *model.GormStrings `json:"ip_address,omitempty" schema:"ip_address"`
	Email     string             `json:"email,omitempty" schema:"email"`
	Owner     string             `json:"owner,omitempty" schema:"owner"`

	Cpu     string  `json:"cpu,omitempty" schema:"cpu"`
	Memory  string  `json:"memory,omitempty" schema:"memory"`
	Storage string  `json:"storage,omitempty" schema:"storage"`
	Network string  `json:"network,omitempty" schema:"network"`
	Cost    float64 `json:"cost,omitempty" schema:"cost"`

	DeployedApps *model.GormStrings `json:"deployed_apps,omitempty" schema:"deployed_apps"`

	ProjectId   string                 `json:"project_id,omitempty" schema:"project_id"`
	VendorId    string                 `json:"vendor_id,omitempty" schema:"vendor_id"`
	ProjectName string                 `json:"project_name,omitempty" schema:"project_name"`
	VendorName  string                 `json:"vendor_name,omitempty" schema:"vendor_name"`
	Project     *model_setting.Project `json:"project,omitempty" gorm:"-"`
	Vendor      *model_setting.Vendor  `json:"vendor,omitempty" gorm:"-"`

	model.Base
}

func (m *Machine) GetAfter() error  { return m.expendFields() }
func (m *Machine) ListAfter() error { return m.expendFields() }
func (m *Machine) expendFields() error {
	if len(m.ProjectId) > 0 {
		p := new(model_setting.Project)
		if err := database.Database[*model_setting.Project]().WithCache().Get(p, m.ProjectId); err != nil {
			return err
		}
		m.Project = p
	}

	if len(m.VendorId) > 0 {
		v := new(model_setting.Vendor)
		if err := database.Database[*model_setting.Vendor]().WithCache().Get(v, m.VendorId); err != nil {
			return err
		}
		m.Vendor = v
	}

	return nil
}
