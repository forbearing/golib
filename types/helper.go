package types

import "github.com/forbearing/golib/types/consts"

// LoggerWithContext build *types.Logger from *types.ServiceContext.
func LoggerWithContext(ctx *ServiceContext, l Logger, phase consts.Phase) Logger {
	return l.With(consts.PHASE, string(phase)).
		With(consts.CTX_USERNAME, ctx.Username).
		With(consts.CTX_USER_ID, ctx.UserId).
		With(consts.REQUEST_ID, ctx.RequestId)
}
