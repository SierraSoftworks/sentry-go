package sentry

import "strings"

// ErrType represents an error which may contain hierarchical error information.
type ErrType string

// IsInstance will tell you whether a given error is an instance
// of this ErrType
func (e ErrType) IsInstance(err error) bool {
	return strings.Contains(err.Error(), string(e))
}

// Unwrap will unwrap this error and return the underlying error which caused
// it to be triggered.
func (e ErrType) Unwrap() error {
	return nil
}

// Error gets the error message for this ErrType
func (e ErrType) Error() string {
	return string(e)
}
