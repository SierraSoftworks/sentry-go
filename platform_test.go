package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExamplePlatform() {
	cl := NewClient(
		// You can set the platform at a client level
		Platform("go"),
	)

	cl.Capture(
		// Or override it when sending the event
		Platform("go"),
	)
}

func TestPlatform(t *testing.T) {
	Convey("Platform", t, func() {
		Convey("Should register itself with the default providers", func() {
			opt := testGetOptionsProvider(Platform("go"))
			So(opt, ShouldNotBeNil)
		})

		Convey("Platform()", func() {
			Convey("Should use the correct Class()", func() {
				So(Platform("go").Class(), ShouldEqual, "platform")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					Convey("Should marshal to a string", func() {
						b, err := json.Marshal(Platform("go"))
						So(err, ShouldBeNil)
						So(string(b), ShouldEqual, `"go"`)
					})
				})
			})
		})
	})
}
