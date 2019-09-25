package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRelease() {
	cl := NewClient(
		// You can set the release when you create a client
		Release("v1.0.0"),
	)

	cl.Capture(
		// You can also set it when you send an event
		Release("v1.0.0-dev"),
	)
}

func TestRelease(t *testing.T) {
	o := Release("test")
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "release", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, "test", testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
