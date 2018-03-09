package sentry

import "fmt"

type message struct {
	Message   string
	Params    []interface{}
	Formatted string
}

func (m *message) Class() string {
	return "sentry.interfaces.Message"
}

// Message generates a new message entry for Sentry, optionally
// using a format string with standard fmt.Sprintf params.
func Message(format string, params ...interface{}) Option {
	if len(params) == 0 {
		return &message{
			Message: format,
		}
	}

	return &message{
		Message:   format,
		Params:    params,
		Formatted: fmt.Sprintf(format, params...),
	}
}
