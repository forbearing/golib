package zap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	casbinl "github.com/casbin/casbin/v2/log"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gorml "gorm.io/gorm/logger"
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
	readConf()
	zap.ReplaceGlobals(zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.FatalLevel),
	))

	logger.Runtime = New("runtime.log")
	logger.Cronjob = New("cronjob.log")
	logger.Task = New("task.log")

	logger.Controller = New("controller.log")
	logger.Service = New("service.log")
	logger.Database = New("database.log")
	logger.Cache = New("cache.log")
	logger.Redis = New("redis.log")

	logger.Authz = New("authz.log")
	logger.Cassandra = New("cassandra.log")
	logger.Elastic = New("elastic.log")
	logger.Etcd = New("etcd.log")
	logger.Feishu = New("feishu.log")
	logger.Influxdb = New("influxdb.log")
	logger.Kafka = New("kafka.log")
	logger.Ldap = New("ldap.log")
	logger.Memcached = New("memcached.log")
	logger.Minio = New("minio.log")
	logger.Mongo = New("mongo.log")
	logger.Mqtt = New("mqtt.log")
	logger.Nats = New("nats.log")
	logger.Scylla = New("scylla.log")
	logger.RethinkDB = New("rethinkdb.log")
	logger.RocketMQ = New("rocketmq.log")

	logger.Protocol = New("protocol.log")
	logger.Binary = New("binary.log")

	logger.Gin = NewGin("access.log")
	logger.Gorm = NewGorm("gorm.log")
	logger.Casbin = NewCasbin("casbin.log")

	return nil
}

func Clean() {
	// types.Logger
	zap.L().Sync()
	logs := []types.Logger{
		logger.Runtime,
		logger.Cronjob,
		logger.Task,

		logger.Controller,
		logger.Service,
		logger.Database,
		logger.Cache,
		logger.Redis,

		logger.Authz,
		logger.Cassandra,
		logger.Elastic,
		logger.Etcd,
		logger.Feishu,
		logger.Influxdb,
		logger.Kafka,
		logger.Ldap,
		logger.Memcached,
		logger.Minio,
		logger.Mongo,
		logger.Mqtt,
		logger.Nats,
		logger.Scylla,
		logger.RethinkDB,
		logger.RocketMQ,

		logger.Protocol,
		logger.Binary,
	}
	for _, log := range logs {
		if l, ok := log.(*Logger); ok {
			l.zlog.Sync()
		}
	}

	// Gin logger
	logger.Gin.Sync()

	// gorm logger
	gormLogs := []gorml.Interface{
		logger.Gorm,
	}
	for _, glog := range gormLogs {
		if log, ok := glog.(*GormLogger); ok {
			if l, ok := log.l.(*Logger); ok {
				l.zlog.Sync()
			}
		}
	}

	// casbin logger
	casbinLogs := []casbinl.Logger{
		logger.Casbin,
	}
	for _, clog := range casbinLogs {
		if log, ok := clog.(*CasbinLogger); ok {
			if l, ok := log.l.(*Logger); ok {
				l.zlog.Sync()
			}
		}
	}
}

// New returns *Logger instance that contains *zap.Logger and *zap.SugaredLogger
// and implements types.Logger.
func New(filename ...string) *Logger {
	readConf()
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
	return &Logger{zlog: logger}
}

// NewGorm returns a *GormLogger instance that implements gorm logger.Interface.
func NewGorm(filename ...string) gorml.Interface {
	readConf()
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
	return &GormLogger{l: &Logger{zlog: logger}}
}

// NewCasbin returns a *GormLogger instance that implements casbin Logger interface.
// This logger without 'caller' field.
func NewCasbin(filename ...string) casbinl.Logger {
	readConf()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	logger := zap.New(
		zapcore.NewCore(newLogEncoder(), newLogWriter(), newLogLevel()),
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	return &CasbinLogger{l: &Logger{zlog: logger}}
}

// NewGin returns a *Logger instance that contains *zap.Logger.
func NewGin(filename ...string) *zap.Logger {
	readConf()
	if len(filename) > 0 {
		if len(filename[0]) > 0 {
			logFile = filename[0]
		}
	}
	return zap.New(zapcore.NewCore(newLogEncoder(option{disableMsg: true, disableLevel: true}), newLogWriter(), newLogLevel()))
}

// NewStdLog returns a *log.Logger instance that contains *zap.Logger.
func NewStdLog() *log.Logger {
	return zap.NewStdLog(NewZap())
}

// NewZap new a *zap.Logger instance according to Server/Client configurations.
// The instance implements types.Logger interface.
// log file default to config.Server.LoggerConfig.LogFile or config.Client.LoggerConfig.LogFile,
// you can create a custom *zap.Logger by pass log filename to this function.
func NewZap(filename ...string) *zap.Logger {
	readConf()
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
	readConf()
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
			Filename:   filepath.Join(config.App.Dir, logFile),
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
	encConfig.EncodeTime = zapcore.TimeEncoderOfLayout(consts.LayoutTimeEncoder)
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
		// return newCustomConsoleEncoder(encConfig)
		return zapcore.NewConsoleEncoder(encConfig)
	default:
		return zapcore.NewJSONEncoder(encConfig)
	}
}

func readConf() {
	mode = config.App.Mode
	logFile = config.App.Logger.File
	logLevel = config.App.Logger.Level
	logFormat = config.App.Logger.Format
	logEncoder = config.App.Logger.Encoder
	logMaxAge = int(config.App.Logger.MaxAge)
	logMaxSize = int(config.App.Logger.MaxSize)
	logMaxBackups = int(config.App.Logger.MaxBackups)
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

func newCustomConsoleEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	return &customConsoleEncoder{zapcore.NewConsoleEncoder(config)}
}

type customConsoleEncoder struct {
	zapcore.Encoder
}

func (e *customConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line, err := e.Encoder.EncodeEntry(ent, nil)
	if err != nil {
		return nil, err
	}

	// if contains extra fields, append them in key=value format
	if len(fields) > 0 {
		line.TrimNewline() // remove trailing newline
		// add extra fields
		for i, f := range fields {
			if i > 0 {
				line.AppendString("\t")
			} else {
				line.AppendString("\t")
			}
			line.AppendString(f.Key)
			line.AppendString("=")
			// according to the field type, format the value
			switch f.Type {
			case zapcore.StringType:
				line.AppendString(f.String)
			case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
				line.AppendString(strconv.FormatInt(f.Integer, 10))
			// you can add more types here
			default:
				line.AppendString(fmt.Sprint(f.Interface))
			}
		}
		line.AppendString("\n")
	}

	return line, nil
}
