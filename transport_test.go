package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleUseTransport() {
	var myTransport Transport

	cl := NewClient(
		// You can configure the transport to be used on a client level
		UseTransport(myTransport),
	)

	cl.Capture(
		// Or for a specific event when it is sent
		UseTransport(myTransport),
	)
}

func TestTransport(t *testing.T) {
	assert.Nil(t, UseTransport(nil), "it should return nil if no transport is provided")

	tr := newHTTPTransport()
	o := UseTransport(tr)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "sentry-go.transport", o.Class(), "it should use the right option class")
	
	if assert.Implements(t, (*Option)(nil), o, "it should implement the OmitableOption interface") {
		oo := o.(OmitableOption)
		assert.True(t, oo.Omit(), "it should always return true for calls to Omit()")
	}
}

func testNewTestTransport() *testTransport {
	return &testTransport{
		ch: make(chan Packet),
	}
}

type testTransport struct {
	ch  chan Packet
	err error
}

func (t *testTransport) Send(dsn string, packet Packet) error {
	t.ch <- packet
	return t.err
}
