package sentry

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func Example() {
	cl := NewClient(
		// Your DSN is fetched from the $SENTRY_DSN environment
		// variable automatically. But you can override it if you
		// prefer...
		DSN("https://key:secret@example.com/sentry/1"),
		Release("v1.0.0"),

		// Your environment is fetched from $ENV/$ENVIRONMENT automatically,
		// but you can override it here if you prefer.
		Environment("example"),

		Logger("example"),
	)

	err := errors.New("something went wrong")

	// The HTTP request that was being handled when this error occurred
	var req *http.Request

	e := cl.Capture(
		Culprit("GET /api/v1/explode"),
		ExceptionForError(err),
		HTTPRequest(req).WithHeaders().WithCookies(),
	)

	if err := e.Error(); err != nil {
		fmt.Printf("Failed to send event: %s", err.Error())
	} else {
		fmt.Printf("Sent event (id: %s)\n", e.EventID())
	}
}
