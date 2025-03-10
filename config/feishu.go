package config

const (
	FEISHU_APP_ID     = "FEISHU_APP_ID"
	FEISHU_APP_SECRET = "FEISHU_APP_SECRET"
	FEISHU_ENABLE     = "FEISHU_ENABLE"
)

type Feishu struct {
	AppID     string `json:"app_id" mapstructure:"app_id" ini:"app_id" yaml:"app_id"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" ini:"app_secret" yaml:"app_secret"`
	Enable    bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Feishu) setDefault() {
	cv.SetDefault("feishu.app_id", "")
	cv.SetDefault("feishu.app_secret", "")
	cv.SetDefault("feishu.enable", false)
}
