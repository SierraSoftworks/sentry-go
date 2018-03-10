package sentry

import (
	"fmt"
)

// A SendQueue is used by the Sentry client to coordinate the transmission
// of events. Custom queues can be used to control parallelism and circuit
// breaking as necessary.
type SendQueue interface {
	// Enqueue is called by clients wishing to send an event to Sentry.
	// It is provided with a Config and Packet and is expected to return
	// a QueuedEvent compatible object which an application can use to
	// access information about whether the event was sent successfully
	// or not.
	Enqueue(conf Config, packet Packet) QueuedEvent

	// Shutdown is called by a client that wishes to stop the flow of
	// events through a SendQueue.
	Shutdown(wait bool)
}

var (
	// The ErrSendQueueFull error is used when an attempt to enqueue a
	// new event fails as a result of no buffer space being available.
	ErrSendQueueFull = fmt.Errorf("sentry: send queue was full")

	// The ErrSendQueueShutdown error is used when an attempt to enqueue
	// a new event fails as a result of the queue having been shutdown
	// already.
	ErrSendQueueShutdown = fmt.Errorf("sentry: send queue was shutdown")
)

var defaultSendQueue SendQueue

func init() {
	defaultSendQueue = NewSequentialSendQueue(100)
}

// The DefaultSendQueue is used by all clients which have not been configured
// to use a specific send queue themselves.
func DefaultSendQueue() SendQueue {
	return defaultSendQueue
}

// SetDefaultSendQueue allows you to change the default queue implementation
// used to send events to Sentry.
func SetDefaultSendQueue(queue SendQueue) {
	defaultSendQueue = queue
}
