package sentry

import (
	"net/http"
	"os"
)

func ExampleHTTPRequest() {
	cl := NewClient(
		Release("v1.0.0"),
	)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		cl := cl.With(
			HTTPRequest(req).WithHeaders(),
		)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(404)
		res.Write([]byte(`{"error":"Not Found","message":"We could not find the route you requested, please check your URL and try again."}`))

		cl.Capture(
			Message("Route Not Found: [%s] %s", req.Method, req.URL.Path),
			Level(Warning),
		)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		cl.Capture(
			ExceptionForError(err),
			Level(Fatal),
			Extra(map[string]interface{}{
				"port": 8080,
			}),
		)

		os.Exit(1)
	}

	//Output:

}
