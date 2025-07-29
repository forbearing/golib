package model_cmdb

import (
	"time"

	"nebula/model"
	model_setting "nebula/model/setting"

	model_internal "github.com/forbearing/golib/model"
)

type DNS struct {
	Env      model.Env `json:"env,omitempty" schema:"env"`
	TenantId string    `json:"tenant_id,omitempty" schema:"tenant_id"`
	Domain   string    `json:"domain,omitempty" schema:"domain"`

	// TODO: owner should be id
	Owner string `json:"owner,omitempty" schema:"owner"`

	// TODO: icp_status should be enum
	IcpStatus string `json:"icp_status,omitempty" schema:"icp_status"`

	Usage            string                      `json:"usage,omitempty" schema:"usage"`
	IsBuiltInLine    bool                        `json:"is_builtin_line,omitempty" schema:"is_builtin_line"`
	IsDedicatedLine  bool                        `json:"is_dedicated_line,omitempty" schema:"is_dedicated_line"`
	Hosts            *model_internal.GormStrings `json:"hosts,omitempty" schema:"hosts"`
	ServerIps        *model_internal.GormStrings `json:"server_ips,omitempty" schema:"server_ips"`
	SourceIps        *model_internal.GormStrings `json:"source_ips,omitempty" schema:"source_ips"`
	PurchaseAt       *model_internal.GormTime    `json:"purchase_at,omitempty" schema:"purchase_at"`
	PurchaseDuration time.Duration               `json:"purchase_duration,omitempty" schema:"purchase_duration"`
	ExpireAt         *model_internal.GormTime    `json:"expire_at,omitempty" schema:"expire_at"`
	CdnProvider      string

	ProjectId    string                 `json:"project_id,omitempty" schema:"project_id"`
	SubProjectId string                 `json:"sub_project_id,omitempty" schema:"sub_project_id"`
	RegionId     string                 `json:"region_id,omitempty" schema:"region_id"`
	VendorId     string                 `json:"vendor_id,omitempty" schema:"vendor_id"`
	Project      *model_setting.Project `json:"project,omitempty" gorm:"-"`
	SubProject   *model_setting.Project `json:"sub_project,omitempty" gorm:"-"`
	Vendor       *model_setting.Vendor  `json:"vendor,omitempty" gorm:"-"`
	Region       *model_setting.Region  `json:"region,omitempty" gorm:"-"`

	model_internal.Base
}
