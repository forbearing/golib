package model_keycloak

import (
	"github.com/forbearing/golib/model"
	"gorm.io/datatypes"
)

type User struct {
	// 常用索引字段
	Username      *string `json:"username,omitempty"`
	Email         *string `json:"email,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
	EmailVerified *bool   `json:"email_verified,omitempty"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`

	RealmId   *string `json:"realm_id,omitempty"`
	RealmName *string `json:"realm_name,omitempty"`

	// 完整的 User 配置存储在 JSON 字段中
	Data datatypes.JSON `json:"data,omitempty"`

	model.Base
}

type UserData struct {
	ID                         string                     `json:"id"`
	CreatedTimestamp           int64                      `json:"created_timestamp"`
	Username                   string                     `json:"username"`
	Enabled                    bool                       `json:"enabled"`
	Totp                       bool                       `json:"totp"`
	EmailVerified              bool                       `json:"email_verified"`
	FirstName                  string                     `json:"first_name"`
	LastName                   string                     `json:"last_name"`
	Email                      string                     `json:"email"`
	FederationLink             string                     `json:"federation_link"`
	Attributes                 map[string][]string        `json:"attributes"`
	DisableableCredentialTypes []interface{}              `json:"disableable_credential_types"`
	RequiredActions            []string                   `json:"required_actions"`
	Access                     map[string]bool            `json:"access"`
	ClientRoles                map[string][]string        `json:"client_roles"`
	RealmRoles                 []string                   `json:"realm_roles"`
	Groups                     []string                   `json:"groups"`
	ServiceAccountClientID     string                     `json:"service_account_client_id"`
	Credentials                []CredentialRepresentation `json:"credentials"`
}

// type User struct {
// 	// RealmID 是用户所属的 Realm
// 	// 与 Username 组合形成唯一索引
// 	RealmID string `json:"realm_id,omitempty" gorm:"type:varchar(36);uniqueIndex:idx_realm_username,priority:1;not null"`
//
// 	// CreatedTimestamp 是用户账号创建时间的 Unix 时间戳（毫秒）
// 	// 由系统自动生成，记录用户首次创建的精确时间
// 	// 示例: 1623456789000 (代表 2021-06-12 10:13:09)
// 	CreatedTimestamp *int64 `json:"created_timestamp,omitempty" gorm:"type:bigint"`
//
// 	// Username 是用户登录时使用的用户名
// 	// 在同一个 realm 中必须唯一，可以是邮箱、手机号或自定义字符串
// 	// 示例: "john.doe" 或 "john@example.com"
// 	Username *string `json:"username,omitempty" gorm:"type:varchar(255);uniqueIndex:idx_realm_username,priority:2"`
//
// 	// Enabled 表示用户账号是否启用
// 	// true: 账号正常，可以登录; false: 账号被禁用，无法登录
// 	// 管理员可以通过此字段临时或永久禁用用户访问
// 	Enabled *bool `json:"enabled,omitempty" gorm:"type:boolean;default:true"`
//
// 	// Totp 表示用户是否启用了基于时间的一次性密码（Time-based One-Time Password）
// 	// true: 已配置 TOTP 双因素认证（如 Google Authenticator）
// 	// false: 未启用双因素认证
// 	Totp *bool `json:"totp,omitempty" gorm:"type:boolean;default:false"`
//
// 	// EmailVerified 表示用户的电子邮件地址是否已通过验证
// 	// true: 邮箱已验证; false: 邮箱未验证
// 	// 某些操作可能需要邮箱验证后才能执行
// 	EmailVerified *bool `json:"email_verified,omitempty" gorm:"type:boolean;default:false"`
//
// 	// FirstName 是用户的名字/名
// 	// 用于显示和个性化，不用于身份验证
// 	// 示例: "John"
// 	FirstName *string `json:"first_name,omitempty" gorm:"type:varchar(255)"`
//
// 	// LastName 是用户的姓氏/姓
// 	// 用于显示和个性化，不用于身份验证
// 	// 示例: "Doe"
// 	LastName *string `json:"last_name,omitempty" gorm:"type:varchar(255)"`
//
// 	// Email 是用户的电子邮件地址
// 	// 可用于登录（如果启用）、密码重置和通知
// 	// 示例: "john.doe@example.com"
// 	Email *string `json:"email,omitempty" gorm:"type:varchar(255);index"`
//
// 	// FederationLink 是外部身份提供者的链接标识
// 	// 当用户通过 LDAP、Active Directory、SAML 或其他 IdP 同步/导入时设置
// 	// 包含外部系统中用户的唯一标识符
// 	// 示例: "f:d5f8e9a0-1234-5678-9abc-def012345678:john.doe"
// 	FederationLink *string `json:"federation_link,omitempty" gorm:"type:varchar(255)"`
//
// 	// Attributes 是用户的自定义属性映射
// 	// 可存储任意额外信息，键是属性名，值是字符串数组（支持多值属性）
// 	// 示例: {"phone": ["123-456-7890"], "department": ["IT", "Security"], "employee_id": ["EMP001"]}
// 	Attributes datatypes.JSONType[map[string][]string] `json:"attributes,omitempty" gorm:"type:json"`
//
// 	// DisableableCredentialTypes 是可以被禁用的凭证类型列表
// 	// 显示哪些认证方式可以为该用户禁用
// 	// 常见类型: "password", "otp", "webauthn"
// 	DisableableCredentialTypes datatypes.JSONType[[]any] `json:"disableable_credential_types,omitempty" gorm:"type:json"`
//
// 	// RequiredActions 是用户下次登录时必须执行的操作列表
// 	// 完成这些操作前，用户无法正常使用系统
// 	// 常见值: "VERIFY_EMAIL"（验证邮箱）, "UPDATE_PASSWORD"（更新密码）,
// 	// "UPDATE_PROFILE"（更新资料）, "CONFIGURE_TOTP"（配置双因素认证）
// 	RequiredActions datatypes.JSONType[[]string] `json:"required_actions,omitempty" gorm:"type:json"`
//
// 	// Access 是用户的访问权限映射
// 	// 定义用户对自身账号的管理权限
// 	// 示例: {"manageGroupMembership": true, "view": true, "mapRoles": false, "impersonate": false}
// 	Access datatypes.JSONType[map[string]bool] `json:"access,omitempty" gorm:"type:json"`
//
// 	// ClientRoles 是用户在特定客户端应用中的角色映射
// 	// 键是客户端ID，值是该客户端中的角色列表
// 	// 这些角色仅在访问对应客户端时有效
// 	// 示例: {"account": ["manage-account", "view-profile"], "admin-cli": ["admin"]}
// 	ClientRoles datatypes.JSONType[map[string][]string] `json:"client_roles,omitempty" gorm:"type:json"`
//
// 	// RealmRoles 是用户在 realm 级别拥有的角色列表
// 	// 这些是全局角色，对整个 realm 中的所有客户端都有效
// 	// 示例: ["user", "admin", "offline_access"]
// 	RealmRoles datatypes.JSONType[[]string] `json:"realm_roles,omitempty" gorm:"type:json"`
//
// 	// Groups 是用户所属的组 ID 列表
// 	// 组可用于批量管理用户权限、角色和属性
// 	// 用户会继承所属组的角色和属性
// 	// 示例: ["/employees", "/employees/developers", "/partners"]
// 	Groups datatypes.JSONType[[]string] `json:"groups,omitempty" gorm:"type:json"`
//
// 	// ServiceAccountClientID 标识这是否是一个服务账号用户
// 	// 如果非空，表示这是代表客户端应用的服务账号，值为关联的客户端ID
// 	// 服务账号用于客户端凭证流（client credentials flow）
// 	// 普通用户此字段为 nil
// 	ServiceAccountClientID *string `json:"service_account_client_id,omitempty" gorm:"type:varchar(255);index"`
//
// 	// Credentials 是用户的凭证列表
// 	// 包含密码、OTP 设备、WebAuthn 设备等认证信息
// 	// 每个凭证包含类型、值、是否临时、创建时间等详细信息
// 	// 注意: 获取用户时通常不返回凭证的敏感信息（如密码值）
// 	Credentials datatypes.JSONType[[]CredentialRepresentation] `json:"credentials,omitempty" gorm:"type:json"`
//
// 	model.Base
// }

// CredentialRepresentation is a representations of the credentials
// v7: https://www.keycloak.org/docs-api/7.0/rest-api/index.html#_credentialrepresentation
// v8: https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_credentialrepresentation
type CredentialRepresentation struct {
	// Common part
	CreatedDate *int64  `json:"createdDate,omitempty"`
	Temporary   *bool   `json:"temporary,omitempty"`
	Type        *string `json:"type,omitempty"`
	Value       *string `json:"value,omitempty"`

	// <= v7
	Algorithm         *string                                 `json:"algorithm,omitempty"`
	Config            datatypes.JSONType[*MultiValuedHashMap] `json:"config,omitempty"`
	Counter           *int32                                  `json:"counter,omitempty"`
	Device            *string                                 `json:"device,omitempty"`
	Digits            *int32                                  `json:"digits,omitempty"`
	HashIterations    *int32                                  `json:"hash_iterations,omitempty"`
	HashedSaltedValue *string                                 `json:"hashed_salted_value,omitempty"`
	Period            *int32                                  `json:"period,omitempty"`
	Salt              *string                                 `json:"salt,omitempty"`

	// >= v8
	CredentialData *string `json:"credentialData,omitempty"`
	ID             *string `json:"id,omitempty"`
	Priority       *int32  `json:"priority,omitempty"`
	SecretData     *string `json:"secretData,omitempty"`
	UserLabel      *string `json:"userLabel,omitempty"`
}

// MultiValuedHashMap represents something
type MultiValuedHashMap struct {
	Empty      *bool    `json:"empty,omitempty"`
	LoadFactor *float32 `json:"load_factor,omitempty"`
	Threshold  *int32   `json:"threshold,omitempty"`
}

func (*User) TableName() string {
	return "keycloak_users"
}
