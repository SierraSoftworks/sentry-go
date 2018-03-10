package sentry

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleHTTPRequest() {
	cl := NewClient(
		Release("v1.0.0"),
	)

	// Add your 404 handler to the default mux
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		cl := cl.With(
			// Set the HTTP request context for your request's client
			HTTPRequest(req).WithHeaders(),
		)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(404)
		res.Write([]byte(`{"error":"Not Found","message":"We could not find the route you requested, please check your URL and try again."}`))

		// Capture the problem using your request's client
		cl.Capture(
			Message("Route Not Found: [%s] %s", req.Method, req.URL.Path),
			Level(Warning),
		)
	})
}

func TestHTTPRequest(t *testing.T) {
	Convey("HTTPRequest", t, func() {
		r, err := http.NewRequest("GET", "https://example.com/test?testing=1&password=test", nil)
		So(err, ShouldBeNil)

		r.RemoteAddr = "127.0.0.1:12835"
		r.Header.Set("Host", "example.com")
		r.Header.Set("X-Forwarded-Proto", "https")
		r.Header.Set("Cookie", "testing=1")
		r.Header.Set("X-Testing", "1")

		Convey("HTTPRequest()", func() {
			Convey("Should return an Option", func() {
				So(HTTPRequest(r), ShouldImplement, (*Option)(nil))
			})

			Convey("Should not return nil if request is nil", func() {
				So(HTTPRequest(nil), ShouldNotBeNil)
			})
		})

		Convey("Should use the correct Class()", func() {
			So(HTTPRequest(r).Class(), ShouldEqual, "request")
		})

		Convey("Omit()", func() {
			Convey("Should return false with a valid request", func() {
				So(HTTPRequest(r).(*httpRequestOption).Omit(), ShouldBeFalse)
			})

			Convey("Should return true if no request was provided", func() {
				So(HTTPRequest(nil).(*httpRequestOption).Omit(), ShouldBeTrue)
			})
		})

		Convey("buildData()", func() {
			Convey("With the default config", func() {
				opt := HTTPRequest(r)
				hr, ok := opt.(*httpRequestOption)
				So(ok, ShouldBeTrue)

				d := hr.buildData()
				So(d, ShouldNotBeNil)

				So(d.Method, ShouldEqual, "GET")
				So(d.URL, ShouldEqual, "https://example.com/test")
				So(d.Query, ShouldEqual, url.Values{
					"testing":  {"1"},
					"password": {"********"},
				}.Encode())

				So(d.Data, ShouldBeNil)
				So(d.Headers, ShouldResemble, map[string]string{})
				So(d.Env, ShouldResemble, map[string]string{})
				So(d.Cookies, ShouldEqual, "")
			})

			Convey("With cookies enabled", func() {
				opt := HTTPRequest(r).WithCookies()
				hr, ok := opt.(*httpRequestOption)
				So(ok, ShouldBeTrue)

				d := hr.buildData()
				So(d, ShouldNotBeNil)

				So(d.Method, ShouldEqual, "GET")
				So(d.URL, ShouldEqual, "https://example.com/test")
				So(d.Query, ShouldEqual, url.Values{
					"testing":  {"1"},
					"password": {"********"},
				}.Encode())

				So(d.Data, ShouldBeNil)
				So(d.Headers, ShouldResemble, map[string]string{})
				So(d.Env, ShouldResemble, map[string]string{})
				So(d.Cookies, ShouldEqual, "testing=1")
			})

			Convey("With headers enabled", func() {
				opt := HTTPRequest(r).WithHeaders()
				hr, ok := opt.(*httpRequestOption)
				So(ok, ShouldBeTrue)

				d := hr.buildData()
				So(d, ShouldNotBeNil)

				So(d.Method, ShouldEqual, "GET")
				So(d.URL, ShouldEqual, "https://example.com/test")
				So(d.Query, ShouldEqual, url.Values{
					"testing":  {"1"},
					"password": {"********"},
				}.Encode())

				So(d.Data, ShouldBeNil)
				So(d.Headers, ShouldResemble, map[string]string{
					"Host":              "example.com",
					"Cookie":            "testing=1",
					"X-Testing":         "1",
					"X-Forwarded-Proto": "https",
				})
				So(d.Env, ShouldResemble, map[string]string{})
				So(d.Cookies, ShouldEqual, "")
			})

			Convey("With env enabled", func() {
				opt := HTTPRequest(r).WithEnv()
				hr, ok := opt.(*httpRequestOption)
				So(ok, ShouldBeTrue)

				d := hr.buildData()
				So(d, ShouldNotBeNil)

				So(d.Method, ShouldEqual, "GET")
				So(d.URL, ShouldEqual, "https://example.com/test")
				So(d.Query, ShouldEqual, url.Values{
					"testing":  {"1"},
					"password": {"********"},
				}.Encode())

				So(d.Data, ShouldBeNil)
				So(d.Headers, ShouldResemble, map[string]string{})
				So(d.Env["REMOTE_ADDR"], ShouldEqual, "127.0.0.1")
				So(d.Env["REMOTE_PORT"], ShouldEqual, "12835")
				So(d.Cookies, ShouldEqual, "")
			})

			Convey("With data provided", func() {
				opt := HTTPRequest(r).WithData("testing")
				hr, ok := opt.(*httpRequestOption)
				So(ok, ShouldBeTrue)

				d := hr.buildData()
				So(d, ShouldNotBeNil)

				So(d.Method, ShouldEqual, "GET")
				So(d.URL, ShouldEqual, "https://example.com/test")
				So(d.Query, ShouldEqual, url.Values{
					"testing":  {"1"},
					"password": {"********"},
				}.Encode())

				So(d.Data, ShouldEqual, "testing")
				So(d.Headers, ShouldResemble, map[string]string{})
				So(d.Env, ShouldResemble, map[string]string{})
				So(d.Cookies, ShouldEqual, "")
			})
		})
	})
}
