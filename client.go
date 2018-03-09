package sentry

// A Client is responsible for letting you interact with the Sentry API.
// You can create derivative clients
type Client interface {
	// With creates a new derivative client with the provided options
	// set as part of its defaults.
	With(options ...Option) Client

	// UseSendQueue allows you to switch out the SendQueue implementation
	// used by this client. It will be copied to all future derivative
	// clients created using With().
	// Specifying `nil` as your queue will tell this client to use the
	// global DefaultSendQueue().
	UseSendQueue(queue SendQueue) Client

	// Capture will queue an event for sending to Sentry and return a
	// QueuedEvent object which can be used to keep tabs on when it is
	// actually sent, if you are curious.
	Capture(options ...Option) QueuedEvent
}

type client struct {
	parent  *client
	options []Option
	queue   SendQueue
}

// NewClient will create a new client instance with the provided
// default options and config.
func NewClient(options ...Option) Client {
	return &client{
		parent:  nil,
		options: options,
	}
}

func (c *client) Capture(options ...Option) QueuedEvent {
	p := NewPacket().SetOptions(c.fullDefaultOptions()...).SetOptions(options...)
	conf := c.getConfig(options)

	q := c.queue
	if c.queue == nil {
		q = DefaultSendQueue()
	}

	return q.Enqueue(conf, p)
}

func (c *client) With(options ...Option) Client {
	return &client{
		parent:  c,
		options: append(c.options, options...),

		queue: c.queue,
	}
}

func (c *client) UseSendQueue(queue SendQueue) Client {
	c.queue = queue
	return c
}

func (c *client) fullDefaultOptions() []Option {
	return append(c.parent.fullDefaultOptions(), c.options...)
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
