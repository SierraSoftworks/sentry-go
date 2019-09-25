package sentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleTimestamp() {
	cl := NewClient()

	cl.Capture(
		// You can specify the timestamp when sending an event to Sentry
		Timestamp(time.Now()),
	)
}

func TestTimestamp(t *testing.T) {
	assert.NotNil(t, testGetOptionsProvider(t, Timestamp(time.Now())), "it should be registered as a default option")

	now := time.Now()
	o := Timestamp(now)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "timestamp", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, now.UTC().Format("2006-01-02T15:04:05"), testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
