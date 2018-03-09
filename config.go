package sentry

import (
	"os"
)

func init() {
	addDefaultOptionProvider(func() Option {
		return DSN(os.Getenv("SENTRY_DSN"))
	})
}

// DSN lets you specify the unique Sentry DSN used to submit events for
// your application. Specifying an empty DSN will disable the client.
func DSN(dsn string) Option {
	return &configOption{
		dsn: &dsn,
	}
}

type configOption struct {
	dsn *string
}

func (o *configOption) Class() string {
	return "sentry-go.config"
}

func (o *configOption) Ommit() bool {
	return true
}

func (o *configOption) Clone() *configOption {
	return &configOption{
		dsn: o.dsn,
	}
}

func (o *configOption) Merge(old Option) Option {
	if old, ok := old.(*configOption); ok {
		c := old.Clone()

		if o.dsn != nil {
			c.dsn = o.dsn
		}

		return c
	}

	return o
}

func (o *configOption) DSN() string {
	if o.dsn == nil {
		return ""
	}

	return *o.dsn
}
