package sentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleBreadcrumb() {
	b := DefaultBreadcrumbs().NewDefault(nil)

	// You can set the severity level for the breadcrumb
	b.WithLevel(Error)

	// You can configure the category that the breadcrumb belongs to
	b.WithCategory("auth")

	// You can also specify a message describing the breadcrumb
	b.WithMessage("User's credentials were invalid")

	// And if you need to change the timestamp, you can do that too
	b.WithTimestamp(time.Now())

	// All together now!
	DefaultBreadcrumbs().
		NewDefault(nil).
		WithLevel(Error).
		WithCategory("auth").
		WithMessage("User's credentials were invalid").
		WithTimestamp(time.Now())
}

func TestBreadcrumb(t *testing.T) {
	data := map[string]interface{}{
		"test": true,
	}

	t.Run("newBreadcrumb", func(t *testing.T) {
		b := newBreadcrumb("default", data)

		if assert.NotNil(t, b) {
			assert.Implements(t, (*Breadcrumb)(nil), b)
			assert.Equal(t, "", b.Type, "It should set the correct type")
			assert.NotEqual(t, 0, b.Timestamp, "It should set the timestamp")
			assert.Equal(t, data, b.Data, "It should set the correct data")
		}
	})

	t.Run("WithMessage()", func(t *testing.T) {
		b := newBreadcrumb("default", data)

		if assert.NotNil(t, b) {
			bb := b.WithMessage("test")
			assert.Equal(t, b, bb, "It should return the breadcrumb for chaining")
			assert.Equal(t, "test", b.Message)
		}
	})

	t.Run("WithCategory()", func(t *testing.T) {
		b := newBreadcrumb("default", data)

		if assert.NotNil(t, b) {
			bb := b.WithCategory("test")
			assert.Equal(t, b, bb, "It should return the breadcrumb for chaining")
			assert.Equal(t, "test", b.Category)
		}
	})

	t.Run("WithLevel()", func(t *testing.T) {
		b := newBreadcrumb("default", data)

		if assert.NotNil(t, b) {
			bb := b.WithLevel(Error)
			assert.Equal(t, b, bb, "It should return the breadcrumb for chaining")
			assert.Equal(t, Error, b.Level)
		}
	})

	t.Run("WithTimestamp()", func(t *testing.T) {
		b := newBreadcrumb("default", data)
		now := time.Now()

		if assert.NotNil(t, b) {
			bb := b.WithTimestamp(now)
			assert.Equal(t, b, bb, "It should return the breadcrumb for chaining")
			assert.Equal(t, now.UTC().Unix(), b.Timestamp)
		}
	})
}
