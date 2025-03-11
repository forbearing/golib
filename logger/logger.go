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
	Elastic    types.Logger
	Mongo      types.Logger
	Mqtt       types.Logger
	Task       types.Logger
	Runtime    types.Logger
	Cronjob    types.Logger
	Job        types.Logger

	Protocol types.Logger
	Binary   types.Logger

	Gin  *zap.Logger
	Gorm gorml.Interface
)
