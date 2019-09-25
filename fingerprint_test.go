package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleFingerprint() {
	cl := NewClient()

	cl.Capture(
		// You can specify a fingerprint that extends the default behaviour
		Fingerprint("{{ default }}", "http://example.com/my.url"),

		// Or you can define your own
		Fingerprint("myrpc", "POST", "/foo.bar"),
	)
}

func TestFingerprint(t *testing.T) {
	o := Fingerprint("test")
	assert.NotNil(t, o, "it should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "fingerprint", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func (t *testing.T) {
		assert.Equal(t, []interface{}{"test"}, testOptionsSerialize(t, o), "it should serialize as a list of fingerprint keys")
	})
}
