package sentry

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
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
	r, err := http.NewRequest("GET", "https://example.com/test?testing=1&password=test", nil)
	assert.Nil(t, err, "should be able to create an HTTP request object")

	r.RemoteAddr = "127.0.0.1:12835"
	r.Header.Set("Host", "example.com")
	r.Header.Set("X-Forwarded-Proto", "https")
	r.Header.Set("Cookie", "testing=1")
	r.Header.Set("X-Testing", "1")
	r.Header.Set("X-API-Key", "secret")

	assert.NotNil(t, HTTPRequest(nil), "it should not return nil if no request is provided")

	o := HTTPRequest(r)
	assert.NotNil(t, o, "should not return a nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "request", o.Class(), "it should use the right option class")

	if assert.Implements(t, (*OmitableOption)(nil), o, "it should implement the OmitableOption interface") {
		assert.False(t, o.(OmitableOption).Omit(), "it should return false if there is a request")
		assert.True(t, HTTPRequest(nil).(OmitableOption).Omit(), "it should return true if there is no request")
	}

	tm := "GET"
	tu := "https://example.com/test"
	tq := url.Values{
		"testing":  {"1"},
		"password": {sanitizationString},
	}
	var td interface{} = nil
	th := map[string]string{}
	te := map[string]string{}
	tc := ""

	cases := []struct {
		Name  string
		Opt   Option
		Setup func()
	}{
		{"Default", HTTPRequest(r), func() {}},
		{"Default.Sanitize()", HTTPRequest(r).Sanitize("testing"), func() {
			tq = url.Values{
				"testing":  {sanitizationString},
				"password": {sanitizationString},
			}
		}},
		{"WithCookies()", HTTPRequest(r).WithCookies(), func() {
			tc = "testing=1"
		}},
		{"WithHeaders()", HTTPRequest(r).WithHeaders(), func() {
			th = map[string]string{
				"Host":              "example.com",
				"Cookie":            "testing=1",
				"X-Testing":         "1",
				"X-Forwarded-Proto": "https",
				"X-Api-Key":         "secret",
			}
		}},
		{"WithHeaders().Sanitize()", HTTPRequest(r).WithHeaders().Sanitize("key"), func() {
			th = map[string]string{
				"Host":              "example.com",
				"Cookie":            "testing=1",
				"X-Testing":         "1",
				"X-Forwarded-Proto": "https",
				"X-Api-Key":         sanitizationString,
			}
		}},
		{"WithEnv()", HTTPRequest(r).WithEnv(), func() {
			te = map[string]string{
				"REMOTE_ADDR": "127.0.0.1",
				"REMOTE_PORT": "12835",
			}
		}},
		{"WithData()", HTTPRequest(r).WithData("testing"), func() {
			td = "testing"
		}},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			tq = url.Values{
				"testing":  {"1"},
				"password": {sanitizationString},
			}
			td = nil
			th = map[string]string{}
			te = map[string]string{}
			tc = ""

			testCase.Setup()

			hr, ok := testCase.Opt.(*httpRequestOption)
			assert.True(t, ok, "the option should actually be a *httpRequestOption")

			d := hr.buildData()
			assert.NotNil(t, d, "the built data should not be nil")

			assert.Equal(t, tm, d.Method, "the method should be correct")
			assert.Equal(t, tu, d.URL, "the url should be correct")
			assert.Equal(t, tq.Encode(), d.Query, "the query should be correct")
			assert.Equal(t, td, d.Data, "the data should be correct")
			assert.Equal(t, th, d.Headers, "the headers should be correct")

			for k, v := range te {
				if assert.Contains(t, d.Env, k, "the environment should include the %s entry", k) {
					assert.Equal(t, v, d.Env[k], "the value of the %s environment variable should be correct", k)
				}
			}

			assert.Equal(t, tc, d.Cookies, "the cookies should be correct")
		})
	}

	t.Run("MarshalJSON()", func(t *testing.T) {

	})
}
