package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	Convey("Logger", t, func() {
		Convey("Should register itself with the default providers", func() {
			opt := testGetOptionsProvider(Logger(""))
			So(opt, ShouldNotBeNil)
		})

		Convey("Logger()", func() {
			Convey("Should use the correct Class()", func() {
				So(Logger("test").Class(), ShouldEqual, "logger")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					Convey("Should marshal to a string", func() {
						b, err := json.Marshal(Logger("test"))
						So(err, ShouldBeNil)
						So(string(b), ShouldEqual, `"test"`)
					})
				})
			})
		})
	})
}
