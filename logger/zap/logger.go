package zap

import (
	"context"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"go.uber.org/zap"
	gorml "gorm.io/gorm/logger"
)

// Logger implements types.Logger interface.
type Logger struct {
	zlog *zap.Logger
}

var _ types.Logger = (*Logger)(nil)

func (l *Logger) Debug(args ...any) { l.zlog.Sugar().Debug(args...) }
func (l *Logger) Info(args ...any)  { l.zlog.Sugar().Info(args...) }
func (l *Logger) Warn(args ...any)  { l.zlog.Sugar().Warn(args...) }
func (l *Logger) Error(args ...any) { l.zlog.Sugar().Error(args...) }
func (l *Logger) Fatal(args ...any) { l.zlog.Sugar().Fatal(args...) }

func (l *Logger) Debugf(format string, args ...any) { l.zlog.Sugar().Debugf(format, args...) }
func (l *Logger) Infof(format string, args ...any)  { l.zlog.Sugar().Infof(format, args...) }
func (l *Logger) Warnf(format string, args ...any)  { l.zlog.Sugar().Warnf(format, args...) }
func (l *Logger) Errorf(format string, args ...any) { l.zlog.Sugar().Errorf(format, args...) }
func (l *Logger) Fatalf(format string, args ...any) { l.zlog.Sugar().Fatalf(format, args...) }

func (l *Logger) Debugw(msg string, keysValues ...any) { l.zlog.Sugar().Debugw(msg, keysValues...) }
func (l *Logger) Infow(msg string, keysValues ...any)  { l.zlog.Sugar().Infow(msg, keysValues...) }
func (l *Logger) Warnw(msg string, keysValues ...any)  { l.zlog.Sugar().Warnw(msg, keysValues...) }
func (l *Logger) Errorw(msg string, keysValues ...any) { l.zlog.Sugar().Errorw(msg, keysValues...) }
func (l *Logger) Fatalw(msg string, keysValues ...any) { l.zlog.Sugar().Fatalw(msg, keysValues...) }

func (l *Logger) Debugz(msg string, fields ...zap.Field) { l.zlog.Debug(msg, fields...) }
func (l *Logger) Infoz(msg string, fields ...zap.Field)  { l.zlog.Info(msg, fields...) }
func (l *Logger) Warnz(msg string, fields ...zap.Field)  { l.zlog.Warn(msg, fields...) }
func (l *Logger) Errorz(msg string, fields ...zap.Field) { l.zlog.Error(msg, fields...) }
func (l *Logger) Fatalz(msg string, fields ...zap.Field) { l.zlog.Fatal(msg, fields...) }

func (l *Logger) ZapLogger() *zap.Logger { return l.zlog }

// With creates a new logger with additional string key-value pairs.
// Each pair of arguments must be a key(string) followed by its value(string).
// If an odd number of arguments is provided, an empty string will be appended as the last value.
//
// Example 1 - Multiple With calls:
//
//	logger.With("phase", "update").
//	      With("user", "admin").
//	      With("request_id", "123")
//
// Example 2 - Single With call with multiple fields:
//
//	logger.With(
//	    "phase", "update",
//	    "user", "admin",
//	    "request_id", "123",
//	)
//
// Returns the original logger if no fields are provided or if only an empty key is provided.
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

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		if len(fields[i]) == 0 {
			continue
		}
		zapFields = append(zapFields, zap.String(fields[i], fields[i+1]))
	}
	return &Logger{zlog: l.zlog.With(zapFields...)}
}

// WithControllerContext creates a new logger with controller context fields.
// It extends the base logger with phase, username, user ID, and request ID from *types.ControllerContext.
//
// examples:
//
// log := logger.Controller.WithControllerContext(ctx, consts.PHASE_LIST)
func (l *Logger) WithControllerContext(ctx *types.ControllerContext, phase consts.Phase) types.Logger {
	return l.With(
		consts.PHASE, string(phase),
		consts.CTX_ROUTE, ctx.Route,
		consts.CTX_USERNAME, ctx.Username,
		consts.CTX_USER_ID, ctx.UserId,
		consts.TRACE_ID, ctx.TraceId,
	)
}

// WithServiceContext creates a new logger with service context fields.
// It extends the base logger with phase, username, user ID, and request ID from *types.ServiceContext.
//
// examples:
//
// log := logger.Service.WithServiceContext(ctx, consts.PHASE_LIST_BEFORE)
func (l *Logger) WithServiceContext(ctx *types.ServiceContext, phase consts.Phase) types.Logger {
	return l.With(
		consts.PHASE, string(phase),
		consts.CTX_ROUTE, ctx.Route,
		consts.CTX_USERNAME, ctx.Username,
		consts.CTX_USER_ID, ctx.UserId,
		consts.TRACE_ID, ctx.TraceId,
	)
}

// WithDatabaseContext creates a new logger with database context fields.
// It extends the base logger with phase, username, user ID, and request ID from *types.DatabaseContext.
//
// examples:
//
// log := logger.Database.WithDatabaseContext(ctx, consts.PHASE_LIST_BEFORE)
func (l *Logger) WithDatabaseContext(ctx *types.DatabaseContext, phase consts.Phase) (clone types.Logger) {
	return l.With(
		consts.PHASE, string(phase),
		consts.CTX_ROUTE, ctx.Route,
		consts.CTX_USERNAME, ctx.Username,
		consts.CTX_USER_ID, ctx.UserId,
		consts.TRACE_ID, ctx.TraceId,
	)
}

// GormLogger implements gorm logger.Interface
type GormLogger struct{ l types.Logger }

var _ gorml.Interface = (*GormLogger)(nil)

func (g *GormLogger) LogMode(gorml.LogLevel) gorml.Interface           { return g }
func (g *GormLogger) Info(_ context.Context, str string, args ...any)  { g.l.Infow(str, args) }
func (g *GormLogger) Warn(_ context.Context, str string, args ...any)  { g.l.Warnw(str, args) }
func (g *GormLogger) Error(_ context.Context, str string, args ...any) { g.l.Errorw(str, args) }
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	username, _ := ctx.Value(consts.CTX_USERNAME).(string)
	userId, _ := ctx.Value(consts.CTX_USER_ID).(string)
	traceId, _ := ctx.Value(consts.TRACE_ID).(string)
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		g.l.Errorz("", zap.String("sql", sql), zap.Int64("rows", rows), zap.String("elapsed", elapsed.String()), zap.Error(err))
	} else {
		if elapsed > config.App.Database.SlowQueryThreshold {
			g.l.Warnz("slow SQL detected",
				zap.String(consts.CTX_USERNAME, username),
				zap.String(consts.CTX_USER_ID, userId),
				zap.String(consts.TRACE_ID, traceId),
				zap.String("sql", sql),
				zap.String("elapsed", elapsed.String()),
				zap.String("threshold", config.App.Database.SlowQueryThreshold.String()),
				zap.Int64("rows", rows))
		} else {
			g.l.Infoz("sql executed",
				zap.String(consts.CTX_USERNAME, username),
				zap.String(consts.CTX_USER_ID, userId),
				zap.String(consts.TRACE_ID, traceId),
				zap.String("sql", sql),
				zap.String("elapsed", elapsed.String()),
				zap.Int64("rows", rows))
		}
	}
}
