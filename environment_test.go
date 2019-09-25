package sentry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleEnvironment() {
	cl := NewClient(
		// You can configure your environment at the client level
		Environment("development"),
	)

	cl.Capture(
		// ...or at the event level
		Environment("prod"),
	)
}

func TestEnvironment(t *testing.T) {
	o := Environment("testing")
	assert.NotNil(t, o, "it should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the option interface")
	assert.Equal(t, "environment", o.Class(), "it should use the correct option class")

	t.Run("No Environment", func(t *testing.T) {
		os.Unsetenv("ENV")
		os.Unsetenv("ENVIRONMENT")

		opt := testGetOptionsProvider(t, &environmentOption{})
		assert.Nil(t, opt, "it should not be registered as a default option provider")
	})

	t.Run("$ENV=...", func(t *testing.T){
		os.Setenv("ENV", "testing")
		defer os.Unsetenv("ENV")

		opt := testGetOptionsProvider(t, &environmentOption{})
		assert.NotNil(t, opt, "it should be registered with the default option providers")
		assert.IsType(t, &environmentOption{}, opt, "it should actually be an *environmentOption")

		oo := opt.(*environmentOption)
		assert.Equal(t, "testing", oo.env, "it should set the environment to the same value as the $ENV variable")
	})

	t.Run("$ENVIRONMENT=...", func(t *testing.T){
		os.Setenv("ENVIRONMENT", "testing")
		defer os.Unsetenv("ENVIRONMENT")

		opt := testGetOptionsProvider(t, &environmentOption{})
		assert.NotNil(t, opt, "it should be registered with the default option providers")
		assert.IsType(t, &environmentOption{}, opt, "it should actually be an *environmentOption")

		oo := opt.(*environmentOption)
		assert.Equal(t, "testing", oo.env, "it should set the environment to the same value as the $ENVIRONMENT variable")
	})

	t.Run("MarshalJSON()", func(t *testing.T) {
		s := testOptionsSerialize(t, o)
		assert.Equal(t, "testing", s, "it should serialize to the name of the environment")
	})
}
