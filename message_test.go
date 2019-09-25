package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleMessage() {
	cl := NewClient()

	cl.Capture(
		// You can either use just a simple message
		Message("this is a simple message"),
	)

	cl.Capture(
		// Or you can provide formatting entries as you would with
		// fmt.Sprintf() calls.
		Message("this is a %s message (%d/7 would use again)", "formatted", 5),
	)
}

func TestMessage(t *testing.T) {
	o := Message("test")
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "sentry.interfaces.Message", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, map[string]interface{}{"message":"test"}, testOptionsSerialize(t, o), "it should serialize to an object")
	})

	t.Run("parameters", func(t *testing.T) {
		o := Message("this is a %s", "test")
		assert.NotNil(t, o, "should not return a nil option")

		mi, ok := o.(*messageOption)
		assert.True(t, ok, "it should actually be a *messageOption")
		assert.Equal(t, "this is a %s", mi.Message, "it should use the right message")
		assert.Equal(t, []interface{}{"test"}, mi.Params, "it should have the correct parameters")
		assert.Equal(t, "this is a test", mi.Formatted, "it should format the message when requested")
	})
}
