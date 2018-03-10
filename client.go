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

var defaultClient = NewClient()

// The DefaultClient is a singleton client instance which can be configured
// and used throughout your application.
func DefaultClient() Client {
	return defaultClient
}

// SetDefaultClient allows you to replace the DefaultClient with your own
// version after you have configured it the way you wish.
func SetDefaultClient(client Client) {
	if client == nil {
		client = NewClient()
	}

	defaultClient = client
}

// UpdateDefaultClient allows you to add options to the DefaultClient.
// Changing these options will not affect derivative clients of the
// previous DefaultClient.
func UpdateDefaultClient(options ...Option) Client {
	cl := DefaultClient().With(options...)
	SetDefaultClient(cl)
	return cl
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
	conf := c.getConfig(options...)

	return c.getQueue().Enqueue(conf, p)
}

func (c *client) With(options ...Option) Client {
	return &client{
		parent:  c,
		options: options,

		queue: nil,
	}
}

func (c *client) UseSendQueue(queue SendQueue) Client {
	c.queue = queue
	return c
}

func (c *client) getQueue() SendQueue {
	if c.queue != nil {
		return c.queue
	}

	if c.parent == nil {
		return DefaultSendQueue()
	}

	return c.parent.getQueue()
}

func (c *client) fullDefaultOptions() []Option {
	if c.parent == nil {
		return c.options
	}

	return append(c.parent.fullDefaultOptions(), c.options...)
}

func (c *client) getConfig(options ...Option) *configOption {
	cnf := &configOption{}
	for _, opt := range append(c.fullDefaultOptions(), options...) {
		if oc, ok := opt.(*configOption); ok {
			cnf = oc.Merge(cnf).(*configOption)
		}
	}

	return cnf
}
