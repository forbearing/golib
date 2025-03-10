package config

const (
	ELASTICSEARCH_ADDRS                   = "ELASTICSEARCH_ADDRS"
	ELASTICSEARCH_USERNAME                = "ELASTICSEARCH_USERNAME"
	ELASTICSEARCH_PASSWORD                = "ELASTICSEARCH_PASSWORD"
	ELASTICSEARCH_CLOUD_ID                = "ELASTICSEARCH_CLOUD_ID"
	ELASTICSEARCH_API_KEY                 = "ELASTICSEARCH_API_KEY"
	ELASTICSEARCH_SERVICE_TOKEN           = "ELASTICSEARCH_SERVICE_TOKEN"
	ELASTICSEARCH_CERTIFICATE_FINGERPRINT = "ELASTICSEARCH_CERTIFICATE_FINGERPRINT"
	ELASTICSEARCH_ENABLE                  = "ELASTICSEARCH_ENABLE"
)

type Elasticsearch struct {
	Addrs                  []string `json:"addrs" mapstructure:"addrs" ini:"addrs" yaml:"addrs"`
	Username               string   `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password               string   `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	CloudID                string   `json:"cloud_id" mapstructure:"cloud_id" ini:"cloud_id" yaml:"cloud_id"`
	APIKey                 string   `json:"api_key" mapstructure:"api_key" ini:"api_key" yaml:"api_key"`
	ServiceToken           string   `json:"service_token" mapstructure:"service_token" ini:"service_token" yaml:"service_token"`
	CertificateFingerprint string   `json:"certificate_fingerprint" mapstructure:"certificate_fingerprint" ini:"certificate_fingerprint" yaml:"certificate_fingerprint"`
	Enable                 bool     `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Elasticsearch) setDefault() {
	cv.SetDefault("elasticsearch.addrs", []string{"http://127.0.0.1:9200"})
	cv.SetDefault("elasticsearch.username", "")
	cv.SetDefault("elasticsearch.password", "")
	cv.SetDefault("elasticsearch.cloud_id", "")
	cv.SetDefault("elasticsearch.api_key", "")
	cv.SetDefault("elasticsearch.service_token", "")
	cv.SetDefault("elasticsearch.certificate_fingerprint", "")
	cv.SetDefault("elasticsearch.enable", false)
}
