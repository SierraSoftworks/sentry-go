package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleConfig() {
	cl := NewClient(
		// Specify the DSN to use for sending events, or ""
		// to disable sending events altogether.
		DSN(""),

		// Specify a custom transport to use for sending the
		// events to Sentry. Nil resets this to its default (HTTP)
		UseTransport(nil),
	)

	cl.Capture(Message("Example"))
}

func TestConfig(t *testing.T) {
	Convey("Config", t, func() {
		Convey("Should register itself as a default option provider", func() {
			provider := testGetOptionsProvider(&configOption{})
			So(provider, ShouldNotBeNil)
		})

		opt := &configOption{}

		Convey("Should implement the Option interface", func() {
			So(opt, ShouldImplement, (*Option)(nil))
		})

		Convey("Should implement the OmitableOption interface", func() {
			So(opt, ShouldImplement, (*OmitableOption)(nil))
		})

		Convey("Should implement the MergableOption interface", func() {
			So(opt, ShouldImplement, (*OmitableOption)(nil))
		})

		Convey("Omit()", func() {
			Convey("Should always return true", func() {
				So(opt.Omit(), ShouldBeTrue)
			})
		})

		Convey("Merge()", func() {
			Convey("Should not modify the original config objects", func() {
				opt2 := DSN("test2")

				c := opt.Merge(opt2)
				So(c, ShouldNotEqual, opt)
				So(c, ShouldNotEqual, opt2)
				So(c, ShouldHaveSameTypeAs, opt2)
			})

			Convey("Should overwrite the DSN option if a DSN is provided", func() {
				opt := DSN("newDSN").(*configOption)
				c := opt.Merge(DSN("oldDSN"))
				So(c, ShouldHaveSameTypeAs, opt)
				So(*c.(*configOption).dsn, ShouldEqual, "newDSN")
			})

			Convey("Should not overwrite the DSN option if it is not provided", func() {
				opt := &configOption{}
				c := opt.Merge(DSN("oldDSN"))
				So(c, ShouldHaveSameTypeAs, opt)
				So(*c.(*configOption).dsn, ShouldEqual, "oldDSN")
			})

			Convey("Should overwrite if it doesn't recognize the old option's type", func() {
				c := opt.Merge(&testOption{})
				So(c, ShouldEqual, opt)
			})
		})

		Convey("Clone()", func() {
			Convey("Should return a new configOption", func() {
				c := opt.Clone()
				So(c, ShouldNotBeNil)
				So(c, ShouldNotEqual, opt)
			})
		})

		Convey("Transport()", func() {
			Convey("When transport is nil it should return the DefaultTransport", func() {
				c := &configOption{}
				So(c.Transport(), ShouldEqual, DefaultTransport())
			})

			Convey("When transport is defined, it should return that", func() {
				t := newHTTPTransport()
				c := &configOption{
					transport: t,
				}
				So(c.Transport(), ShouldEqual, t)
			})
		})

		Convey("DSN", func() {
			Convey("When DSN is nil it should return an empty DSN", func() {
				c := &configOption{}
				So(c.DSN(), ShouldEqual, "")
			})

			Convey("When DSN is defined, it should return that", func() {
				d := "https://key:secret@example.com/sentry/1"
				c := &configOption{
					dsn: &d,
				}
				So(c.DSN(), ShouldEqual, "https://key:secret@example.com/sentry/1")
			})
		})
	})
}
