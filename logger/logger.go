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
	Elastic    types.Logger
	Redis      types.Logger
	Task       types.Logger
	Visitor    types.Logger
	Cronjob    types.Logger
	Job        types.Logger

	Gin  *zap.Logger
	Gorm gorml.Interface
)
