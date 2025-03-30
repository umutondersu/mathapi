package routes

// APIError represents a common interface for all API errors
type APIError interface {
	Error() string
	StatusCode() int
}

// BaseAPIError implements the APIError interface
type BaseAPIError struct {
	message    string
	statusCode int
}

func (e BaseAPIError) Error() string {
	return e.message
}

func (e BaseAPIError) String() string {
	return e.message
}

func (e BaseAPIError) StatusCode() int {
	return e.statusCode
}

// NewAPIError creates a new API error with the given message and status code
func NewAPIError(message string, statusCode int) APIError {
	return BaseAPIError{
		message:    message,
		statusCode: statusCode,
	}
}

// Predefined errors
var (
	ErrInvalidKeys     = NewAPIError("Bad Request", 400)
	ErrInvalidBOValues = NewAPIError("number1 and number2 must be numbers", 400)
	ErrInvalidSOValues = NewAPIError("numbers must be an array of numbers", 400)
	ErrDivisionByZero  = NewAPIError("Cannot divide by zero", 400)
	ErrInternalServer  = NewAPIError("Internal server error", 500)
)
