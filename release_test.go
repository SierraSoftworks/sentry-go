package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRelease() {
	cl := NewClient(
		// You can set the release when you create a client
		Release("v1.0.0"),
	)

	cl.Capture(
		// You can also set it when you send an event
		Release("v1.0.0-dev"),
	)
}

func TestRelease(t *testing.T) {
	Convey("Release", t, func() {
		Convey("Release()", func() {
			Convey("Should use the correct Class()", func() {
				So(Release("test").Class(), ShouldEqual, "release")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					Convey("Should marshal to a string", func() {
						b, err := json.Marshal(Release("test"))
						So(err, ShouldBeNil)
						So(string(b), ShouldEqual, `"test"`)
					})
				})
			})
		})
	})
}
