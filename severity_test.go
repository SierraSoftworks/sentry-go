package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleLevel() {
	cl := NewClient(
		// You can set the severity level when you create your client
		Level(Debug),
	)

	cl.Capture(
		// You can also specify it when sending an event
		Level(Error),
	)
}

func TestSeverity(t *testing.T) {
	Convey("Severity", t, func() {
		Convey("Level()", func() {
			Convey("Should use the correct Class()", func() {
				So(Level(Error).Class(), ShouldEqual, "level")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					b, err := json.Marshal(Level(Error))
					So(err, ShouldBeNil)
					So(string(b), ShouldEqual, `"error"`)
				})
			})
		})

		Convey("Fatal should use the correct name", func() {
			So(string(Fatal), ShouldEqual, "fatal")
		})

		Convey("Error should use the correct name", func() {
			So(string(Error), ShouldEqual, "error")
		})

		Convey("Warning should use the correct name", func() {
			So(string(Warning), ShouldEqual, "warning")
		})

		Convey("Info should use the correct name", func() {
			So(string(Info), ShouldEqual, "info")
		})

		Convey("Debug should use the correct name", func() {
			So(string(Debug), ShouldEqual, "debug")
		})
	})
}
