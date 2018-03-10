package sentry

// A Config allows you to control how events are sent to Sentry.
// It is usually populated through the standard build pipeline
// through the DSN() and UseTransport() options.
type Config interface {
	DSN() string
	Transport() Transport
}

type configOption struct {
	dsn       *string
	transport Transport
}

func (o *configOption) Class() string {
	return "sentry-go.config"
}

func (o *configOption) Ommit() bool {
	return true
}

func (o *configOption) Clone() *configOption {
	return &configOption{
		dsn:       o.dsn,
		transport: o.transport,
	}
}

func (o *configOption) Merge(old Option) Option {
	if old, ok := old.(*configOption); ok {
		c := old.Clone()

		if o.dsn != nil {
			c.dsn = o.dsn
		}

		if o.transport != nil {
			c.transport = o.transport
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

func (o *configOption) Transport() Transport {
	if o.transport == nil {
		return DefaultTransport()
	}

	return o.transport
}
