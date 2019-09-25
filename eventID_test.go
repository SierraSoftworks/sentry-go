package sentry

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleEventID() {
	id, err := NewEventID()
	if err != nil {
		log.Fatalln(err)
	}

	cl := NewClient()

	ctxCl := cl.With(
		// You could set the event ID for a context specific
		// client if you wanted (but you probably shouldn't).
		EventID(id),
	)

	ctxCl.Capture(
		// The best place to set it is when you are ready to send
		// an event to Sentry.
		EventID(id),
	)
}

func TestEventID(t *testing.T) {
	id, err := NewEventID()
	assert.Nil(t, err, "creating an event ID shouldn't return an error")
	assert.Regexp(t, "^[0-9a-f]{32}$", id, "the event ID should be 32 characters long and only alphanumeric characters")

	t.Run("EventID()", func(t *testing.T) {
		assert.Nil(t, EventID("invalid"), "it should return nil if the ID is invalid")

		o := EventID(id)
		assert.NotNil(t, o, "it should return a non-nil option if the ID is valid")
		assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")

		assert.Equal(t, "event_id", o.Class(), "it should use the correct option class")

		t.Run("MarshalJSON()", func(t *testing.T) {
			assert.Equal(t, id, testOptionsSerialize(t, EventID(id)), "it should serialize to the ID")
		})
	})

	t.Run("Packet Extensions", func (t *testing.T) {
		t.Run("getEventID()", func(t *testing.T) {
			p := NewPacket()
			assert.NotNil(t, p, "the packet should not be nil")

			pp, ok := p.(*packet)
			assert.True(t, ok, "the packet should actually be a *packet")
			assert.Equal(t, "", pp.getEventID(), "it should return an empty event ID if there is no EventID option")

			p = NewPacket().SetOptions(EventID(id))
			assert.NotNil(t, p, "the packet should not be nil")

			pp, ok = p.(*packet)
			assert.True(t, ok, "the packet should actually be a *packet")
			assert.Equal(t, id, pp.getEventID(), "it should return the event ID")
		})
	})
}
