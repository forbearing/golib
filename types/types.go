package types

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
	ParamName string
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
