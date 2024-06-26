// Package logger provides global logger used by server, client and cli.
package logger

import (
	"github.com/forbearing/golib/types"
	"go.uber.org/zap"
	gorml "gorm.io/gorm/logger"
)

// Global is the application global logger.
var Global types.Logger

// Internal is the application internal logger.
var Internal types.Logger

// Controller is a logger for controller layer.
var Controller types.Logger

// Service is a logger for service layer.
var Service types.Logger

// Database is a logger for database.
var Database types.Logger

// Cache is a logger for cache.
var Cache types.Logger

// Redis is a logger for redis.
var Redis types.Logger

// Task is a logger for task.
var Task types.Logger

// Visitor logging system and application runtime state.
var Visitor types.Logger

// Cronjob logging cronjob log.
var Cronjob types.Logger

// Job logging job log.
var Job types.Logger

// Gin is a logger for gin framework.
var Gin *zap.Logger

// Gorm is a logger for gorm.
var Gorm gorml.Interface

// GormDLP is a logger for gorm
var GormDLP gorml.Interface

// GormSOC is a logger for gorm
var GormSOC gorml.Interface

// GormSocAgent is a logger for gorm.
var GormSocAgent gorml.Interface

// GormSoftware is a logger for gorm
var GormSoftware gorml.Interface

// GormCulture is a logger for gorm
var GormCulture gorml.Interface
