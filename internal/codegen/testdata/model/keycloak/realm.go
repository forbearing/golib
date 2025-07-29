package model_keycloak

import (
	"github.com/forbearing/golib/model"
	"gorm.io/datatypes"
)

// Realm 数据库模型 - 简化版本
type Realm struct {
	Name        *string `json:"name" gorm:"index"`
	DisplayName *string `json:"display_name" gorm:"index"`
	Enabled     *bool   `json:"enabled" gorm:"index"`

	// 完整的 Realm 配置存储在 JSON 字段中
	Data datatypes.JSON `json:"data,omitempty"`

	model.Base
}

// // Realm represents a realm
// type Realm struct {
// 	// AccessCodeLifespan 访问代码的生命周期（秒）
// 	// 授权码的有效期，用户必须在此时间内完成授权流程
// 	// 默认值通常为 60 秒
// 	AccessCodeLifespan *int `json:"access_code_lifespan,omitempty" gorm:"type:int"`
//
// 	// AccessCodeLifespanLogin 登录操作中访问代码的生命周期（秒）
// 	// 专门用于登录流程的授权码有效期
// 	// 可以与标准访问代码生命周期不同
// 	AccessCodeLifespanLogin *int `json:"access_code_lifespan_login,omitempty" gorm:"type:int"`
//
// 	// AccessCodeLifespanUserAction 用户操作中访问代码的生命周期（秒）
// 	// 用于需要用户交互的操作（如重置密码）的代码有效期
// 	// 通常比标准访问代码更长
// 	AccessCodeLifespanUserAction *int `json:"access_code_lifespan_user_action,omitempty" gorm:"type:int"`
//
// 	// AccessTokenLifespan 访问令牌的生命周期（秒）
// 	// JWT 访问令牌的有效期
// 	// 默认值通常为 300 秒（5分钟）
// 	AccessTokenLifespan *int `json:"access_token_lifespan,omitempty" gorm:"type:int"`
//
// 	// AccessTokenLifespanForImplicitFlow 隐式流程中访问令牌的生命周期（秒）
// 	// 专门用于隐式授权流程的访问令牌有效期
// 	// 由于安全考虑，通常比标准流程更短
// 	AccessTokenLifespanForImplicitFlow *int `json:"access_token_lifespan_for_implicit_flow,omitempty" gorm:"type:int"`
//
// 	// AccountTheme 账户管理界面的主题
// 	// 用户自助服务界面的视觉主题
// 	// 示例: "keycloak", "custom-theme"
// 	AccountTheme *string `json:"account_theme,omitempty" gorm:"type:varchar(255)"`
//
// 	// ActionTokenGeneratedByAdminLifespan 管理员生成的操作令牌生命周期（秒）
// 	// 管理员为用户生成的操作令牌（如密码重置链接）的有效期
// 	// 默认值通常为 43200 秒（12小时）
// 	ActionTokenGeneratedByAdminLifespan *int `json:"action_token_generated_by_admin_lifespan,omitempty" gorm:"type:int"`
//
// 	// ActionTokenGeneratedByUserLifespan 用户生成的操作令牌生命周期（秒）
// 	// 用户自己触发的操作令牌（如邮箱验证）的有效期
// 	// 默认值通常为 300 秒（5分钟）
// 	ActionTokenGeneratedByUserLifespan *int `json:"action_token_generated_by_user_lifespan,omitempty" gorm:"type:int"`
//
// 	// AdminEventsDetailsEnabled 是否启用管理事件详情
// 	// true: 记录管理操作的详细信息
// 	// false: 只记录基本信息
// 	AdminEventsDetailsEnabled *bool `json:"admin_events_details_enabled,omitempty" gorm:"type:boolean;default:false"`
//
// 	// AdminEventsEnabled 是否启用管理事件记录
// 	// true: 记录所有管理操作（创建用户、修改配置等）
// 	// false: 不记录管理事件
// 	AdminEventsEnabled *bool `json:"admin_events_enabled,omitempty" gorm:"type:boolean;default:false"`
//
// 	// AdminTheme 管理控制台的主题
// 	// Keycloak 管理界面的视觉主题
// 	// 示例: "keycloak", "custom-admin-theme"
// 	AdminTheme *string `json:"admin_theme,omitempty" gorm:"type:varchar(255)"`
//
// 	// Attributes realm 的自定义属性
// 	// 可存储任意配置信息
// 	// 示例: {"custom_config": "value", "feature_flag": "enabled"}
// 	Attributes datatypes.JSONType[map[string]string] `json:"attributes,omitempty" gorm:"type:json"`
//
// 	// AuthenticationFlows 认证流程配置
// 	// 定义各种认证场景的流程
// 	// 包含浏览器登录、直接授权等流程的详细步骤
// 	AuthenticationFlows datatypes.JSON `json:"authentication_flows,omitempty" gorm:"type:json"`
//
// 	// AuthenticatorConfig 认证器配置
// 	// 各种认证器（OTP、WebAuthn等）的具体配置
// 	AuthenticatorConfig datatypes.JSON `json:"authenticator_config,omitempty" gorm:"type:json"`
//
// 	// BrowserFlow 浏览器认证流程名称
// 	// 标准浏览器登录使用的认证流程
// 	// 默认值: "browser"
// 	BrowserFlow *string `json:"browser_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// BrowserSecurityHeaders 浏览器安全头配置
// 	// 设置各种 HTTP 安全响应头
// 	// 示例: {"X-Frame-Options": "SAMEORIGIN", "Content-Security-Policy": "frame-src 'self'"}
// 	BrowserSecurityHeaders datatypes.JSONType[map[string]string] `json:"browser_security_headers,omitempty" gorm:"type:json"`
//
// 	// BruteForceProtected 是否启用暴力破解保护
// 	// true: 启用登录失败次数限制和账户锁定
// 	// false: 不限制登录尝试次数
// 	BruteForceProtected *bool `json:"brute_force_protected,omitempty" gorm:"type:boolean;default:false"`
//
// 	// ClientAuthenticationFlow 客户端认证流程名称
// 	// 客户端进行身份验证时使用的流程
// 	// 默认值: "clients"
// 	ClientAuthenticationFlow *string `json:"client_authentication_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// ClientPolicies 客户端策略配置
// 	// 定义客户端必须遵守的策略
// 	// 用于强制执行安全标准
// 	ClientPolicies datatypes.JSON `json:"client_policies,omitempty" gorm:"type:json"`
//
// 	// ClientProfiles 客户端配置文件
// 	// 定义客户端的配置模板
// 	// 可以被多个客户端共享
// 	ClientProfiles datatypes.JSON `json:"client_profiles,omitempty" gorm:"type:json"`
//
// 	// ClientScopeMappings 客户端作用域映射
// 	// 定义客户端之间的作用域映射关系
// 	ClientScopeMappings datatypes.JSON `json:"client_scope_mappings,omitempty" gorm:"type:json"`
//
// 	// ClientScopes 客户端作用域列表
// 	// 定义可用的客户端作用域
// 	// 用于管理客户端可以访问的资源范围
// 	ClientScopes datatypes.JSONType[[]ClientScope] `json:"client_scopes,omitempty" gorm:"type:json"`
//
// 	// Clients 客户端列表
// 	// realm 中注册的所有客户端应用
// 	Clients datatypes.JSONType[[]Client] `json:"clients,omitempty" gorm:"type:json"`
//
// 	// Components 组件配置
// 	// 各种可插拔组件的配置（存储提供者、密钥提供者等）
// 	Components datatypes.JSON `json:"components,omitempty" gorm:"type:json"`
//
// 	// DefaultDefaultClientScopes 默认的默认客户端作用域
// 	// 新客户端自动获得的默认作用域
// 	// 示例: ["openid", "profile", "email", "web-origins", "roles"]
// 	DefaultDefaultClientScopes datatypes.JSONType[[]string] `json:"default_default_client_scopes,omitempty" gorm:"type:json"`
//
// 	// DefaultGroups 默认组列表
// 	// 新用户自动加入的组
// 	// 示例: ["/default-users", "/newsletter-subscribers"]
// 	DefaultGroups datatypes.JSONType[[]string] `json:"default_groups,omitempty" gorm:"type:json"`
//
// 	// DefaultLocale 默认语言环境
// 	// realm 的默认语言设置
// 	// 示例: "en", "zh-CN", "ja"
// 	DefaultLocale *string `json:"default_locale,omitempty" gorm:"type:varchar(255)"`
//
// 	// DefaultOptionalClientScopes 默认的可选客户端作用域
// 	// 新客户端可以选择的作用域
// 	// 示例: ["address", "phone", "offline_access"]
// 	DefaultOptionalClientScopes datatypes.JSONType[[]string] `json:"default_optional_client_scopes,omitempty" gorm:"type:json"`
//
// 	// DefaultRole 默认角色对象
// 	// 包含默认角色的详细配置
// 	DefaultRole datatypes.JSONType[*Role] `json:"default_role,omitempty" gorm:"type:json"`
//
// 	// DefaultRoles 默认角色名称列表（已弃用）
// 	// 新用户自动获得的角色
// 	// 注意: 推荐使用 DefaultRole 字段
// 	DefaultRoles datatypes.JSONType[[]string] `json:"default_roles,omitempty" gorm:"type:json"`
//
// 	// DefaultSignatureAlgorithm 默认签名算法
// 	// 用于签名令牌的算法
// 	// 示例: "RS256", "ES256", "HS256"
// 	DefaultSignatureAlgorithm *string `json:"default_signature_algorithm,omitempty" gorm:"type:varchar(255)"`
//
// 	// DirectGrantFlow 直接授权流程名称
// 	// 资源所有者密码凭证流使用的认证流程
// 	// 默认值: "direct grant"
// 	DirectGrantFlow *string `json:"direct_grant_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// DisplayName realm 的显示名称
// 	// 在用户界面中显示的友好名称
// 	// 示例: "公司员工系统", "客户门户"
// 	DisplayName *string `json:"display_name,omitempty" gorm:"type:varchar(255)"`
//
// 	// DisplayNameHTML realm 显示名称的 HTML 版本
// 	// 支持 HTML 格式的显示名称
// 	// 可包含样式和格式化
// 	DisplayNameHTML *string `json:"display_name_html,omitempty" gorm:"type:text"`
//
// 	// DockerAuthenticationFlow Docker 认证流程名称
// 	// Docker 客户端认证使用的流程
// 	// 用于 Docker Registry 集成
// 	DockerAuthenticationFlow *string `json:"docker_authentication_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// DuplicateEmailsAllowed 是否允许重复邮箱
// 	// true: 多个用户可以使用相同邮箱
// 	// false: 邮箱必须唯一
// 	DuplicateEmailsAllowed *bool `json:"duplicate_emails_allowed,omitempty" gorm:"type:boolean;default:false"`
//
// 	// EditUsernameAllowed 是否允许编辑用户名
// 	// true: 用户可以修改自己的用户名
// 	// false: 用户名一旦设置不可更改
// 	EditUsernameAllowed *bool `json:"edit_username_allowed,omitempty" gorm:"type:boolean;default:false"`
//
// 	// EmailTheme 邮件主题
// 	// 发送邮件时使用的模板主题
// 	// 示例: "keycloak", "custom-email-theme"
// 	EmailTheme *string `json:"email_theme,omitempty" gorm:"type:varchar(255)"`
//
// 	// Enabled realm 是否启用
// 	// true: realm 正常运行
// 	// false: realm 被禁用，所有请求都会被拒绝
// 	Enabled *bool `json:"enabled,omitempty" gorm:"type:boolean;default:true"`
//
// 	// EnabledEventTypes 启用的事件类型列表
// 	// 需要记录的事件类型
// 	// 示例: ["LOGIN", "LOGIN_ERROR", "LOGOUT", "REGISTER"]
// 	EnabledEventTypes datatypes.JSONType[[]string] `json:"enabled_event_types,omitempty" gorm:"type:json"`
//
// 	// EventsEnabled 是否启用事件记录
// 	// true: 记录用户活动事件
// 	// false: 不记录事件
// 	EventsEnabled *bool `json:"events_enabled,omitempty" gorm:"type:boolean;default:false"`
//
// 	// EventsExpiration 事件过期时间（秒）
// 	// 事件记录保留的时间
// 	// 0 表示永不过期
// 	EventsExpiration *int64 `json:"events_expiration,omitempty" gorm:"type:bigint"`
//
// 	// EventsListeners 事件监听器列表
// 	// 处理事件的监听器
// 	// 示例: ["jboss-logging", "email", "custom-listener"]
// 	EventsListeners datatypes.JSONType[[]string] `json:"events_listeners,omitempty" gorm:"type:json"`
//
// 	// FailureFactor 失败因子
// 	// 暴力破解保护中，触发锁定的失败次数
// 	// 默认值通常为 30
// 	FailureFactor *int `json:"failure_factor,omitempty" gorm:"type:int"`
//
// 	// FederatedUsers 联邦用户列表
// 	// 从外部身份提供者导入的用户
// 	FederatedUsers datatypes.JSON `json:"federated_users,omitempty" gorm:"type:json"`
//
// 	// Groups 组列表
// 	// realm 中定义的所有用户组
// 	Groups datatypes.JSON `json:"groups,omitempty" gorm:"type:json"`
//
// 	// ID realm 的唯一标识符
// 	// 由 Keycloak 自动生成的 UUID
// 	// 示例: "master", "my-realm"
// 	ID *string `json:"id,omitempty" gorm:"type:varchar(36);primaryKey"`
//
// 	// IdentityProviderMappers 身份提供者映射器
// 	// 定义如何映射外部 IdP 的属性到 Keycloak 用户
// 	IdentityProviderMappers datatypes.JSON `json:"identity_provider_mappers,omitempty" gorm:"type:json"`
//
// 	// IdentityProviders 身份提供者列表
// 	// 配置的外部身份提供者（Google、Facebook、SAML等）
// 	IdentityProviders datatypes.JSON `json:"identity_providers,omitempty" gorm:"type:json"`
//
// 	// InternationalizationEnabled 是否启用国际化
// 	// true: 支持多语言
// 	// false: 只使用默认语言
// 	InternationalizationEnabled *bool `json:"internationalization_enabled,omitempty" gorm:"type:boolean;default:false"`
//
// 	// KeycloakVersion Keycloak 版本
// 	// 创建此 realm 的 Keycloak 版本
// 	// 示例: "21.0.0"
// 	KeycloakVersion *string `json:"keycloak_version,omitempty" gorm:"type:varchar(255)"`
//
// 	// LoginTheme 登录页面主题
// 	// 登录界面使用的视觉主题
// 	// 示例: "keycloak", "custom-login-theme"
// 	LoginTheme *string `json:"login_theme,omitempty" gorm:"type:varchar(255)"`
//
// 	// LoginWithEmailAllowed 是否允许使用邮箱登录
// 	// true: 用户可以使用邮箱或用户名登录
// 	// false: 只能使用用户名登录
// 	LoginWithEmailAllowed *bool `json:"login_with_email_allowed,omitempty" gorm:"type:boolean;default:true"`
//
// 	// MaxDeltaTimeSeconds 最大时间偏差（秒）
// 	// 客户端和服务器之间允许的最大时间差
// 	// 用于防止重放攻击
// 	MaxDeltaTimeSeconds *int `json:"max_delta_time_seconds,omitempty" gorm:"type:int"`
//
// 	// MaxFailureWaitSeconds 最大失败等待时间（秒）
// 	// 账户锁定的最长时间
// 	// 默认值通常为 900 秒（15分钟）
// 	MaxFailureWaitSeconds *int `json:"max_failure_wait_seconds,omitempty" gorm:"type:int"`
//
// 	// MinimumQuickLoginWaitSeconds 最小快速登录等待时间（秒）
// 	// 防止暴力破解的最小等待时间
// 	// 默认值通常为 60 秒
// 	MinimumQuickLoginWaitSeconds *int `json:"minimum_quick_login_wait_seconds,omitempty" gorm:"type:int"`
//
// 	// NotBefore 生效时间戳
// 	// 在此时间之前签发的令牌都被视为无效
// 	// 用于令牌撤销
// 	NotBefore *int `json:"not_before,omitempty" gorm:"type:int"`
//
// 	// OfflineSessionIdleTimeout 离线会话空闲超时（秒）
// 	// 离线会话在空闲多久后过期
// 	// 默认值通常为 2592000 秒（30天）
// 	OfflineSessionIdleTimeout *int `json:"offline_session_idle_timeout,omitempty" gorm:"type:int"`
//
// 	// OfflineSessionMaxLifespan 离线会话最大生命周期（秒）
// 	// 离线会话的最长有效期
// 	// 默认值通常为 5184000 秒（60天）
// 	OfflineSessionMaxLifespan *int `json:"offline_session_max_lifespan,omitempty" gorm:"type:int"`
//
// 	// OfflineSessionMaxLifespanEnabled 是否启用离线会话最大生命周期
// 	// true: 强制执行最大生命周期限制
// 	// false: 离线会话永不过期
// 	OfflineSessionMaxLifespanEnabled *bool `json:"offline_session_max_lifespan_enabled,omitempty" gorm:"type:boolean;default:false"`
//
// 	// OtpPolicyAlgorithm OTP 策略算法
// 	// 生成 OTP 的算法
// 	// 示例: "HmacSHA1", "HmacSHA256", "HmacSHA512"
// 	OtpPolicyAlgorithm *string `json:"otp_policy_algorithm,omitempty" gorm:"type:varchar(255)"`
//
// 	// OtpPolicyDigits OTP 位数
// 	// OTP 代码的位数
// 	// 通常为 6 或 8
// 	OtpPolicyDigits *int `json:"otp_policy_digits,omitempty" gorm:"type:int"`
//
// 	// OtpPolicyInitialCounter OTP 初始计数器
// 	// HOTP（基于计数器的OTP）的初始值
// 	// 默认为 0
// 	OtpPolicyInitialCounter *int `json:"otp_policy_initial_counter,omitempty" gorm:"type:int"`
//
// 	// OtpPolicyLookAheadWindow OTP 前瞻窗口
// 	// 接受的 OTP 偏移范围
// 	// 用于处理时间不同步
// 	OtpPolicyLookAheadWindow *int `json:"otp_policy_look_ahead_window,omitempty" gorm:"type:int"`
//
// 	// OtpPolicyPeriod OTP 周期（秒）
// 	// TOTP 代码的有效期
// 	// 默认为 30 秒
// 	OtpPolicyPeriod *int `json:"otp_policy_period,omitempty" gorm:"type:int"`
//
// 	// OtpPolicyType OTP 策略类型
// 	// OTP 的类型
// 	// 值: "totp"（基于时间）或 "hotp"（基于计数器）
// 	OtpPolicyType *string `json:"otp_policy_type,omitempty" gorm:"type:varchar(255)"`
//
// 	// OtpSupportedApplications 支持的 OTP 应用列表
// 	// 推荐给用户的 OTP 应用
// 	// 示例: ["FreeOTP", "Google Authenticator"]
// 	OtpSupportedApplications datatypes.JSONType[[]string] `json:"otp_supported_applications,omitempty" gorm:"type:json"`
//
// 	// PasswordPolicy 密码策略
// 	// 定义密码复杂度要求
// 	// 示例: "length(8) and upperCase(1) and lowerCase(1) and digits(1) and specialChars(1)"
// 	PasswordPolicy *string `json:"password_policy,omitempty" gorm:"type:text"`
//
// 	// PermanentLockout 是否永久锁定
// 	// true: 达到失败次数后永久锁定账户
// 	// false: 临时锁定，过期后自动解锁
// 	PermanentLockout *bool `json:"permanent_lockout,omitempty" gorm:"type:boolean;default:false"`
//
// 	// ProtocolMappers 协议映射器
// 	// 定义如何映射用户信息到令牌
// 	ProtocolMappers datatypes.JSON `json:"protocol_mappers,omitempty" gorm:"type:json"`
//
// 	// QuickLoginCheckMilliSeconds 快速登录检查间隔（毫秒）
// 	// 防止暴力破解的检查间隔
// 	// 默认值通常为 1000 毫秒
// 	QuickLoginCheckMilliSeconds *int64 `json:"quick_login_check_milli_seconds,omitempty" gorm:"type:bigint"`
//
// 	// Realm realm 名称
// 	// realm 的唯一标识名称
// 	// 示例: "master", "production", "development"
// 	Realm *string `json:"realm,omitempty" gorm:"type:varchar(255);uniqueIndex"`
//
// 	// RefreshTokenMaxReuse 刷新令牌最大重用次数
// 	// 刷新令牌可以被使用的最大次数
// 	// 0 表示无限制
// 	RefreshTokenMaxReuse *int `json:"refresh_token_max_reuse,omitempty" gorm:"type:int"`
//
// 	// RegistrationAllowed 是否允许用户注册
// 	// true: 开放用户自助注册
// 	// false: 只能由管理员创建用户
// 	RegistrationAllowed *bool `json:"registration_allowed,omitempty" gorm:"type:boolean;default:false"`
//
// 	// RegistrationEmailAsUsername 注册时是否使用邮箱作为用户名
// 	// true: 自动使用邮箱地址作为用户名
// 	// false: 用户需要单独设置用户名
// 	RegistrationEmailAsUsername *bool `json:"registration_email_as_username,omitempty" gorm:"type:boolean;default:false"`
//
// 	// RegistrationFlow 注册流程名称
// 	// 用户注册时使用的认证流程
// 	// 默认值: "registration"
// 	RegistrationFlow *string `json:"registration_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// RememberMe 是否启用"记住我"功能
// 	// true: 登录页面显示"记住我"选项
// 	// false: 不提供记住我功能
// 	RememberMe *bool `json:"remember_me,omitempty" gorm:"type:boolean;default:false"`
//
// 	// RequiredActions 必需操作配置
// 	// 定义各种必需操作的配置
// 	RequiredActions datatypes.JSON `json:"required_actions,omitempty" gorm:"type:json"`
//
// 	// ResetCredentialsFlow 重置凭证流程名称
// 	// 密码重置使用的认证流程
// 	// 默认值: "reset credentials"
// 	ResetCredentialsFlow *string `json:"reset_credentials_flow,omitempty" gorm:"type:varchar(255)"`
//
// 	// ResetPasswordAllowed 是否允许重置密码
// 	// true: 用户可以自助重置密码
// 	// false: 只能由管理员重置密码
// 	ResetPasswordAllowed *bool `json:"reset_password_allowed,omitempty" gorm:"type:boolean;default:false"`
//
// 	// RevokeRefreshToken 是否撤销刷新令牌
// 	// true: 使用后立即撤销刷新令牌
// 	// false: 刷新令牌可重复使用
// 	RevokeRefreshToken *bool `json:"revoke_refresh_token,omitempty" gorm:"type:boolean;default:false"`
//
// 	// Roles 角色配置
// 	// realm 和客户端角色的完整配置
// 	Roles datatypes.JSONType[*RolesRepresentation] `json:"roles,omitempty" gorm:"type:json"`
//
// 	// ScopeMappings 作用域映射
// 	// 定义角色和作用域的映射关系
// 	ScopeMappings datatypes.JSON `json:"scope_mappings,omitempty" gorm:"type:json"`
//
// 	// SMTPServer SMTP 服务器配置
// 	// 邮件发送配置
// 	// 示例: {"host": "smtp.gmail.com", "port": "587", "from": "noreply@example.com"}
// 	SMTPServer datatypes.JSONType[map[string]string] `json:"smtp_server,omitempty" gorm:"type:json"`
//
// 	// SslRequired SSL 要求级别
// 	// 定义 SSL/HTTPS 的要求程度
// 	// 值: "all"（所有请求）, "external"（外部请求）, "none"（不要求）
// 	SslRequired *string `json:"ssl_required,omitempty" gorm:"type:varchar(255)"`
//
// 	// SsoSessionIdleTimeout SSO 会话空闲超时（秒）
// 	// 用户无活动后会话过期时间
// 	// 默认值通常为 1800 秒（30分钟）
// 	SsoSessionIdleTimeout *int `json:"sso_session_idle_timeout,omitempty" gorm:"type:int"`
//
// 	// SsoSessionIdleTimeoutRememberMe 记住我状态下的 SSO 会话空闲超时（秒）
// 	// 启用"记住我"时的会话空闲超时
// 	// 通常比标准超时更长
// 	SsoSessionIdleTimeoutRememberMe *int `json:"sso_session_idle_timeout_remember_me,omitempty" gorm:"type:int"`
//
// 	// SsoSessionMaxLifespan SSO 会话最大生命周期（秒）
// 	// 会话的最长有效期，无论是否活跃
// 	// 默认值通常为 36000 秒（10小时）
// 	SsoSessionMaxLifespan *int `json:"sso_session_max_lifespan,omitempty" gorm:"type:int"`
//
// 	// SsoSessionMaxLifespanRememberMe 记住我状态下的 SSO 会话最大生命周期（秒）
// 	// 启用"记住我"时的会话最大生命周期
// 	// 通常比标准生命周期更长
// 	SsoSessionMaxLifespanRememberMe *int `json:"sso_session_max_lifespan_remember_me,omitempty" gorm:"type:int"`
//
// 	// SupportedLocales 支持的语言环境列表
// 	// realm 支持的语言
// 	// 示例: ["en", "zh-CN", "ja", "de"]
// 	SupportedLocales datatypes.JSONType[[]string] `json:"supported_locales,omitempty" gorm:"type:json"`
//
// 	// UserFederationMappers 用户联邦映射器
// 	// 定义如何映射外部用户存储的属性
// 	UserFederationMappers datatypes.JSON `json:"user_federation_mappers,omitempty" gorm:"type:json"`
//
// 	// UserFederationProviders 用户联邦提供者
// 	// 配置的外部用户存储（LDAP、AD等）
// 	UserFederationProviders datatypes.JSON `json:"user_federation_providers,omitempty" gorm:"type:json"`
//
// 	// UserManagedAccessAllowed 是否允许用户管理访问
// 	// true: 用户可以管理自己资源的访问权限
// 	// false: 只能由管理员管理
// 	UserManagedAccessAllowed *bool `json:"user_managed_access_allowed,omitempty" gorm:"type:boolean;default:false"`
//
// 	// Users 用户列表
// 	// realm 中的所有用户
// 	Users datatypes.JSONType[[]User] `json:"users,omitempty" gorm:"type:json"`
//
// 	// VerifyEmail 是否验证邮箱
// 	// true: 新用户必须验证邮箱
// 	// false: 不需要邮箱验证
// 	VerifyEmail *bool `json:"verify_email,omitempty" gorm:"type:boolean;default:false"`
//
// 	// WaitIncrementSeconds 等待增量（秒）
// 	// 暴力破解保护中，每次失败后增加的等待时间
// 	// 默认值通常为 60 秒
// 	WaitIncrementSeconds *int `json:"wait_increment_seconds,omitempty" gorm:"type:int"`
//
// 	// WebAuthnPolicyAcceptableAaguids WebAuthn 策略：可接受的 AAGUID 列表
// 	// 允许的认证器证明 GUID
// 	// 空列表表示接受所有认证器
// 	WebAuthnPolicyAcceptableAaguids datatypes.JSONType[[]string] `json:"web_authn_policy_acceptable_aaguids,omitempty" gorm:"type:json"`
//
// 	// WebAuthnPolicyAttestationConveyancePreference WebAuthn 策略：证明传输偏好
// 	// 定义如何处理认证器的证明
// 	// 值: "none", "indirect", "direct", "enterprise"
// 	WebAuthnPolicyAttestationConveyancePreference *string `json:"web_authn_policy_attestation_conveyance_preference,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyAuthenticatorAttachment WebAuthn 策略：认证器附件
// 	// 限制认证器类型
// 	// 值: "platform"（内置）, "cross-platform"（外部）, 或不限制
// 	WebAuthnPolicyAuthenticatorAttachment *string `json:"web_authn_policy_authenticator_attachment,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyAvoidSameAuthenticatorRegister WebAuthn 策略：避免重复注册同一认证器
// 	// true: 防止同一认证器多次注册
// 	// false: 允许重复注册
// 	WebAuthnPolicyAvoidSameAuthenticatorRegister *bool `json:"web_authn_policy_avoid_same_authenticator_register,omitempty" gorm:"type:boolean;default:false"`
//
// 	// WebAuthnPolicyCreateTimeout WebAuthn 策略：创建超时（秒）
// 	// 用户完成 WebAuthn 注册的超时时间
// 	// 默认值通常为 0（使用浏览器默认值）
// 	WebAuthnPolicyCreateTimeout *int `json:"web_authn_policy_create_timeout,omitempty" gorm:"type:int"`
//
// 	// WebAuthnPolicyPasswordlessAcceptableAaguids 无密码 WebAuthn 策略：可接受的 AAGUID 列表
// 	// 无密码登录允许的认证器证明 GUID
// 	WebAuthnPolicyPasswordlessAcceptableAaguids datatypes.JSONType[[]string] `json:"web_authn_policy_passwordless_acceptable_aaguids,omitempty" gorm:"type:json"`
//
// 	// WebAuthnPolicyPasswordlessAttestationConveyancePreference 无密码 WebAuthn 策略：证明传输偏好
// 	// 无密码登录的证明处理方式
// 	WebAuthnPolicyPasswordlessAttestationConveyancePreference *string `json:"web_authn_policy_passwordless_attestation_conveyance_preference,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyPasswordlessAuthenticatorAttachment 无密码 WebAuthn 策略：认证器附件
// 	// 无密码登录的认证器类型限制
// 	WebAuthnPolicyPasswordlessAuthenticatorAttachment *string `json:"web_authn_policy_passwordless_authenticator_attachment,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyPasswordlessAvoidSameAuthenticatorRegister 无密码 WebAuthn 策略：避免重复注册
// 	// 无密码登录是否防止同一认证器多次注册
// 	WebAuthnPolicyPasswordlessAvoidSameAuthenticatorRegister *bool `json:"web_authn_policy_passwordless_avoid_same_authenticator_register,omitempty" gorm:"type:boolean;default:false"`
//
// 	// WebAuthnPolicyPasswordlessCreateTimeout 无密码 WebAuthn 策略：创建超时（秒）
// 	// 无密码登录注册的超时时间
// 	WebAuthnPolicyPasswordlessCreateTimeout *int `json:"web_authn_policy_passwordless_create_timeout,omitempty" gorm:"type:int"`
//
// 	// WebAuthnPolicyPasswordlessRequireResidentKey 无密码 WebAuthn 策略：要求常驻密钥
// 	// 是否要求认证器支持常驻密钥
// 	// 值: "Yes", "No", "not specified"
// 	WebAuthnPolicyPasswordlessRequireResidentKey *string `json:"web_authn_policy_passwordless_require_resident_key,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyPasswordlessRpEntityName 无密码 WebAuthn 策略：RP 实体名称
// 	// 无密码登录时显示的依赖方名称
// 	WebAuthnPolicyPasswordlessRpEntityName *string `json:"web_authn_policy_passwordless_rp_entity_name,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyPasswordlessRpID 无密码 WebAuthn 策略：RP ID
// 	// 无密码登录的依赖方标识符
// 	// 通常是域名
// 	WebAuthnPolicyPasswordlessRpID *string `json:"web_authn_policy_passwordless_rp_id,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyPasswordlessSignatureAlgorithms 无密码 WebAuthn 策略：签名算法列表
// 	// 无密码登录接受的签名算法
// 	// 示例: ["ES256", "RS256"]
// 	WebAuthnPolicyPasswordlessSignatureAlgorithms datatypes.JSONType[[]string] `json:"web_authn_policy_passwordless_signature_algorithms,omitempty" gorm:"type:json"`
//
// 	// WebAuthnPolicyPasswordlessUserVerificationRequirement 无密码 WebAuthn 策略：用户验证要求
// 	// 无密码登录的用户验证级别
// 	// 值: "required", "preferred", "discouraged"
// 	WebAuthnPolicyPasswordlessUserVerificationRequirement *string `json:"web_authn_policy_passwordless_user_verification_requirement,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyRequireResidentKey WebAuthn 策略：要求常驻密钥
// 	// 是否要求认证器支持常驻密钥（可发现凭证）
// 	// 值: "Yes", "No", "not specified"
// 	WebAuthnPolicyRequireResidentKey *string `json:"web_authn_policy_require_resident_key,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyRpEntityName WebAuthn 策略：RP 实体名称
// 	// 在 WebAuthn 提示中显示的依赖方名称
// 	// 示例: "我的应用程序"
// 	WebAuthnPolicyRpEntityName *string `json:"web_authn_policy_rp_entity_name,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicyRpID WebAuthn 策略：RP ID
// 	// 依赖方标识符，通常是域名
// 	// 示例: "example.com"
// 	WebAuthnPolicyRpID *string `json:"web_authn_policy_rp_id,omitempty" gorm:"type:varchar(255)"`
//
// 	// WebAuthnPolicySignatureAlgorithms WebAuthn 策略：签名算法列表
// 	// 接受的签名算法
// 	// 示例: ["ES256", "ES384", "ES512", "RS256", "RS384", "RS512"]
// 	WebAuthnPolicySignatureAlgorithms datatypes.JSONType[[]string] `json:"web_authn_policy_signature_algorithms,omitempty" gorm:"type:json"`
//
// 	// WebAuthnPolicyUserVerificationRequirement WebAuthn 策略：用户验证要求
// 	// 定义用户验证的严格程度
// 	// 值: "required"（必须）, "preferred"（首选）, "discouraged"（不鼓励）
// 	WebAuthnPolicyUserVerificationRequirement *string `json:"web_authn_policy_user_verification_requirement,omitempty" gorm:"type:varchar(255)"`
//
// 	model.Base
// }
//
// // Role is a role
// type Role struct {
// 	ID                 *string                                       `json:"id,omitempty"`
// 	Name               *string                                       `json:"name,omitempty"`
// 	ScopeParamRequired *bool                                         `json:"scope_param_required,omitempty"`
// 	Composite          *bool                                         `json:"composite,omitempty"`
// 	Composites         datatypes.JSONType[*CompositesRepresentation] `json:"composites,omitempty"`
// 	ClientRole         *bool                                         `json:"client_role,omitempty"`
// 	ContainerID        *string                                       `json:"container_id,omitempty"`
// 	Description        *string                                       `json:"description,omitempty"`
// 	Attributes         *map[string][]string                          `json:"attributes,omitempty"`
// }
//
// // CompositesRepresentation represents the composite roles of a role
// type CompositesRepresentation struct {
// 	Client datatypes.JSONType[map[string][]string] `json:"client,omitempty"`
// 	Realm  *[]string                               `json:"realm,omitempty"`
// }
//
// // RolesRepresentation represents the roles of a realm
// type RolesRepresentation struct {
// 	Client datatypes.JSONType[map[string][]Role] `json:"client,omitempty"`
// 	Realm  *[]Role                               `json:"realm,omitempty"`
// }

func (*Realm) TableName() string {
	return "keycloak_realms"
}
