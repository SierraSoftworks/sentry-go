package sentry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleQueuedEvent() {
	cl := NewClient()

	e := cl.Capture(
		Message("Example Event"),
	)

	// You can then wait on the event to be sent
	e.Wait()

	// Or you can use the WaitChannel if you want support for timeouts
	select {
	case err := <-e.WaitChannel():
		if err != nil {
			fmt.Println("failed to send event: ", err)
		} else {
			// You can also get the EventID for this event
			fmt.Println("sent event: ", e.EventID())
		}
	case <-time.After(time.Second):
		fmt.Println("timed out waiting for send")
	}
}

func ExampleQueuedEventInternal() {
	// If you're implementing your own send queue, you will want to use
	// the QueuedEventInternal to control when events are finished and
	// to access the packet and config related to them.

	cl := NewClient()
	e := cl.Capture()

	if ei, ok := e.(QueuedEventInternal); ok {
		// Get the packet for the event
		p := ei.Packet()

		// Get the config for the event
		cfg := ei.Config()

		// Use the configured transport to send the packet
		err := cfg.Transport().Send(cfg.DSN(), p)

		// Complete the event (with the error, if not nil)
		ei.Complete(err)
	}
}

func TestQueuedEvent(t *testing.T) {
	id, err := NewEventID()
	require.Nil(t, err, "there should be no errors creating an event ID")

	cl := NewClient(DSN(""))
	require.NotNil(t, cl, "the client should not be nil")

	cfg, ok := cl.(Config)
	require.True(t, ok, "the client should implement the Config interface")

	p := NewPacket().SetOptions(EventID(id))
	require.NotNil(t, p, "the packet should not be nil")

	t.Run("NewQueuedEvent()", func(t *testing.T) {
		e := NewQueuedEvent(cfg, p)
		require.NotNil(t, e, "the event should not be nil")
		assert.Implements(t, (*QueuedEvent)(nil), e, "it should implement the QueuedEvent interface")

		ei, ok := e.(*queuedEvent)
		require.True(t, ok, "it should actually be a *queuedEvent")
		assert.Same(t, cfg, ei.cfg, "it should use the same config provider")
		assert.Same(t, p, ei.packet, "it should use the same packet")

		t.Run("EventID()", func(t *testing.T) {
			assert.Equal(t, id, e.EventID(), "it should have the right event ID")

			assert.Empty(t, NewQueuedEvent(cfg, nil).EventID(), "it should have an empty EventID for an invalid packet")
		})

		t.Run("Wait()", func(t *testing.T) {
			cases := []struct {
				Name            string
				Waiter          func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
				PreWaiterStart  func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
				PostWaiterStart func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
			}{
				{
					Name: "SuccessSlow",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ch <- struct{}{}
						ei.Wait()
						assert.Nil(t, ei.Error(), "there should have been no error raised")
					},
					PostWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						<-ch
						ei.Complete(nil)
					},
				},
				{
					Name: "SuccessFast",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Wait()
						assert.Nil(t, ei.Error(), "there should have been no error raised")
					},
					PreWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Complete(nil)
					},
				},
				{
					Name: "FailSlow",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ch <- struct{}{}
						ei.Wait()
						assert.EqualError(t, ei.Error(), "test error", "there should have been an error raised")
					},
					PostWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						<-ch
						ei.Complete(fmt.Errorf("test error"))
					},
				},
				{
					Name: "FailFast",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Wait()
						assert.EqualError(t, ei.Error(), "test error", "there should have been an error raised")
					},
					PreWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Complete(fmt.Errorf("test error"))
					},
				},
			}

			for _, tc := range cases {
				tc := tc
				t.Run(tc.Name, func(t *testing.T) {
					e := NewQueuedEvent(cfg, p)
					require.NotNil(t, e, "the event should not be nil")

					require.Implements(t, (*QueuedEventInternal)(nil), e, "it should implement the QueuedEventInternal interface")
					ei := e.(QueuedEventInternal)

					ch := make(chan struct{})
					defer close(ch)

					if tc.PreWaiterStart != nil {
						tc.PreWaiterStart(t, ch, ei)
					}

					go func() {
						if tc.Waiter != nil {
							tc.Waiter(t, ch, ei)
						}

						ch <- struct{}{}
					}()

					if tc.PostWaiterStart != nil {
						tc.PostWaiterStart(t, ch, ei)
					}

					select {
					case <-ch:
					case <-time.After(100 * time.Millisecond):
						t.Error("timed out after 100ms with no response")
					}
				})
			}
		})

		t.Run("WaitChannel()", func(t *testing.T) {
			cases := []struct {
				Name     string
				Complete func(ei QueuedEventInternal)
				Error    error
			}{
				{"SucceedFast", func(ei QueuedEventInternal) { ei.Complete(nil) }, nil},
				{"SucceedSlow", func(ei QueuedEventInternal) { go func() { ei.Complete(nil) }() }, nil},
				{"FailFast", func(ei QueuedEventInternal) { ei.Complete(fmt.Errorf("test error")) }, fmt.Errorf("test error")},
				{"FailSlow", func(ei QueuedEventInternal) { go func() { ei.Complete(fmt.Errorf("test error")) }() }, fmt.Errorf("test error")},
			}

			for _, tc := range cases {
				tc := tc
				t.Run(tc.Name, func(t *testing.T) {
					e := NewQueuedEvent(cfg, p)
					require.NotNil(t, e, "the event should not be nil")

					require.Implements(t, (*QueuedEventInternal)(nil), e, "it should implement the QueuedEventInternal interface")
					ei := e.(QueuedEventInternal)

					tc.Complete(ei)

					select {
					case err := <-e.WaitChannel():
						if tc.Error != nil {
							assert.EqualError(t, err, tc.Error.Error(), "the right error should have been raised")
						} else {
							assert.NoError(t, err, "no error should have been raised")
						}
					case <-time.After(100 * time.Millisecond):
						t.Error("timeout after 100ms")
					}
				})
			}
		})

		t.Run("Error()", func(t *testing.T) {
			cases := []struct {
				Name            string
				Waiter          func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
				PreWaiterStart  func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
				PostWaiterStart func(t *testing.T, ch chan struct{}, ei QueuedEventInternal)
			}{
				{
					Name: "SuccessSlow",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ch <- struct{}{}
						assert.Nil(t, ei.Error(), "there should have been no error raised")
					},
					PostWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						<-ch
						ei.Complete(nil)
					},
				},
				{
					Name: "SuccessFast",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						assert.Nil(t, ei.Error(), "there should have been no error raised")
					},
					PreWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Complete(nil)
					},
				},
				{
					Name: "FailSlow",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ch <- struct{}{}
						assert.EqualError(t, ei.Error(), "test error", "there should have been an error raised")
					},
					PostWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						<-ch
						ei.Complete(fmt.Errorf("test error"))
					},
				},
				{
					Name: "FailFast",
					Waiter: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						assert.EqualError(t, ei.Error(), "test error", "there should have been an error raised")
					},
					PreWaiterStart: func(t *testing.T, ch chan struct{}, ei QueuedEventInternal) {
						ei.Complete(fmt.Errorf("test error"))
					},
				},
			}

			for _, tc := range cases {
				tc := tc
				t.Run(tc.Name, func(t *testing.T) {
					e := NewQueuedEvent(cfg, p)
					require.NotNil(t, e, "the event should not be nil")

					require.Implements(t, (*QueuedEventInternal)(nil), e, "it should implement the QueuedEventInternal interface")
					ei := e.(QueuedEventInternal)

					ch := make(chan struct{})
					defer close(ch)

					if tc.PreWaiterStart != nil {
						tc.PreWaiterStart(t, ch, ei)
					}

					go func() {
						if tc.Waiter != nil {
							tc.Waiter(t, ch, ei)
						}

						ch <- struct{}{}
					}()

					if tc.PostWaiterStart != nil {
						tc.PostWaiterStart(t, ch, ei)
					}

					select {
					case <-ch:
					case <-time.After(100 * time.Millisecond):
						t.Error("timed out after 100ms with no response")
					}
				})
			}
		})
	})

	t.Run("Complete()", func(t *testing.T) {
		e := NewQueuedEvent(cfg, p)
		require.NotNil(t, e, "the event should not be nil")

		require.Implements(t, (*QueuedEventInternal)(nil), e, "it should implement the QueuedEventInternal interface")
		ei := e.(QueuedEventInternal)

		ei.Complete(fmt.Errorf("test error"))
		assert.EqualError(t, e.Error(), "test error", "it should set the error correctly")

		ei.Complete(nil)
		assert.NotNil(t, e.Error(), "it shouldn't modify the status of the event after it has been set")
	})
}
