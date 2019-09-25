package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSDKOption(t *testing.T) {
	o := testGetOptionsProvider(t, &sdkOption{})
	assert.NotNil(t, o, "it should be registered as a default option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "sdk", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, map[string]interface{}{
			"integrations": []interface{}{},
			"name": "SierraSoftworks/sentry-go",
			"version": version,
		}, testOptionsSerialize(t, o), "it should serialize to a string")
	})
}
