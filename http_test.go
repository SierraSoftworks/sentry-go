package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	Convey("HTTP", t, func() {
		r := &HTTPRequestInfo{
			URL:    "http://example.com/my.url",
			Method: "GET",
			Query:  "q=test",

			Cookies: "testing=1",
		}

		Convey("HTTP()", func() {
			Convey("Should return an Option", func() {
				So(HTTP(r), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if the data is nil", func() {
				So(HTTP(nil), ShouldBeNil)
			})
		})

		Convey("HTTPRequestInfo", func() {
			Convey("Should use the correct Class()", func() {
				So(r.Class(), ShouldEqual, "request")
				So(HTTP(r).Class(), ShouldEqual, "request")
			})
		})
	})
}
