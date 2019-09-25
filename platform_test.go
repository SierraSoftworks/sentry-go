package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePlatform() {
	cl := NewClient(
		// You can set the platform at a client level
		Platform("go"),
	)

	cl.Capture(
		// Or override it when sending the event
		Platform("go"),
	)
}

func TestPlatform(t *testing.T) {
	assert.NotNil(t, testGetOptionsProvider(t, Platform("go")), "it should be registered as a default option")

	o := Platform("go")
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "platform", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, "go", testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
