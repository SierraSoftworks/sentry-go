package sentry

// Transport is the interface that any network transport must implement
// if it wishes to be used to send Sentry events
type Transport interface {
	Send(dsn string, packet Packet) error
}

// UseTransport allows you to control which transport is used to
// send events for a specific client or packet.
func UseTransport(transport Transport) Option {
	if transport == nil {
		return nil
	}

	return &configOption{
		transport: transport,
	}
}

var defaultTransport Transport

func init() {
	defaultTransport = newHTTPTransport()
}

// DefaultTransport retrieves the transport that will be used
// by default for all new ClientQueues.
func DefaultTransport() Transport {
	return defaultTransport
}

// SetDefaultTransport allows you to change the transport that
// is used by default for all new packets.
func SetDefaultTransport(transport Transport) {
	defaultTransport = transport
}
