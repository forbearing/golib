package zap

import (
	"context"
	"time"

	"github.com/forbearing/golib/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gorml "gorm.io/gorm/logger"
)

// Logger implements types.Logger interface.
// https://github.com/moul/zapgorm2 may be the alternative choice.
// eg: gorm.Open(mysql.Open(dsnAsset), &gorm.Config{Logger: zapgorm2.New(pkgzap.NewZap("./logs/gorm_asset.log"))})
type Logger struct {
	zlog *zap.Logger
	slog *zap.SugaredLogger
}

var _ types.Logger = (*Logger)(nil)

func (l *Logger) Debug(args ...any) { l.slog.Debug(args...) }
func (l *Logger) Info(args ...any)  { l.slog.Info(args...) }
func (l *Logger) Warn(args ...any)  { l.slog.Warn(args...) }
func (l *Logger) Error(args ...any) { l.slog.Error(args...) }
func (l *Logger) Fatal(args ...any) { l.slog.Fatal(args...) }

func (l *Logger) Debugf(format string, args ...any) { l.slog.Debugf(format, args...) }
func (l *Logger) Infof(format string, args ...any)  { l.slog.Infof(format, args...) }
func (l *Logger) Warnf(format string, args ...any)  { l.slog.Warnf(format, args...) }
func (l *Logger) Errorf(format string, args ...any) { l.slog.Errorf(format, args...) }
func (l *Logger) Fatalf(format string, args ...any) { l.slog.Fatalf(format, args...) }

func (l *Logger) Debugw(msg string, keysAndValues ...any) { l.slog.Debugw(msg, keysAndValues...) }
func (l *Logger) Infow(msg string, keysAndValues ...any)  { l.slog.Infow(msg, keysAndValues...) }
func (l *Logger) Warnw(msg string, keysAndValues ...any)  { l.slog.Warnw(msg, keysAndValues...) }
func (l *Logger) Errorw(msg string, keysAndValues ...any) { l.slog.Errorw(msg, keysAndValues...) }
func (l *Logger) Fatalw(msg string, keysAndValues ...any) { l.slog.Fatalw(msg, keysAndValues...) }

func (l *Logger) Debugz(msg string, fields ...zap.Field) { l.zlog.Debug(msg, fields...) }
func (l *Logger) Infoz(msg string, fields ...zap.Field)  { l.zlog.Info(msg, fields...) }
func (l *Logger) Warnz(msg string, fields ...zap.Field)  { l.zlog.Warn(msg, fields...) }
func (l *Logger) Errorz(msg string, fields ...zap.Field) { l.zlog.Error(msg, fields...) }
func (l *Logger) Fatalz(msg string, fields ...zap.Field) { l.zlog.Fatal(msg, fields...) }

// With one or multiple fields.
// Examples:
//
// log := logger.Controller.
//
//	With(types.PHASE, string(types.PHASE_UPDATE)).
//	With(types.CTX_USERNAME, c.GetString(types.CTX_USERNAME)).
//	With(types.CTX_USER_ID, c.GetString(types.CTX_USER_ID)).
//	With(types.REQUEST_ID, c.GetString(types.REQUEST_ID))
//
// log := logger.Controller.With(
//
//	types.PHASE, string(types.PHASE_DELETE),
//	types.CTX_USERNAME, c.GetString(types.CTX_USERNAME),
//	types.CTX_USER_ID, c.GetString(types.CTX_USER_ID),
//	types.REQUEST_ID, c.GetString(types.REQUEST_ID),
//	)
func (l *Logger) With(fields ...string) types.Logger {
	if len(fields) == 0 {
		return l
	}
	if len(fields) == 1 {
		if len(fields[0]) == 0 {
			return l
		}
	}
	if len(fields)%2 != 0 {
		fields = append(fields, "")
	}
	var zapFields []zapcore.Field
	var sugaredFields []any
	for i := 0; i < len(fields); i += 2 {
		key, val := fields[i], fields[i+1]
		zapFields = append(zapFields, zap.String(key, val))
		sugaredFields = append(sugaredFields, key, val)
	}
	clone := new(Logger)
	clone.zlog = l.zlog.With(zapFields...)
	clone.slog = l.slog.With(sugaredFields...)
	return clone
}

// GormLogger implements gorm logger.Interface
// https://github.com/moul/zapgorm2 may be the alternative choice.
// eg: gorm.Open(mysql.Open(dsnAsset), &gorm.Config{Logger: zapgorm2.New(pkgzap.NewZap("./logs/gorm_asset.log"))})
type GormLogger struct {
	l *Logger
}

var _ gorml.Interface = (*GormLogger)(nil)

func (g *GormLogger) LogMode(gorml.LogLevel) gorml.Interface { return g }
func (g *GormLogger) Info(_ context.Context, str string, args ...any) {
	g.l.Infow(str, args)
}

func (g *GormLogger) Warn(_ context.Context, str string, args ...any) {
	g.l.Warnw(str, args)
}

func (g *GormLogger) Error(_ context.Context, str string, args ...any) {
	g.l.Errorw(str, args)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		// g.l.Error(fmt.Sprintf("%s [%.3fms] %s rows:%v error:%v", "query", float64(elapsed.Nanoseconds())/1e6, sql, rows, err))
		g.l.Errorz(sql, zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()), zap.Error(err))
	} else {
		// g.l.Info(fmt.Sprintf("%s [%.3fms] %s rows:%v", "query", float64(elapsed.Nanoseconds())/1e6, sql, rows))
		g.l.Infow(sql, zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()))
	}
}
