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

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPTransport(t *testing.T) {
	deserializePacket := func(dataType string, data io.Reader) (interface{}, error) {
		var out interface{}

		if strings.Contains(dataType, "application/json") {
			if err := json.NewDecoder(data).Decode(&out); err != nil {
				return nil, err
			}
		} else if strings.Contains(dataType, "application/octet-stream") {
			b64 := base64.NewDecoder(base64.StdEncoding, data)
			deflate, err := zlib.NewReader(b64)
			defer deflate.Close()
			if err != nil {
				return nil, err
			}

			if err := json.NewDecoder(deflate).Decode(&out); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unknown datatype for packet: %s", dataType)
		}

		return out, nil
	}

	longMessage := func(minLength int) Option {
		msg := " "
		for len(msg) < 1000 {
			msg = fmt.Sprintf("%s%s", msg, msg)
		}

		return Message(msg)
	}

	Convey("HTTPTransport", t, func() {
		t := newHTTPTransport()
		So(t, ShouldNotBeNil)

		ht, ok := t.(*httpTransport)
		So(ok, ShouldBeTrue)

		Convey("newHTTPTransport", func() {
			So(ht.client, ShouldNotEqual, http.DefaultClient)
		})

		Convey("Send()", func(c C) {
			p := NewPacket()

			received := false

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
				received = true
				res.WriteHeader(200)
				res.Write([]byte("No data"))

				c.So(req.Method, ShouldEqual, "POST")
				c.So(req.RequestURI, ShouldEqual, "/api/1/store/")
				c.So(req.Header.Get("X-Sentry-Auth"), ShouldEqual, "Sentry sentry_version=4, sentry_key=user, sentry_secret=pass")

				expectedData, err := testSerializePacket(p)
				c.So(err, ShouldBeNil)

				data, err := deserializePacket(req.Header.Get("Content-Type"), req.Body)
				c.So(err, ShouldBeNil)
				c.So(data, ShouldNotBeNil)
				c.So(data, ShouldResemble, expectedData)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			uri, err := url.Parse(ts.URL)
			So(err, ShouldBeNil)
			uri.User = url.UserPassword("user", "pass")
			uri.Path = "/1"

			dsn := uri.String()

			Convey("With a small packet", func() {
				So(t.Send(dsn, p), ShouldBeNil)
				So(received, ShouldBeTrue)
			})

			Convey("With a large packet", func() {
				p.SetOptions(longMessage(1000))
				So(t.Send(dsn, p), ShouldBeNil)
				So(received, ShouldBeTrue)
			})
		})

		Convey("serializePacket()", func() {
			p := NewPacket()

			Convey("Short Packet", func() {
				data, dataType, err := ht.serializePacket(p)
				So(err, ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(dataType, ShouldContainSubstring, "application/json")

				pd, err := deserializePacket(dataType, data)
				So(err, ShouldBeNil)

				ped, err := testSerializePacket(p)
				So(err, ShouldBeNil)
				So(pd, ShouldResemble, ped)
			})

			Convey("Long Packet", func() {
				p.SetOptions(longMessage(10000))
				data, dataType, err := ht.serializePacket(p)
				So(err, ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(dataType, ShouldContainSubstring, "application/octet-stream")

				pd, err := deserializePacket(dataType, data)
				So(err, ShouldBeNil)

				ped, err := testSerializePacket(p)
				So(err, ShouldBeNil)
				So(pd, ShouldResemble, ped)
			})
		})

		Convey("parseDSN()", func() {
			Convey("With an empty DSN", func() {
				url, authheader, err := ht.parseDSN("")
				So(err, ShouldBeNil)
				So(url, ShouldEqual, "")
				So(authheader, ShouldEqual, "")
			})

			Convey("With an invalid DSN", func() {
				url, authheader, err := ht.parseDSN("@")
				So(err, ShouldNotBeNil)
				So(url, ShouldEqual, "")
				So(authheader, ShouldEqual, "")
			})

			Convey("With a valid DSN", func() {
				url, authHeader, err := ht.parseDSN("https://user:pass@example.com/sentry/1")
				So(err, ShouldBeNil)
				So(url, ShouldEqual, "https://example.com/sentry/api/1/store/")
				So(authHeader, ShouldEqual, "Sentry sentry_version=4, sentry_key=user, sentry_secret=pass")
			})
		})

		// If you set $SENTRY_DSN you can send events to a live Sentry instance
		// to confirm that this library functions correctly.
		if liveTestDSN := os.Getenv("SENTRY_DSN"); liveTestDSN != "" {
			Convey("Live Test", func() {
				log.SetLevel(log.DebugLevel)
				defer log.SetLevel(log.InfoLevel)

				p := NewPacket().SetOptions(
					Message("Ran Live Test"),
					Release(version),
					Level(Debug),
				)

				So(t.Send(liveTestDSN, p), ShouldBeNil)
			})
		}
	})
}
