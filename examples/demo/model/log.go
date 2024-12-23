package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Log]()
}

type LogLevel string

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type Log struct {
	pkgmodel.Base

	Level    LogLevel `json:"level"`
	Message  string   `json:"message"`
	SourceID string   `json:"source_id"`
}
