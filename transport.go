package sentry

// Transport is the interface that any network transport must implement
// if it wishes to be used to send Sentry events
type Transport interface {
	Send(dsn string, packet Packet) error
}
