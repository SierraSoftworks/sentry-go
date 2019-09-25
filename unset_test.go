package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleUnset() {
	cl := NewClient(
		// You can remove specific default fields from your final packet if you do
		// not wish to send them.
		Unset("runtime"),
	)

	cl.Capture(
		// You can also remove things that you may have added later
		Unset("message"),
	)
}

func TestUnset(t *testing.T) {
	o := Unset("runtime")
	assert.Equal(t, o.Class(), "runtime", "it should use the correct class name")

	o = Unset("device")
	assert.Equal(t, o.Class(), "device", "it should use the correct class name")
	
	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, nil, testOptionsSerialize(t, o), "it should serialize to nil")
	})

	if assert.Implements(t, (*AdvancedOption)(nil), o, "it should implement the AdvancedOption interface") {
		p := map[string]Option{
			"level": Level(Error),
			"release": Release("1.0.0"),
		}

		Unset("level").(AdvancedOption).Apply(p)
		assert.NotContains(t, "level", p, "it should remove the property from the packet")
	}
}
