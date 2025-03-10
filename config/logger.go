package config

const (
	LOGGER_DIR         = "LOGGER_DIR"
	LOGGER_PREFIX      = "LOGGER_PREFIX"
	LOGGER_FILE        = "LOGGER_FILE"
	LOGGER_LEVEL       = "LOGGER_LEVEL"
	LOGGER_FORMAT      = "LOGGER_FORMAT"
	LOGGER_ENCODER     = "LOGGER_ENCODER"
	LOGGER_MAX_AGE     = "LOGGER_MAX_AGE"
	LOGGER_MAX_SIZE    = "LOGGER_MAX_SIZE"
	LOGGER_MAX_BACKUPS = "LOGGER_MAX_BACKUPS"
)

// Logger represents section "logger" for client-side or server-side configuration,
// and there is only one copy during the application entire lifetime.
type Logger struct {
	// Dir specifies which direcotory log to.
	Dir string `json:"dir" ini:"dir" yaml:"dir" mapstructure:"dir"`

	// Prefix specifies the log prefix.
	// You can set the prefix name to your project name.
	Prefix string `json:"prefix" ini:"prefix" yaml:"prefix" mapstructure:"prefix"`

	// File specifies the which file log to.
	// If value is "/dev/stdout", log to os.Stdout.
	// If value is "/dev/stderr", log to os.Stderr.
	// If value is empty(length is zero), log to os.Stdout.
	File string `json:"file" ini:"file" yaml:"file" mapstructure:"file"`

	// Level specifies the log level,  supported values are: (error|warn|warning|info|debug).
	// The value default to "info" and ignore case.
	Level string `json:"level" ini:"level" yaml:"level" mapstructure:"level"`

	// Format specifies the log format, supported values are: (json|text).
	// The Value default to "text" and ignore case.
	Format string `json:"format" ini:"format" yaml:"format" mapstructure:"format"`

	// Encoder is the same as LogFormat.
	Encoder string `json:"encoder" ini:"encoder" yaml:"encoder" mapstructure:"encoder"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.
	// uint is "day" and default to 7.
	MaxAge uint `json:"max_age" ini:"max_age" yaml:"max_age" mapstructure:"max_age"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated, default to 1MB.
	MaxSize uint `json:"max_size" ini:"max_size" yaml:"max_size" mapstructure:"max_size"`

	// MaxBackups is the maximum number of old log files to retain.
	// The value default to 3.
	MaxBackups uint `json:"max_backups" ini:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`
}

func (*Logger) setDefault() {
	cv.SetDefault("logger.dir", "./logs")
	cv.SetDefault("logger.prefix", "")
	cv.SetDefault("logger.file", "")
	cv.SetDefault("logger.level", "info")
	cv.SetDefault("logger.format", "json")
	cv.SetDefault("logger.encoder", "json")
	cv.SetDefault("logger.max_age", 30)
	cv.SetDefault("logger.max_size", 100)
	cv.SetDefault("logger.max_backups", 1)
}
