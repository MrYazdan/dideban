package types

// Error represents error information in API responses
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse creates an error API response
func ErrorResponse(code, message, details string) Response {
	return Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(details string) Response {
	return ErrorResponse("VALIDATION_ERROR", "Invalid input data", details)
}

// AuthenticationErrorResponse creates an authentication error response
func AuthenticationErrorResponse(details string) Response {
	return ErrorResponse("AUTHENTICATION_ERROR", "Authentication failed", details)
}

// AuthorizationErrorResponse creates an authorization error response
func AuthorizationErrorResponse(details string) Response {
	return ErrorResponse("AUTHORIZATION_ERROR", "Access denied", details)
}

// NotFoundErrorResponse creates a not found error response
func NotFoundErrorResponse(resource string) Response {
	return ErrorResponse("NOT_FOUND", "Resource not found", resource+" not found")
}

// ConflictErrorResponse creates a conflict error response
func ConflictErrorResponse(details string) Response {
	return ErrorResponse("CONFLICT", "Resource conflict", details)
}

// InternalErrorResponse creates an internal server error response
func InternalErrorResponse(details string) Response {
	return ErrorResponse("INTERNAL_ERROR", "Internal server error", details)
}

// TimeoutErrorResponse create a timeout server error response
func TimeoutErrorResponse(details string) Response {
	return ErrorResponse("TIMEOUT", "Request timeout", details)
}
