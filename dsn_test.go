package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleDSN() {
	cl := NewClient(
		// You can configure the DSN when creating a client
		DSN("https://key:pass@example.com/sentry/1"),
	)

	cl.Capture(
		// You can also configure the DSN when sending an event
		DSN(""),
		Message("This won't be sent"),
	)
}

func TestDSN(t *testing.T) {
	Convey("DSN", t, func() {
		Convey("Parse()", func() {
			cases := []struct {
				Name  string
				URL   string
				Error error
			}{
				{"With a valid URL", "https://u:p@example.com/sentry/1", nil},
				{"With a badly formatted URL", ":", ErrBadURL},
				{"Without a public key", "https://example.com/sentry/1", ErrMissingPublicKey},
				{"Without a private key", "https://u@example.com/sentry/1", ErrMissingPrivateKey},
				{"Without a project ID", "https://u:p@example.com", ErrMissingProjectID},
			}

			for _, tc := range cases {
				Convey(tc.Name, func() {
					d := &dsn{}
					err := d.Parse(tc.URL)
					if tc.Error == nil {
						So(err, ShouldBeNil)
					} else {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldStartWith, tc.Error.Error())
					}
				})
			}
		})
	})
}
