package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleHTTP() {
	// You can manually populate all this request info in situations
	// where you aren't using `net/http` as your HTTP server (or don't
	// have access to the http.Request object).
	// In all other situations, you're better off using `HTTPRequest(r)`
	// and saving yourself the effort of building this up manually.
	ri := &HTTPRequestInfo{
		URL:    "http://example.com/my.url",
		Method: "GET",
		Query:  "q=test",

		Cookies: "testing=1",
		Headers: map[string]string{
			"Host": "example.com",
		},
		Env: map[string]string{
			"REMOTE_ADDR": "127.0.0.1",
			"REMOTE_PORT": "18204",
		},
		Data: map[string]interface{}{
			"awesome": true,
		},
	}

	cl := NewClient()

	ctxCl := cl.With(
		// You can provide the HTTP request context in a context-specific
		// derived client
		HTTP(ri),
	)

	ctxCl.Capture(
		// Or you can provide it when sending an event
		HTTP(ri),
	)
}

func TestHTTP(t *testing.T) {
	r := &HTTPRequestInfo{
		URL:    "http://example.com/my.url",
		Method: "GET",
		Query:  "q=test",

		Cookies: "testing=1",
	}

	assert.Equal(t, "request", r.Class(), "request info should use the correct option class")

	assert.Nil(t, HTTP(nil), "it should return nil if it receives a nil request")

	o := HTTP(r)
	assert.NotNil(t, o, "it should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")

	assert.Equal(t, "request", o.Class(), "it should use the correct option class")

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, map[string]interface{}{
			"url": "http://example.com/my.url",
			"method": "GET",
			"query_string": "q=test",
			"cookies": "testing=1",
		}, testOptionsSerialize(t, o), "it should serialize the request info correctly")
	})
}
