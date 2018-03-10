package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleMessage() {
	cl := NewClient()

	cl.Capture(
		// You can either use just a simple message
		Message("this is a simple message"),
	)

	cl.Capture(
		// Or you can provide formatting entries as you would with
		// fmt.Sprintf() calls.
		Message("this is a %s message (%d/7 would use again)", "formatted", 5),
	)
}

func TestMessage(t *testing.T) {
	Convey("Message", t, func() {
		Convey("Message()", func() {
			Convey("Should return an Option", func() {
				So(Message("test"), ShouldImplement, (*Option)(nil))
			})

			Convey("With just a message string", func() {
				m := Message("test")
				So(m, ShouldNotBeNil)

				mi, ok := m.(*messageOption)
				So(ok, ShouldBeTrue)

				So(mi.Message, ShouldEqual, "test")
			})

			Convey("With a formatted message", func() {
				m := Message("this is a %s", "test")
				So(m, ShouldNotBeNil)

				mi, ok := m.(*messageOption)
				So(ok, ShouldBeTrue)

				So(mi.Message, ShouldEqual, "this is a %s")
				So(mi.Params, ShouldResemble, []interface{}{"test"})
				So(mi.Formatted, ShouldEqual, "this is a test")
			})
		})

		Convey("Should use the correct Class()", func() {
			So(Message("test").Class(), ShouldEqual, "sentry.interfaces.Message")
		})
	})
}
