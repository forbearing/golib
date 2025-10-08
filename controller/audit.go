package controller

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/database"
	"github.com/forbearing/gst/ds/queue/circularbuffer"
	model_log "github.com/forbearing/gst/model/log"
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/types/consts"
	"github.com/gin-gonic/gin"
)

// AuditParams contains parameters for audit logging.
// This struct encapsulates all necessary information for recording operation logs
// in a centralized and consistent manner across all Factory functions.
type AuditParams struct {
	OP          consts.OP         // Operation type (create, update, delete, etc.)
	Table       string            // Table/model name
	Model       string            // Model name
	RecordID    string            // Record ID
	RecordName  string            // Record name/title
	Record      any               // Full record content
	Request     any               // Request data
	Response    any               // Response data
	OldRecord   any               // Old record content (for updates)
	NewRecord   any               // New record content (for updates)
	QueryParams map[string]string // Query parameters
}

// AuditManager manages audit logging based on configuration.
// It provides a centralized way to handle operation logging across all Factory functions,
// replacing the previous direct enqueuing of OperationLog records.
// The manager supports configurable filtering, field exclusion, and data truncation.
type AuditManager struct {
	config *config.Audit
	cb     *circularbuffer.CircularBuffer[*model_log.OperationLog]
}

// NewAuditManager creates a new audit manager instance.
// This replaces the previous direct usage of circular buffer for operation logging.
func NewAuditManager(auditConfig *config.Audit, cb *circularbuffer.CircularBuffer[*model_log.OperationLog]) *AuditManager {
	return &AuditManager{
		config: auditConfig,
		cb:     cb,
	}
}

// RecordOperation records a single operation audit log.
// This method is now used by all Factory functions instead of directly enqueuing OperationLog records.
// It provides centralized audit logging with configurable filtering and supports both sync and async writing.
func (am *AuditManager) RecordOperation(c *gin.Context, params *AuditParams) error {
	if !am.IsEnabled(params.OP, params.Table) {
		return nil
	}

	operationLog := am.buildOperationLog(c, params)

	if am.config.AsyncWrite {
		// Use existing circular buffer for async writing
		am.cb.Enqueue(operationLog)
		return nil
	}

	// Synchronous writing
	if err := database.Database[*model_log.OperationLog](types.NewDatabaseContext(c)).Create(operationLog); err != nil {
		return errors.Wrap(err, "failed to write audit log")
	}
	return nil
}

// RecordBatchOperations records multiple operations audit logs
func (am *AuditManager) RecordBatchOperations(c *gin.Context, batchParams []AuditParams) error {
	if !am.config.Enable {
		return nil
	}

	var operationLogs []*model_log.OperationLog

	for _, params := range batchParams {
		if am.IsEnabled(params.OP, params.Table) {
			operationLog := am.buildOperationLog(c, &params)
			operationLogs = append(operationLogs, operationLog)
		}
	}

	if len(operationLogs) == 0 {
		return nil
	}

	if am.config.AsyncWrite {
		// Enqueue all logs to circular buffer
		for _, log := range operationLogs {
			am.cb.Enqueue(log)
		}
		return nil
	}

	// Synchronous batch writing
	if err := database.Database[*model_log.OperationLog](types.NewDatabaseContext(c)).Create(operationLogs...); err != nil {
		return errors.Wrap(err, "failed to write batch audit logs")
	}
	return nil
}

// IsEnabled checks if audit logging is enabled for the given operation and table
func (am *AuditManager) IsEnabled(op consts.OP, table string) bool {
	if !am.config.Enable {
		return false
	}

	return am.isOperationEnabled(op) && am.isTableEnabled(table)
}

// isOperationEnabled checks if a specific operation should be audited
func (am *AuditManager) isOperationEnabled(op consts.OP) bool {
	if !am.config.Enable {
		return false
	}

	// If no operations specified, audit all operations
	if len(am.config.Operations) == 0 {
		return true
	}

	return slices.Contains(am.config.Operations, strings.ToLower(string(op)))
}

// isTableEnabled checks if a specific table should be audited
func (am *AuditManager) isTableEnabled(table string) bool {
	if !am.config.Enable {
		return false
	}

	// Check exclude list first
	if len(am.config.ExcludeTables) > 0 && slices.Contains(am.config.ExcludeTables, table) {
		return false
	}

	// If no tables specified, audit all tables (except excluded ones)
	if len(am.config.Tables) == 0 {
		return true
	}

	return slices.Contains(am.config.Tables, table)
}

// shouldRecordField checks if a specific field should be recorded in audit logs
func (am *AuditManager) shouldRecordField(fieldName string) bool {
	fieldNameLower := strings.ToLower(fieldName)

	// If include fields are specified, only record those fields
	if len(am.config.IncludeFields) > 0 {
		for _, includeField := range am.config.IncludeFields {
			if strings.ToLower(includeField) == fieldNameLower {
				return true
			}
		}
		return false
	}

	// Check if field is in exclude list
	excludeFields := am.getExcludeFields()
	for _, excludeField := range excludeFields {
		if strings.ToLower(excludeField) == fieldNameLower {
			return false
		}
	}

	return true
}

// getExcludeFields returns the list of fields to exclude from audit logs
func (am *AuditManager) getExcludeFields() []string {
	// Default sensitive fields to exclude
	defaultExcludeFields := []string{
		"password", "passwd", "pwd", "secret", "token", "key", "private_key",
		"access_token", "refresh_token", "api_key", "auth_token", "session_id",
		"credit_card", "ssn", "social_security_number", "phone", "email",
	}

	// Merge with user-defined exclude fields
	excludeFields := make([]string, 0, len(defaultExcludeFields)+len(am.config.ExcludeFields))
	excludeFields = append(excludeFields, defaultExcludeFields...)
	excludeFields = append(excludeFields, am.config.ExcludeFields...)

	return excludeFields
}

// buildOperationLog builds an OperationLog from the given parameters
func (am *AuditManager) buildOperationLog(c *gin.Context, params *AuditParams) *model_log.OperationLog {
	operationLog := &model_log.OperationLog{
		OP:        string(params.OP),
		Table:     params.Table,
		Model:     params.Model,
		RecordId:  params.RecordID,
		Method:    c.Request.Method,
		URI:       c.Request.RequestURI,
		IP:        c.ClientIP(),
		RequestId: c.GetString(consts.REQUEST_ID),
	}

	// Set user information
	if user, exists := c.Get(consts.CTX_USERNAME); exists {
		if userStr, ok := user.(string); ok {
			operationLog.User = userStr
		}
	}

	// Set record name
	if params.RecordName != "" {
		operationLog.RecordName = params.RecordName
	} else {
		operationLog.RecordName = am.extractRecordName(params.Record)
	}

	// Set request data
	if am.config.RecordRequestBody && params.Request != nil {
		if requestData := am.filterAndTruncateFields(params.Request); requestData != "" {
			operationLog.Request = requestData
		}
	}

	// Set response data
	if am.config.RecordResponseBody && params.Response != nil {
		if responseData := am.filterAndTruncateFields(params.Response); responseData != "" {
			operationLog.Response = responseData
		}
	}

	// Set record data
	if params.Record != nil {
		if recordData := am.filterAndTruncateFields(params.Record); recordData != "" {
			operationLog.Record = recordData
		}
	}

	// Set old record data
	if am.config.RecordOldValues && params.OldRecord != nil {
		if oldRecordData := am.filterAndTruncateFields(params.OldRecord); oldRecordData != "" {
			operationLog.OldRecord = oldRecordData
		}
	}

	// Set new record data
	if am.config.RecordNewValues && params.NewRecord != nil {
		if newRecordData := am.filterAndTruncateFields(params.NewRecord); newRecordData != "" {
			operationLog.NewRecord = newRecordData
		}
	}

	// Set user agent
	if am.config.RecordUserAgent {
		operationLog.UserAgent = c.GetHeader("User-Agent")
	}

	return operationLog
}

// filterAndTruncateFields filters sensitive fields and truncates field values
func (am *AuditManager) filterAndTruncateFields(data any) string {
	if data == nil {
		return ""
	}

	// Convert to map for field filtering
	dataMap, err := am.toMap(data)
	if err != nil {
		// If conversion fails, try to marshal directly
		if jsonData, err := json.Marshal(data); err == nil {
			return am.truncateString(string(jsonData))
		}
		return ""
	}

	// Filter fields
	filteredMap := am.filterFields(dataMap)

	// Marshal to JSON
	jsonData, err := json.Marshal(filteredMap)
	if err != nil {
		return ""
	}

	return am.truncateString(string(jsonData))
}

// toMap converts any data type to map[string]any
func (am *AuditManager) toMap(data any) (map[string]any, error) {
	// First try to marshal and unmarshal to get a clean map
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// filterFields filters out sensitive fields from the data map
func (am *AuditManager) filterFields(dataMap map[string]any) map[string]any {
	if len(dataMap) == 0 {
		return dataMap
	}

	filteredMap := make(map[string]any)

	for key, value := range dataMap {
		if am.shouldRecordField(key) {
			// Recursively filter nested maps
			if nestedMap, ok := value.(map[string]any); ok {
				filteredMap[key] = am.filterFields(nestedMap)
			} else {
				filteredMap[key] = value
			}
		}
	}

	return filteredMap
}

// truncateString truncates string to the configured maximum length
func (am *AuditManager) truncateString(s string) string {
	if am.config.MaxFieldLength <= 0 || len(s) <= am.config.MaxFieldLength {
		return s
	}

	return s[:am.config.MaxFieldLength] + "..."
}

// extractRecordName extracts a meaningful name from the record
func (am *AuditManager) extractRecordName(record any) string {
	if record == nil {
		return ""
	}

	// Common field names that might represent a record name
	nameFields := []string{"name", "title", "username", "email", "code", "key", "id"}

	dataMap, err := am.toMap(record)
	if err != nil {
		return ""
	}

	for _, field := range nameFields {
		if value, exists := dataMap[field]; exists {
			if strValue, ok := value.(string); ok && strValue != "" {
				return strValue
			}
		}
	}

	return ""
}
