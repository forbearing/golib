package zap

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	mode          config.Mode
	logFile       string
	logLevel      string
	logFormat     string
	logEncoder    string
	logMaxAge     int
	logMaxSize    int
	logMaxBackups int
)

// option contains logger constructor options.
type option struct {
	disableMsg   bool   // disable field "msg".
	disableLevel bool   // disable field "level".
	tsLayout     string // fields "ts" layout.
}

// Init will initial global *zap.Logger according to Server/Client configurations.
// log file default to config.Server.LoggerConfig.LogFile or config.Client.LoggerConfig.LogFile.
func Init() error {
	initVar()
	zap.ReplaceGlobals(zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.FatalLevel),
	))

	logger.Global = New()
	logger.Internal = New("internal.log")
	logger.Controller = New("controller.log")
	logger.Service = New("service.log")
	logger.Database = New("database.log")
	logger.Cache = New("cache.log")
	logger.Elastic = New("elastic.log")
	logger.Redis = New("redis.log")
	logger.Task = New("task.log")
	logger.Visitor = New("visitor.log")
	logger.Cronjob = New("cronjob.log")
	logger.Job = New("job.log")
	logger.Gin = NewGin("access.log")
	logger.Gorm = NewGorm("gorm.log")

	return nil
}

// New returns *Logger instance that contains *zap.Logger and *zap.SugaredLogger
// and implements types.Logger.
func New(filename ...string) *Logger {
	initVar()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	logger := zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddCallerSkip(1), // 这里别忘了
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	return &Logger{
		zlog: logger,
		slog: logger.Sugar(),
	}
}

// NewGorm returns a *GormLogger instance that implements gorm logger.Interface.
// The difference between NewGorm and NewLogger is the `zap.AddCallerSkip()`
func NewGorm(filename ...string) *GormLogger {
	initVar()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	logger := zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddCallerSkip(5), // 这个值是后期调出来的.
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	return &GormLogger{l: &Logger{zlog: logger, slog: logger.Sugar()}}
}

// NewGin returns a *Logger instance that contains *zap.Logger.
// The difference between NewGin and New is disable fields "caller", "level" and "msg".
func NewGin(filename ...string) *zap.Logger {
	initVar()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	return zap.New(zapcore.NewCore(newLogEncoder(option{disableMsg: true, disableLevel: true}), newLogWriter(), newLogLevel()))
}

// NewStdLog
func NewStdLog() *log.Logger {
	return zap.NewStdLog(NewZap())
}

// NewZap new a *zap.Logger instance according to Server/Client configurations.
// The instance implements types.Logger interface.
// log file default to config.Server.LoggerConfig.LogFile or config.Client.LoggerConfig.LogFile,
// you can create a custom *zap.Logger by pass log filename to this function.
func NewZap(filename ...string) *zap.Logger {
	initVar()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	return zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.FatalLevel))
}

// NewSugared new a *zap.SugaredLogger instance according to Server/Client configurations.
// The instance implements types.Logger and types.StructuredLogger interface.
// log file default to config.Server.LoggerConfig.LogFile or config.Client.LoggerConfig.LogFile,
// you can create a custom *zap.SugaredLogger by pass log filename to this function.
func NewSugared(filename ...string) *zap.SugaredLogger {
	initVar()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	return zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.FatalLevel)).Sugar()
}

// newLogWriter
func newLogWriter(_ ...option) zapcore.WriteSyncer {
	switch strings.TrimSpace(logFile) {
	case "/dev/stdout":
		return zapcore.AddSync(os.Stdout)
	case "/dev/stderr":
		return zapcore.AddSync(os.Stderr)
	case "":
		return zapcore.AddSync(os.Stdout)
	default:
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(config.App.LogDir, logFile),
			MaxAge:     logMaxAge,
			MaxSize:    logMaxSize,
			MaxBackups: logMaxBackups,
			LocalTime:  true,
			Compress:   false, // openwrt may not support compress.
		})
	}
}

// newLogLevel
func newLogLevel(_ ...option) zapcore.Level {
	if len(logLevel) == 0 {
		return zapcore.InfoLevel
	}
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return zapcore.InfoLevel
	}
	return *level
}

// newLogEncoder
func newLogEncoder(opt ...option) zapcore.Encoder {
	encConfig := zap.NewProductionEncoderConfig()
	// encConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	// encConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encConfig.EncodeDuration = zapcore.MillisDurationEncoder
	// encConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// encConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	// encConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02|15:04:05")
	encConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// encConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// encConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	// encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// encConfig.EncodeLevel = colorfulLevelEncoder
	if len(opt) > 0 {
		o := opt[0]
		if o.disableMsg {
			encConfig.MessageKey = ""
		}
		if o.disableLevel {
			encConfig.LevelKey = ""
		}
		if len(o.tsLayout) > 0 {
			encConfig.EncodeTime = zapcore.TimeEncoderOfLayout(o.tsLayout)
		}
	}
	switch strings.ToLower(logFormat) {
	case "json":
		return zapcore.NewJSONEncoder(encConfig)
	case "text", "console":
		return zapcore.NewConsoleEncoder(encConfig)
	default:
		return zapcore.NewJSONEncoder(encConfig)
	}
}

func initVar() {
	mode = config.App.Mode
	logFile = config.App.LogFile
	logLevel = config.App.LogLevel
	logFormat = config.App.LogFormat
	logEncoder = config.App.LogEncoder
	logMaxAge = int(config.App.LogMaxAge)
	logMaxSize = int(config.App.LogMaxSize)
	logMaxBackups = int(config.App.LogMaxBackups)
}

// colorfulLevelEncoder 自定义 Level Encoder，为不同的日志级别添加颜色
func colorfulLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[36m" // Cyan
	case zapcore.InfoLevel:
		color = "\033[32m" // Green
	case zapcore.WarnLevel:
		color = "\033[33m" // Yellow
	case zapcore.ErrorLevel:
		color = "\033[31m" // Red
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = "\033[35m" // Magenta
	default:
		color = "\033[0m" // Reset
	}
	// 使用颜色代码包装原始 Level 字符串
	enc.AppendString(color + level.String() + "\033[0m")
}
