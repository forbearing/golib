package types

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
	ParamName string
}

// QueryConfig configures the behavior of WithQuery method.
//
// Fields:
//   - FuzzyMatch: Enable fuzzy matching (LIKE/REGEXP queries). Default: false (exact match with IN clause)
//   - AllowEmpty: Allow empty query conditions to match all records. Default: false (blocked for safety)
//
// CRITICAL SAFETY FEATURE:
// Empty query conditions (all fields are zero values) are blocked by default to prevent
// catastrophic data loss scenarios, especially when the result is used for Delete operations.
//
// Empty Query Examples:
//   - WithQuery(&User{})                    → all fields are zero values
//   - WithQuery(&User{Name: "", Email: ""}) → all field values are empty strings
//   - WithQuery(&KV{Key: ""})               → happens when removed slice is empty
//
// Usage Examples:
//
//	// Exact match (default)
//	WithQuery(&User{Name: "John"})
//	WithQuery(&User{Name: "John"}, QueryConfig{})
//
//	// Fuzzy match
//	WithQuery(&User{Name: "John"}, QueryConfig{FuzzyMatch: true})
//
//	// Allow empty query (ListFactory with pagination)
//	WithQuery(&User{}, QueryConfig{AllowEmpty: true})
//
//	// Fuzzy match + Allow empty
//	WithQuery(&User{}, QueryConfig{FuzzyMatch: true, AllowEmpty: true})
type QueryConfig struct {
	FuzzyMatch bool // Enable fuzzy matching (LIKE/REGEXP). Default: false
	AllowEmpty bool // Allow empty query conditions. Default: false
}

// ServiceError represents an error with a custom HTTP status code
// that can be returned from service layer methods
type ServiceError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError with the given status code and message
func NewServiceError(statusCode int, message string) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// NewServiceErrorWithCause creates a new ServiceError with the given status code, message and underlying error
func NewServiceErrorWithCause(statusCode int, message string, err error) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}
