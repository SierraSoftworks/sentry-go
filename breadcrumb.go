package sentry

import "time"

// A Breadcrumb keeps track of an action which took place in the application
// leading up to an event.
type Breadcrumb interface {
	WithMessage(msg string) Breadcrumb
	WithCategory(cat string) Breadcrumb
	WithLevel(s Severity) Breadcrumb
	WithTimestamp(ts time.Time) Breadcrumb
}

func newBreadcrumb(typename string, data map[string]interface{}) *breadcrumb {
	if typename == "default" {
		typename = ""
	}

	return &breadcrumb{
		Timestamp: time.Now().UTC().Unix(),
		Type:      typename,
		Data:      data,
	}
}

type breadcrumb struct {
	Timestamp int64                  `json:"timestamp"`
	Type      string                 `json:"type,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Level     Severity               `json:"level,omitempty"`
}

func (b *breadcrumb) WithMessage(msg string) Breadcrumb {
	b.Message = msg
	return b
}

func (b *breadcrumb) WithCategory(cat string) Breadcrumb {
	b.Category = cat
	return b
}

func (b *breadcrumb) WithTimestamp(ts time.Time) Breadcrumb {
	b.Timestamp = ts.UTC().Unix()
	return b
}

func (b *breadcrumb) WithLevel(s Severity) Breadcrumb {
	b.Level = s
	return b
}
