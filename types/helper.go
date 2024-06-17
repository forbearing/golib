package types

// LoggerWithContext build *types.Logger from *types.ServiceContext.
func LoggerWithContext(ctx *ServiceContext, l Logger, phase Phase) Logger {
	return l.With(PHASE, string(phase)).
		With(CTX_USERNAME, ctx.Username).
		With(CTX_USER_ID, ctx.UserId).
		With(REQUEST_ID, ctx.RequestId)
}
