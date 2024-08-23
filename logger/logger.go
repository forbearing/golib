// Package logger provides global logger used by server, client and cli.
package logger

import (
	"github.com/forbearing/golib/types"
	"go.uber.org/zap"
	gorml "gorm.io/gorm/logger"
)

var (
	Global     types.Logger
	Internal   types.Logger
	Controller types.Logger
	Service    types.Logger
	Database   types.Logger
	Cache      types.Logger
	Redis      types.Logger
	Task       types.Logger
	Visitor    types.Logger
	Cronjob    types.Logger
	Job        types.Logger
)
var Gin *zap.Logger
var (
	Gorm         gorml.Interface
	GormDLP      gorml.Interface
	GormSOC      gorml.Interface
	GormSocAgent gorml.Interface
	GormSoftware gorml.Interface
	GormCulture  gorml.Interface
)

func Debug(args ...any)                       { Global.Debug(args...) }
func Info(args ...any)                        { Global.Info(args...) }
func Warn(args ...any)                        { Global.Warn(args...) }
func Error(args ...any)                       { Global.Error(args...) }
func Fatal(args ...any)                       { Global.Fatal(args...) }
func Debugf(format string, args ...any)       { Global.Debugf(format, args...) }
func Infof(format string, args ...any)        { Global.Infof(format, args...) }
func Warnf(format string, args ...any)        { Global.Warnf(format, args...) }
func Errorf(format string, args ...any)       { Global.Errorf(format, args...) }
func Fatalf(format string, args ...any)       { Global.Fatalf(format, args...) }
func Debugw(msg string, keysAndValues ...any) { Global.Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...any)  { Global.Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...any)  { Global.Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...any) { Global.Errorw(msg, keysAndValues...) }
func Fatalw(msg string, keysAndValues ...any) { Global.Fatalw(msg, keysAndValues...) }
func Debugz(msg string, fields ...zap.Field)  { Global.Debugz(msg, fields...) }
func Infoz(msg string, fields ...zap.Field)   { Global.Infoz(msg, fields...) }
func Warnz(msg string, feilds ...zap.Field)   { Global.Warnz(msg, feilds...) }
func Errorz(msg string, fields ...zap.Field)  { Global.Errorz(msg, fields...) }
func Fatalz(msg string, fields ...zap.Field)  { Global.Fatalz(msg, fields...) }
