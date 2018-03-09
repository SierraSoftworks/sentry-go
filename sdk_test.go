package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSDKOption(t *testing.T) {
	Convey("SDK Option", t, func() {
		Convey("Should register itself with the default providers", func() {
			opt := testGetOptionsProvider(&sdkOption{})
			So(opt, ShouldNotBeNil)
		})

		Convey("Should register with the correct name", func() {
			opt := testGetOptionsProvider(&sdkOption{})
			So(opt, ShouldNotBeNil)

			oo := opt.(*sdkOption)
			So(oo.Name, ShouldEqual, "SierraSoftworks/sentry-go")
		})

		Convey("Should register with the correct version", func() {
			opt := testGetOptionsProvider(&sdkOption{})
			So(opt, ShouldNotBeNil)

			oo := opt.(*sdkOption)
			So(oo.Version, ShouldEqual, version)
		})
	})
}
