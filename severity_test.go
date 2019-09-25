package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleLevel() {
	cl := NewClient(
		// You can set the severity level when you create your client
		Level(Debug),
	)

	cl.Capture(
		// You can also specify it when sending an event
		Level(Error),
	)
}

func TestSeverity(t *testing.T) {
	assert.NotNil(t, testGetOptionsProvider(t, Level(Info)), "it should be registered as a default option")

	o := Level(Error)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "level", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, "error", testOptionsSerialize(t, o), "it should serialize to a string")
	})

	assert.EqualValues(t, Fatal, "fatal", "fatal should use the correct name")
	assert.EqualValues(t, Error, "error", "fatal should use the correct name")
	assert.EqualValues(t, Warning, "warning", "fatal should use the correct name")
	assert.EqualValues(t, Info, "info", "fatal should use the correct name")
	assert.EqualValues(t, Debug, "debug", "fatal should use the correct name")
}
