// +build go1.12

package sentry

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModulesWithGomod(t *testing.T) {
	_, ok := debug.ReadBuildInfo()
	if ok {
		assert.NotNil(t, testGetOptionsProvider(t, Modules(map[string]string{"test": "correct"})), "it should be registered as a default provider")
	}
}
