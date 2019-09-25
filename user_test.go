package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleUser() {
	user := UserInfo{
		ID:        "17ba08f7cc89a912bf812918",
		Email:     "test@example.com",
		Username:  "Test User",
		IPAddress: "127.0.0.1",
		Extra: map[string]string{
			"role": "Tester",
		},
	}

	cl := NewClient(
		// You can specify your user when you create your client
		User(&user),
	)

	cl.Capture(
		// Or when you send an event to Sentry
		User(&user),
	)
}

func TestUser(t *testing.T) {
	user := UserInfo{
		ID:        "17ba08f7cc89a912bf812918",
		Email:     "test@example.com",
		Username:  "Test User",
		IPAddress: "127.0.0.1",
		Extra: map[string]string{
			"role": "Tester",
		},
	}

	fields := map[string]interface{}{
		"id":         "17ba08f7cc89a912bf812918",
		"email":      "test@example.com",
		"username":   "Test User",
		"ip_address": "127.0.0.1",
		"role":       "Tester",
	}

	o := User(&user)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "user", o.Class(), "it should use the right option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, fields, testOptionsSerialize(t, o), "it should serialize to the right fields")
	})
}
