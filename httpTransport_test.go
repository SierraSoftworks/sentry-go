package sentry

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPTransport(t *testing.T) {
	deserializePacket := func(t *testing.T, dataType string, data io.Reader) interface{} {
		var out interface{}

		if strings.Contains(dataType, "application/json") {
			require.Nil(t, json.NewDecoder(data).Decode(&out), "there should be no problems deserializing the packet")
		} else if strings.Contains(dataType, "application/octet-stream") {
			b64 := base64.NewDecoder(base64.StdEncoding, data)
			deflate, err := zlib.NewReader(b64)
			defer deflate.Close()

			require.Nil(t, err, "there should be no errors creating the zlib deflator")
			require.Nil(t, json.NewDecoder(deflate).Decode(&out), "there should be no problems deserializing the packet")
		} else {
			t.Fatalf("unknown datatype for packet: %s", dataType)
		}

		return out
	}

	longMessage := func(minLength int) Option {
		msg := " "
		for len(msg) < 1000 {
			msg = fmt.Sprintf("%s%s", msg, msg)
		}

		return Message(msg)
	}

	tr := newHTTPTransport()
	require.NotNil(t, tr, "the transport should not be nil")

	ht, ok := tr.(*httpTransport)
	require.True(t, ok, "it should actually be a *httpTransport")

	t.Run("Send()", func(t *testing.T) {
		p := NewPacket()
		require.NotNil(t, p, "the packet should not be nil")

		received := false
		statusCode := 200

		mux := http.NewServeMux()
		require.NotNil(t, mux, "the http mux should not be nil")
		mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
			received = true
			res.WriteHeader(statusCode)
			res.Write([]byte("No data"))

			assert.Equal(t, "POST", req.Method, "the request should use HTTP POST")
			assert.Equal(t, "/api/1/store/", req.RequestURI, "the request should use the right API endpoint")

			assert.Contains(t, []string{
				"Sentry sentry_version=4, sentry_key=key, sentry_secret=secret",
				"Sentry sentry_version=4, sentry_key=key",
			}, req.Header.Get("X-Sentry-Auth"), "it should use the right auth header")

			expectedData := testSerializePacket(t, p)

			data := deserializePacket(t, req.Header.Get("Content-Type"), req.Body)

			assert.Equal(t, expectedData, data, "the data should match what we expected")
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		makeDSN := func(publicKey, privateKey string) string {
			uri, err := url.Parse(ts.URL)
			require.Nil(t, err, "we should not fail to parse the URI")

			if publicKey != "" {
				uri.User = url.UserPassword(publicKey, privateKey)
			}

			uri.Path = "/1"

			return uri.String()
		}

		cases := []struct {
			Name       string
			Packet     Packet
			DSN        string
			StatusCode int
			Error      error
			Received   bool
		}{
			{"Short Packet", NewPacket(), makeDSN("key", "secret"), 200, nil, true},
			{"Long Packet", NewPacket().SetOptions(longMessage(10000)), makeDSN("key", "secret"), 200, nil, true},
			{"No DSN", NewPacket(), "", 200, nil, false},
			{"Invalid DSN URL", NewPacket(), ":", 400, ErrBadURL, false},
			{"Missing Public Key", NewPacket(), makeDSN("", ""), 401, ErrMissingPublicKey, false},
			{"Invalid Server", NewPacket(), "https://key:secret@invalid_domain.not_a_tld/sentry/1", 404, ErrType("failed to submit request"), false},
			{"Missing Private Key with Required Key", NewPacket(), makeDSN("key", ""), 401, fmt.Errorf("got http status 401, expected 200"), true},
			{"Missing Private Key", NewPacket(), makeDSN("key", ""), 200, nil, true},
		}

		for _, tc := range cases {
			tc := tc

			t.Run(tc.Name, func(t *testing.T) {
				received = false
				statusCode = tc.StatusCode
				p = tc.Packet

				err := tr.Send(tc.DSN, tc.Packet)
				if tc.Error == nil {
					assert.Nil(t, err, "it should not fail to send the packet")
				} else if errType, ok := tc.Error.(ErrType); ok {
					assert.True(t, errType.IsInstance(err), "it should return the right error")
				} else {
					assert.EqualError(t, err, tc.Error.Error(), "it should return the right error")
				}

				if tc.Received {
					assert.True(t, received, "the server should have received the packet")
				} else {
					assert.False(t, received, "the server should not have received the packet")
				}
			})
		}
	})

	t.Run("serializePacket()", func(t *testing.T) {
		cases := []struct {
			Name     string
			Packet   Packet
			DataType string
		}{
			{"Short Packet", NewPacket().SetOptions(Message("short packet")), "application/json; charset=utf8"},
			{"Long Packet", NewPacket().SetOptions(longMessage(10000)), "application/octet-stream"},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				data, dataType, err := ht.serializePacket(tc.Packet)
				assert.Nil(t, err, "there should be no error serializing the packet")
				assert.Equal(t, tc.DataType, dataType, "the request datatype should be %s", tc.DataType)
				assert.NotNil(t, data, "the request data should not be nil")

				assert.Equal(t, testSerializePacket(t, tc.Packet), deserializePacket(t, dataType, data), "the serialized packet should match what we expected")
			})
		}
	})

	t.Run("parseDSN()", func(t *testing.T) {
		cases := []struct {
			Name       string
			DSN        string
			URL        string
			AuthHeader string
			Error      error
		}{
			{"Empty DSN", "", "", "", nil},
			{"Invalid DSN", "@", "", "", fmt.Errorf("sentry: missing public key: missing URL user")},
			{"Full DSN", "https://user:pass@example.com/sentry/1", "https://example.com/sentry/api/1/store/", "Sentry sentry_version=4, sentry_key=user, sentry_secret=pass", nil},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				url, authHeader, err := ht.parseDSN(tc.DSN)
				if tc.Error != nil {
					assert.EqualError(t, err, tc.Error.Error(), "there should be an error with the right message")
				} else {
					assert.Nil(t, err, "there should be no error")
				}
				assert.Equal(t, tc.URL, url, "the parsed URL should be correct")
				assert.Equal(t, tc.AuthHeader, authHeader, "the parsed auth header should be correct")
			})
		}
	})

	// If you set $SENTRY_DSN you can send events to a live Sentry instance
	// to confirm that this library functions correctly.
	if liveTestDSN := os.Getenv("SENTRY_DSN"); liveTestDSN != "" {
		t.Run("Live Test", func(t *testing.T) {
			p := NewPacket().SetOptions(
				Message("Ran Live Test"),
				Release(version),
				Level(Debug),
			)

			assert.Nil(t, tr.Send(liveTestDSN, p), "it should not fail to send the packet")
		})
	}
}
