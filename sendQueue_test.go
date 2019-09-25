package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleUseSendQueue() {
	cl := NewClient(
		// You can override the send queue on your root client
		// All of its derived clients will inherit this queue
		UseSendQueue(NewSequentialSendQueue(10)),
	)

	cl.With(
		// Or you can override it on a derived client
		UseSendQueue(NewSequentialSendQueue(10)),
	)
}

func TestSendQueue(t *testing.T) {
	assert.Nil(t, UseSendQueue(nil), "it should return nil if no transport is provided")

	q := NewSequentialSendQueue(0)
	o := UseSendQueue(q)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "sentry-go.sendqueue", o.Class(), "it should use the right option class")
	
	if assert.Implements(t, (*Option)(nil), o, "it should implement the OmitableOption interface") {
		oo := o.(OmitableOption)
		assert.True(t, oo.Omit(), "it should always return true for calls to Omit()")
	}
}
