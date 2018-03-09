package sentry

import (
	"log"
	"net/http"
)

func ExampleHTTPRequest() {
	cl := NewClient().With(
		DSN("https://demo.getsentry.io/"),
		Release("v1.0.0"),
		Environment("production"),
	)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		event, _ := NewEventID()
		res.Header().Set("X-Sentry-ID", event)
		res.WriteHeader(404)
		res.Write([]byte("This method is not implemented yet"))

		cl.Capture(
			EventID(event),
			Message("Not Implemented"),
			Level(Warning),
			HTTPRequest(req).WithCookies().WithHeaders(),
			StackTrace().WithInternalPrefixes("github.com/SierraSoftworks/sentry-go"),
		)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	//Output:

}
