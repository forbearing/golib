package model

import (
	"net/url"
	"strings"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/util"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

var MenuRoot = &Menu{ParentId: RootId, Base: Base{ID: RootId}}

type MenuPlatform string

const (
	MenuPlatformAll     = "all"
	MenuPlatformWeb     = "web"
	MenuPlatformMobile  = "mobile"
	MenuPlatformDesktop = "desktop"
)

// Menu: 菜单
// TODO: 加一个 api 用来指定后端路由,如果为空则使用 Path.
type Menu struct {
	Api     string `json:"api,omitempty" schema:"api"` // 后端路由, 如果为空则使用 "/api" + Path
	Path    string `json:"path" schema:"path"`         // path should not add `omitempty` tag, empty value means default router in react route6.x.
	Element string `json:"element,omitempty" schema:"element"`
	Label   string `json:"label,omitempty" schema:"label"`
	Icon    string `json:"icon,omitempty" schema:"icon"`

	Visiable *bool  `json:"visiable" schema:"visiable"`                                                                    // 默认路由
	Default  string `json:"default,omitempty" schema:"default"`                                                            // 自路由中的默认路由, 如果有 Children, Default 才可能存在
	Status   *uint  `json:"status" gorm:"type:smallint;default:1;comment:status(0: disabled, 1: enabled)" schema:"status"` // 该路由是否启用

	// RoleIds GormStrings `json:"role_ids,omitempty"`
	// Roles   []*Role     `json:"roles,omitempty" gorm:"-"`

	ParentId string  `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Children []*Menu `json:"children,omitempty" gorm:"foreignKey:ParentId"`             // 子路由
	Parent   *Menu   `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"` // 父路由

	// the empty value of `Platform` means all.
	Platform MenuPlatform `json:"platform" schema:"platform"`

	DomainPattern string `json:"domain_pattern" schema:"domain_pattern"`

	Base
}

func (m *Menu) Expands() []string { return []string{"Children", "Parent"} }
func (m *Menu) Excludes() map[string][]any {
	return map[string][]any{KeyId: {RootId, UnknownId, NoneId}}
}

func (m *Menu) CreateBefore() (err error) {
	return multierr.Combine(m.initDefaultValue(), m.checkPathAndApi())
}

func (m *Menu) UpdateBefore() error {
	return multierr.Combine(m.initDefaultValue(), m.checkPathAndApi())
}

// ListAfter 可能是只查询最顶层的 Menu,并不能拿到最顶层的 Menu
func (m *Menu) ListAfter() (err error) {
	oldPath, oldApi := m.Path, m.Api
	if err = m.checkPathAndApi(); err != nil {
		return err
	}
	if m.Path != oldPath || m.Api != oldApi {
		return database.Database[*Menu]().WithoutHook().Update(m)
	}
	return nil
}

func (m *Menu) GetAfter() (err error) {
	oldPath, oldApi := m.Path, m.Api
	if err = m.checkPathAndApi(); err != nil {
		return err
	}
	if m.Path != oldPath || m.Api != oldApi {
		return database.Database[*Menu]().WithoutHook().Update(m)
	}
	return nil
}

func (m *Menu) initDefaultValue() error {
	if len(m.ParentId) == 0 {
		m.ParentId = RootId
	}
	if m.Visiable == nil {
		m.Visiable = util.ValueOf(true)
	}
	if len(m.DomainPattern) == 0 {
		m.DomainPattern = ".*"
	}
	return nil
}

func (m *Menu) checkPathAndApi() (err error) {
	// 去除空格和尾部所有的 /
	m.Path = strings.TrimSpace(m.Path)
	m.Path = strings.TrimRight(m.Path, "/")

	// 检查是否是有效的 url
	var newPath string
	if newPath, err = url.JoinPath("/", m.Path); err != nil {
		return err
	}

	// m.Path 可能为空,如果这是一个父级菜单的话,则Path为空
	if len(m.Path) == 0 {
		m.Api = ""
	}
	if len(m.Path) > 0 && len(m.Api) == 0 {
		if m.Api, err = url.JoinPath("/api", m.Path); err != nil {
			return err
		}
	}

	// 有一些 path 不是以 / 开头的, 我们需要手动加上
	if len(m.Path) > 0 {
		m.Path = newPath
	}

	return nil
}

func (m *Menu) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if m == nil {
		return nil
	}
	enc.AddString("api", m.Api)
	enc.AddString("path", m.Path)
	enc.AddString("label", m.Label)
	enc.AddString("element", m.Element)
	enc.AddInt("children len", len(m.Children))

	return nil
}
