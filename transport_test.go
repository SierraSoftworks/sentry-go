package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransport(t *testing.T) {
	Convey("Transport", t, func() {
		Convey("UseTransport()", func() {
			Convey("Should return an Option", func() {
				So(UseTransport(newHTTPTransport()), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if the transport is nil", func() {
				So(UseTransport(nil), ShouldBeNil)
			})
		})

		Convey("DefaultTransport()", func() {
			Convey("Should be defined from the start", func() {
				So(DefaultTransport(), ShouldNotBeNil)
			})

			Convey("Should return the correct transport", func() {
				So(DefaultTransport(), ShouldEqual, defaultTransport)
			})
		})

		Convey("SetDefaultTransport()", func() {
			Convey("Should allow you to change the default transport", func() {
				t := newHTTPTransport()
				SetDefaultTransport(t)
				So(DefaultTransport(), ShouldEqual, t)
			})
		})
	})
}
