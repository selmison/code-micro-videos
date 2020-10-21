package errors

type Service interface {
	// ServiceError tells the transport layer whether this error should be translated into the transport format
	// or an internal error should be returned instead.
	ServiceError() bool
	// Validation tells a client that this error is related to a resource being invalid.
	// Can be used to translate the error to eg. status code.
	Validation() bool
}

type validationError struct {
	violations map[string][]string
}

func (v validationError) Error() string {
	return "invalid item"
}

func (e validationError) Violations() map[string][]string {
	return e.violations
}

// Validation tells a client that this error is related to a resource being invalid.
// Can be used to translate the error to eg. status code.
func (validationError) Validation() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (validationError) ServiceError() bool {
	return true
}
