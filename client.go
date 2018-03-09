package sentry

// A Client is responsible for letting you interact with the Sentry API.
// You can create derivative clients
type Client interface {
	With(options ...Option) Client

	// Capture will queue an event for sending to Sentry and return a
	// QueuedEvent object which can be used to keep tabs on when it is
	// actually sent, if you are curious.
	Capture(options ...Option) QueuedEvent
}

type client struct {
	Parent  *client
	Options []Option

	queue *clientQueue
}

// NewClient will create a new client instance with the provided
// default options and config.
func NewClient(options ...Option) Client {
	return &client{
		Parent:  nil,
		Options: options,

		queue: defaultClientQueue(),
	}
}

func (c *client) Capture(options ...Option) QueuedEvent {
	p := NewPacket().SetOptions(c.fullDefaultOptions()...).SetOptions(options...)
	conf := c.getConfig(options)

	return c.queue.Enqueue(conf, p)
}

func (c *client) With(options ...Option) Client {
	return &client{
		Parent:  c,
		Options: append(c.Options, options...),

		queue: c.queue,
	}
}

func (c *client) fullDefaultOptions() []Option {
	return append(c.Parent.fullDefaultOptions(), c.Options...)
}

func (c *client) getConfig(options []Option) *configOption {
	cnf := &configOption{}
	for _, opt := range append(c.fullDefaultOptions(), options...) {
		if oc, ok := opt.(*configOption); ok {
			cnf = cnf.Merge(oc).(*configOption)
		}
	}

	return cnf
}
