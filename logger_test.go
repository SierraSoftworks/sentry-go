package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleLogger() {
	cl := NewClient(
		// You can set the logger when you create your client
		Logger("root"),
	)

	cl.Capture(
		// You can also specify it when sending an event
		Logger("http"),
	)
}

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
