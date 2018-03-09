package sentry

import "sync"

// A QueuedEvent allows you to track the status of sending
// an event to Sentry.
type QueuedEvent interface {
	EventID() string
	Wait() QueuedEvent
	Error() error
}

type queuedEvent struct {
	conf     *configOption
	packet   Packet
	complete bool
	err      error

	ch chan error
}

func (e *queuedEvent) EventID() string {
	if packet, ok := e.packet.(*packet); ok {
		return packet.getEventID()
	}

	return ""
}

func (e *queuedEvent) Wait() QueuedEvent {
	if e.complete {
		return e
	}

	err, ok := <-e.ch

	if ok {
		e.err = err
		e.complete = true
	} else {
		e.complete = true
	}

	return e
}

func (e *queuedEvent) Error() error {
	return e.Wait().(*queuedEvent).err
}

// TODO: Implement this
func defaultClientQueue() *clientQueue {
	return nil
}

func newClientQueue(queueSize int) *clientQueue {
	sendQueue := make(chan *queuedEvent, queueSize)
	closeCh := make(chan struct{})

	q := &clientQueue{
		sendQueue: sendQueue,
		close:     closeCh,
	}

	go q.worker()

	return q
}

type clientQueue struct {
	Transport Transport

	sendQueue chan *queuedEvent
	close     chan struct{}
	closed    bool
	mu        sync.Mutex
}

func (q *clientQueue) WithTransport(transport Transport) *clientQueue {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Transport = transport
	return q
}

func (q *clientQueue) Enqueue(conf *configOption, packet Packet) QueuedEvent {
	errs := make(chan error)
	event := queuedEvent{
		conf:   conf,
		packet: packet,
		ch:     errs,
	}

	select {
	case q.sendQueue <- &event:
		// Sent packet

	default:
		// Queue was full
	}

	return &event
}

func (q *clientQueue) worker() {
	for {
		select {
		case event := <-q.sendQueue:
			err := q.Transport.Send(event.conf.DSN(), event.packet)
			if err != nil {
				event.ch <- err
			}
			close(event.ch)

		case <-q.close:
			q.mu.Lock()
			defer q.mu.Unlock()

			q.closed = true
			close(q.close)
			close(q.sendQueue)
			return
		}
	}
}
