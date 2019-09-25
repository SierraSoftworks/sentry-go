package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	t.Run("Parse()", func (t *testing.T) {
		cases := []struct {
			Name  string
			URL   string
			Error error
		}{
			{"With a valid URL", "https://u:p@example.com/sentry/1", nil},
			{"With a badly formatted URL", ":", ErrBadURL},
			{"Without a public key", "https://example.com/sentry/1", ErrMissingPublicKey},
			{"Without a private key", "https://u@example.com/sentry/1", nil},
			{"Without a project ID", "https://u:p@example.com", ErrMissingProjectID},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				d := &dsn{}
				err := d.Parse(tc.URL)

				if tc.Error == nil {
					assert.Nil(t, err, "it should not return an error")
				} else {
					assert.NotNil(t, err, "it should return an error")
					assert.Regexp(t, "^" + tc.Error.Error(), err.Error(), "it should return the right error")
				}
			})
		}
	})

	t.Run("AuthHeader()", func(t *testing.T) {
		assert.Equal(t, "", (&dsn{PrivateKey: "secret"}).AuthHeader(), "should return no auth header if no public key is provided")
		assert.Equal(t, "Sentry sentry_version=4, sentry_key=key", (&dsn{PublicKey: "key"}).AuthHeader(), "should return an auth header with just the public key if no private key is provided")
		assert.Equal(t, "Sentry sentry_version=4, sentry_key=key, sentry_secret=secret", (&dsn{PublicKey: "key", PrivateKey: "secret"}).AuthHeader(), "should return a full auth header both the public and private key are provided")
	})
}
