package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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

		Convey("Should implement the OmmitableOption interface", func() {
			So(opt, ShouldImplement, (*OmmitableOption)(nil))
		})

		Convey("Should implement the MergableOption interface", func() {
			So(opt, ShouldImplement, (*OmmitableOption)(nil))
		})

		Convey("Ommit()", func() {
			Convey("Should always return true", func() {
				So(opt.Ommit(), ShouldBeTrue)
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
	})
}
