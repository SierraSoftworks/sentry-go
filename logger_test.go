package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleLogger() {
	cl := NewClient(
		// You can set the logger when you create your client
		Logger("root"),
	)

	cl.Capture(
		// You can also specify it when sending an event
		Logger("http"),
	)
}

func TestLogger(t *testing.T) {
	assert.NotNil(t, testGetOptionsProvider(t, Logger("")), "it should be registered as a default option")

	o := Logger("test")
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "logger", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, "test", testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
