package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleUnset() {
	cl := NewClient(
		// You can remove specific default fields from your final packet if you do
		// not wish to send them.
		Unset("runtime"),
	)

	cl.Capture(
		// You can also remove things that you may have added later
		Unset("message"),
	)
}

func TestUnset(t *testing.T) {
	Convey("Unset", t, func() {
		Convey("Unset()", func() {
			Convey("Should use the correct Class()", func() {
				So(Unset("runtime").Class(), ShouldEqual, "runtime")
				So(Unset("device").Class(), ShouldEqual, "device")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to null", func() {
					b, err := json.Marshal(Unset("runtime"))
					So(err, ShouldBeNil)
					So(string(b), ShouldEqual, `null`)
				})
			})
		})

		Convey("Should implement the advanced option interface", func() {
			So(Unset("runtime"), ShouldImplement, (*AdvancedOption)(nil))
		})

		Convey("Should correctly unset fields from the packet", func() {
			p := map[string]Option{}
			p["level"] = Level(Error)
			Unset("level").(AdvancedOption).Apply(p)

			So(p, ShouldResemble, map[string]Option{})
		})
	})
}
