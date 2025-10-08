package config

const (
	AUDIT_ENABLE                 = "AUDIT_ENABLE"
	AUDIT_OPERATIONS             = "AUDIT_OPERATIONS"
	AUDIT_TABLES                 = "AUDIT_TABLES"
	AUDIT_EXCLUDE_TABLES         = "AUDIT_EXCLUDE_TABLES"
	AUDIT_RECORD_REQUEST_BODY    = "AUDIT_RECORD_REQUEST_BODY"
	AUDIT_RECORD_RESPONSE_BODY   = "AUDIT_RECORD_RESPONSE_BODY"
	AUDIT_RECORD_OLD_VALUES      = "AUDIT_RECORD_OLD_VALUES"
	AUDIT_RECORD_NEW_VALUES      = "AUDIT_RECORD_NEW_VALUES"
	AUDIT_EXCLUDE_FIELDS         = "AUDIT_EXCLUDE_FIELDS"
	AUDIT_INCLUDE_FIELDS         = "AUDIT_INCLUDE_FIELDS"
	AUDIT_MAX_FIELD_LENGTH       = "AUDIT_MAX_FIELD_LENGTH"
	AUDIT_RECORD_QUERY_PARAMS    = "AUDIT_RECORD_QUERY_PARAMS"
	AUDIT_RECORD_USER_AGENT      = "AUDIT_RECORD_USER_AGENT"
	AUDIT_ASYNC_WRITE            = "AUDIT_ASYNC_WRITE"
	AUDIT_BATCH_SIZE             = "AUDIT_BATCH_SIZE"
	AUDIT_FLUSH_INTERVAL_SECONDS = "AUDIT_FLUSH_INTERVAL_SECONDS"
)

type Audit struct {
	Enable               bool     `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
	Operations           []string `json:"operations" mapstructure:"operations" ini:"operations" yaml:"operations"`
	Tables               []string `json:"tables" mapstructure:"tables" ini:"tables" yaml:"tables"`
	ExcludeTables        []string `json:"exclude_tables" mapstructure:"exclude_tables" ini:"exclude_tables" yaml:"exclude_tables"`
	RecordRequestBody    bool     `json:"record_request_body" mapstructure:"record_request_body" ini:"record_request_body" yaml:"record_request_body"`
	RecordResponseBody   bool     `json:"record_response_body" mapstructure:"record_response_body" ini:"record_response_body" yaml:"record_response_body"`
	RecordOldValues      bool     `json:"record_old_values" mapstructure:"record_old_values" ini:"record_old_values" yaml:"record_old_values"`
	RecordNewValues      bool     `json:"record_new_values" mapstructure:"record_new_values" ini:"record_new_values" yaml:"record_new_values"`
	ExcludeFields        []string `json:"exclude_fields" mapstructure:"exclude_fields" ini:"exclude_fields" yaml:"exclude_fields"`
	IncludeFields        []string `json:"include_fields" mapstructure:"include_fields" ini:"include_fields" yaml:"include_fields"`
	MaxFieldLength       int      `json:"max_field_length" mapstructure:"max_field_length" ini:"max_field_length" yaml:"max_field_length"`
	RecordQueryParams    bool     `json:"record_query_params" mapstructure:"record_query_params" ini:"record_query_params" yaml:"record_query_params"`
	RecordUserAgent      bool     `json:"record_user_agent" mapstructure:"record_user_agent" ini:"record_user_agent" yaml:"record_user_agent"`
	AsyncWrite           bool     `json:"async_write" mapstructure:"async_write" ini:"async_write" yaml:"async_write"`
	BatchSize            int      `json:"batch_size" mapstructure:"batch_size" ini:"batch_size" yaml:"batch_size"`
	FlushIntervalSeconds int      `json:"flush_interval_seconds" mapstructure:"flush_interval_seconds" ini:"flush_interval_seconds" yaml:"flush_interval_seconds"`
}

func (*Audit) setDefault() {
	cv.SetDefault("audit.enable", false)
	cv.SetDefault("audit.operations", []string{})
	cv.SetDefault("audit.tables", []string{})
	cv.SetDefault("audit.exclude_tables", []string{})
	cv.SetDefault("audit.record_request_body", true)
	cv.SetDefault("audit.record_response_body", true)
	cv.SetDefault("audit.record_old_values", true)
	cv.SetDefault("audit.record_new_values", true)
	cv.SetDefault("audit.exclude_fields", []string{"password", "passwd", "pwd", "secret", "token", "key", "private_key"})
	cv.SetDefault("audit.include_fields", []string{})
	cv.SetDefault("audit.max_field_length", 1000)
	cv.SetDefault("audit.record_query_params", true)
	cv.SetDefault("audit.record_user_agent", true)
	cv.SetDefault("audit.async_write", true)
	cv.SetDefault("audit.batch_size", 100)
	cv.SetDefault("audit.flush_interval_seconds", 30)
}
