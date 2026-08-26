package deadsimple

import "fmt"

// Error represents an API error.
type Error struct {
	Message    string
	Code       string
	StatusCode int
}

func (e *Error) Error() string {
	return fmt.Sprintf("deadsimple: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
}

// IsNotFound returns true if the error is a 404.
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 404
	}
	return false
}

// IsRateLimit returns true if the error is a 429.
func IsRateLimit(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 429
	}
	return false
}

// IsValidation returns true if the error is a 400.
func IsValidation(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 400
	}
	return false
}

// IsAuth returns true if the error is a 401.
func IsAuth(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 401
	}
	return false
}
