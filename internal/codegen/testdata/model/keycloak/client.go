package model_keycloak

import (
	"github.com/forbearing/golib/model"
	"gorm.io/datatypes"
)

type Client struct {
	AdminURL                     *string `json:"admin_url,omitempty"`
	AuthorizationServicesEnabled *bool   `json:"authorization_services_enabled,omitempty"`
	BaseURL                      *string `json:"base_url,omitempty"`
	BearerOnly                   *bool   `json:"bearer_only,omitempty"`
	ClientAuthenticatorType      *string `json:"client_authenticator_type,omitempty"`
	ClientID                     *string `json:"client_id,omitempty"`
	ConsentRequired              *bool   `json:"consent_required,omitempty"`
	Description                  *string `json:"description,omitempty"`
	DirectAccessGrantsEnabled    *bool   `json:"direct_access_grants_enabled,omitempty"`
	Enabled                      *bool   `json:"enabled,omitempty"`
	ImplicitFlowEnabled          *bool   `json:"implicit_flow_enabled,omitempty"`
	Name                         *string `json:"name,omitempty"`
	NotBefore                    *int32  `json:"not_before,omitempty"`
	Origin                       *string `json:"origin,omitempty"`
	Protocol                     *string `json:"protocol,omitempty"`
	PublicClient                 *bool   `json:"public_client,omitempty"`
	RegistrationAccessToken      *string `json:"registration_access_token,omitempty"`
	RootURL                      *string `json:"root_url,omitempty"`
	Secret                       *string `json:"secret,omitempty"`
	ServiceAccountsEnabled       *bool   `json:"service_accounts_enabled,omitempty"`
	StandardFlowEnabled          *bool   `json:"standard_flow_enabled,omitempty"`
	SurrogateAuthRequired        *bool   `json:"surrogate_auth_required,omitempty"`

	RealmId   *string `json:"realm_id,omitempty"`
	RealmName *string `json:"realm_name,omitempty"`

	// 完整的 Client 配置存储在 JSON 字段中
	Data datatypes.JSON `json:"data,omitempty"`

	model.Base
}

func (*Client) TableName() string {
	return "keycloak_clients"
}
