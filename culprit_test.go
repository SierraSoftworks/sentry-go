package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleCulprit() {
	cl := NewClient(
		// You can set this when creating your client
		Culprit("example"),
	)

	cl.Capture(
		// Or you can set it when sending an event
		Culprit("example"),
	)
}

func TestCulprit(t *testing.T) {
	c := Culprit("test")
	assert.Implements(t, (*Option)(nil), c, "it should implement the Option interface")
	assert.Equal(t, "culprit", c.Class(), "it should use the correct option class")
	assert.Equal(t, "test", testOptionsSerialize(t, c), "it should serialize to the value which was passed in the constructor")
}
