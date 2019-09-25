package sentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequentialSendQueue(t *testing.T) {
	q := NewSequentialSendQueue(10)
	require.NotNil(t, q, "the queue should not be nil")
	assert.Implements(t, (*SendQueue)(nil), q, "it should implement the SendQueue interface")
	defer q.Shutdown(true)

	require.IsType(t, &sequentialSendQueue{}, q, "it should actually be a *sequentialSendQueue")

	t.Run("Send()", func(t *testing.T) {
		dsn := "http://user:pass@example.com/sentry/1"
		transport := testNewTestTransport()
		require.NotNil(t, transport, "the transport should not be nil")

		cl := NewClient(DSN(dsn), UseTransport(transport))
		require.NotNil(t, cl, "the client should not be nil")

		cfg, ok := cl.(Config)
		require.True(t, ok, "the client should implement the Config interface")

		p := NewPacket()
		require.NotNil(t, p, "the packet should not be nil")

		t.Run("Normal", func(t *testing.T) {
			q := NewSequentialSendQueue(10)
			require.NotNil(t, q, "the queue should not be nil")
			defer q.Shutdown(true)

			e := q.Enqueue(cfg, p)
			require.NotNil(t, e, "the event should not be nil")

			select {
			case pp := <-transport.ch:
				assert.Equal(t, p, pp, "the packet which was sent should match the packet which was enqueued")
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timed out waiting for send")
			}

			select {
			case err, ok := <-e.WaitChannel():
				assert.False(t, ok, "the channel should have been closed")
				assert.NoError(t, err, "there should have been no error sending the event")
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timed out waiting for event completion")
			}
		})

		t.Run("QueueFull", func(t *testing.T) {
			q := NewSequentialSendQueue(0)
			require.NotNil(t, q, "the queue should not be nil")
			defer q.Shutdown(true)

			// Give the queue time to start
			time.Sleep(1 * time.Millisecond)

			e1 := q.Enqueue(cfg, p)
			require.NotNil(t, e1, "the event should not be nil")

			e2 := q.Enqueue(cfg, p)
			require.NotNil(t, e2, "the event should not be nil")

			select {
			case pp := <-transport.ch:
				assert.Equal(t, p, pp, "the packet which was sent should match the packet which was enqueued")
				assert.Nil(t, e1.Error(), "")
			case err, ok := <-e1.WaitChannel():
				assert.False(t, ok, "the channel should have been closed")
				assert.NoError(t, err, "there should have been no error sending the event")
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timed out waiting for send")
			}

			select {
			case <-transport.ch:
				t.Error("the transport should never have received the event for sending")
			case err, ok := <-e2.WaitChannel():
				assert.True(t, ok, "the channel should not have been closed prematurely")
				assert.True(t, ErrSendQueueFull.IsInstance(err), "the error should be of type ErrSendQueueFull")
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timed out waiting for event completion")
			}
		})

		t.Run("Shutdown", func(t *testing.T) {
			q := NewSequentialSendQueue(10)
			require.NotNil(t, q, "the queue should not be nil")

			// Shutdown the queue
			q.Shutdown(true)

			e := q.Enqueue(cfg, p)
			require.NotNil(t, e, "the event should not be nil")

			select {
			case <-transport.ch:
				t.Error("the transport should never have received the event for sending")
			case err, ok := <-e.WaitChannel():
				assert.True(t, ok, "the channel should not have been closed prematurely")
				assert.True(t, ErrSendQueueShutdown.IsInstance(err), "the error should be of type ErrSendQueueShutdown")
			case <-time.After(100 * time.Millisecond):
				t.Fatal("timed out waiting for event completion")
			}
		})
	})

	t.Run("Shutdown()", func(t *testing.T) {
		q := NewSequentialSendQueue(10)

		// It should be safe to call this repeatedly
		q.Shutdown(true)
		q.Shutdown(true)
	})
}
