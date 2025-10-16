package auditmanager

import (
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/database"
	"github.com/forbearing/gst/ds/queue/circularbuffer"
	modellog "github.com/forbearing/gst/model/log"
	"github.com/forbearing/gst/types"
	"github.com/gertd/go-pluralize"
	"go.uber.org/zap"
)

var pluralizeCli = pluralize.NewClient()

// AuditManager manages audit logging based on configuration.
// It provides a centralized way to handle operation logging across all Factory functions,
// replacing the previous direct enqueuing of OperationLog records.
// The manager supports configurable filtering, field exclusion, and data truncation.
type AuditManager struct {
	config *config.Audit
	cb     *circularbuffer.CircularBuffer[*modellog.OperationLog]
}

// New creates a new audit manager instance.
// This replaces the previous direct usage of circular buffer for operation logging.
func New(auditConfig *config.Audit, cb *circularbuffer.CircularBuffer[*modellog.OperationLog]) *AuditManager {
	return &AuditManager{
		config: auditConfig,
		cb:     cb,
	}
}

// RecordOperation records a single operation audit log.
// This method is now used by all Factory functions instead of directly enqueuing OperationLog records.
// It provides centralized audit logging with configurable filtering and supports both sync and async writing.
func (am *AuditManager) RecordOperation(ctx *types.DatabaseContext, m types.Model, operationLog *modellog.OperationLog) error {
	// Skip if audit is disabled
	if !am.config.Enable {
		return nil
	}

	// Skip if the operation is excluded.
	if slices.Contains(am.config.ExcludeOperations, operationLog.OP) {
		return nil
	}

	// Record the table name
	tableName := m.GetTableName()
	if len(tableName) == 0 {
		typ := reflect.TypeOf(m).Elem()
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
	}
	operationLog.Table = tableName

	if am.config.AsyncWrite {
		// Use existing circular buffer for async writing
		am.cb.Enqueue(operationLog)
		return nil
	}

	// Synchronous writing
	if err := database.Database[*modellog.OperationLog](ctx).Create(operationLog); err != nil {
		return errors.Wrap(err, "failed to write audit log")
	}
	return nil
}

// RecordBatchOperations records multiple operations audit logs
func (am *AuditManager) RecordBatchOperations(ctx *types.DatabaseContext, m types.Model, operationLogs []*modellog.OperationLog) error {
	if !am.config.Enable {
		return nil
	}

	if len(operationLogs) == 0 {
		return nil
	}

	// Record the table name
	tableName := m.GetTableName()
	if len(tableName) == 0 {
		typ := reflect.TypeOf(m).Elem()
		items := strings.Split(typ.Name(), ".")
		if len(items) > 0 {
			tableName = pluralizeCli.Plural(strings.ToLower(items[len(items)-1]))
		}
	}
	for _, operationLog := range operationLogs {
		operationLog.Table = tableName
	}

	if am.config.AsyncWrite {
		// Enqueue all logs to circular buffer
		for _, log := range operationLogs {
			am.cb.Enqueue(log)
		}
		return nil
	}

	// Synchronous batch writing
	if err := database.Database[*modellog.OperationLog](ctx).Create(operationLogs...); err != nil {
		return errors.Wrap(err, "failed to write batch audit logs")
	}
	return nil
}

// Consume operation log.
func (am *AuditManager) Consume() {
	operationLogs := make([]*modellog.OperationLog, 0, config.App.Server.CircularBuffer.SizeOperationLog)
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		operationLogs = operationLogs[:0]
		for !am.cb.IsEmpty() {
			ol, _ := am.cb.Dequeue()
			operationLogs = append(operationLogs, ol)
		}
		if len(operationLogs) > 0 {
			if err := database.Database[*modellog.OperationLog](nil).WithBatchSize(1000).Create(operationLogs...); err != nil {
				zap.S().Error(err)
			}
		}
	}
}
