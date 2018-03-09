package sentry

import "encoding/json"

func init() {
	// Configure the default severity level as Error
	addDefaultOptionProvider(func() Option {
		return Level(Error)
	})
}

// Level is used to set the severity level of an event before it
// is sent to Sentry
func Level(severity Severity) Option {
	return &levelOption{severity}
}

type levelOption struct {
	severity Severity
}

func (o *levelOption) Class() string {
	return "level"
}

func (o *levelOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.severity)
}

// Severity represents a Sentry event severity (ranging from debug to fatal)
type Severity string

// Fatal represents exceptions which result in the application exiting fatally
var Fatal = Severity("fatal")

// Error represents exceptions which break the expected application flow
var Error = Severity("error")

// Warning represents events which are abnormal but do not prevent the application
// from operating correctly
var Warning = Severity("warning")

// Info is used to expose information about events which occur during normal
// operation of the application
var Info = Severity("info")

// Debug is used to expose verbose information about events which occur during
// normal operation of the application
var Debug = Severity("debug")
