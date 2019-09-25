package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleServerName() {
	cl := NewClient(
		// You can set the logger when you create your client
		ServerName("web01"),
	)

	cl.Capture(
		// You can also specify it when sending an event
		ServerName("web01.prod"),
	)
}

func TestServerName(t *testing.T) {
	assert.NotNil(t, testGetOptionsProvider(t, ServerName("")), "it should be registered as a default option")

	o := ServerName("test")
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "server_name", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, "test", testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
