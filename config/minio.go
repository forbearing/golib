package config

const (
	MINIO_ENDPOINT   = "MINIO_ENDPOINT"
	MINIO_REGION     = "MINIO_REGION"
	MINIO_ACCESS_KEY = "MINIO_ACCESS_KEY"
	MINIO_SECRET_KEY = "MINIO_SECRET_KEY"
	MINIO_BUCKET     = "MINIO_BUCKET"
	MINIO_USE_SSL    = "MINIO_USE_SSL"
	MINIO_ENABLE     = "MINIO_ENABLE"
)

type Minio struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint" ini:"endpoint" yaml:"endpoint"`
	Region    string `json:"region" mapstructure:"region" ini:"region" yaml:"region"`
	AccessKey string `json:"access_key" mapstructure:"access_key" ini:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key" ini:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	UseSsl    bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	Enable    bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Minio) setDefault() {
	cv.SetDefault("minio.endpoint", "127.0.0.1:9000")
	cv.SetDefault("minio.region", "")
	cv.SetDefault("minio.access_key", "")
	cv.SetDefault("minio.secret_key", "")
	cv.SetDefault("minio.bucket", "")
	cv.SetDefault("minio.use_ssl", false)
	cv.SetDefault("minio.enable", false)
}
