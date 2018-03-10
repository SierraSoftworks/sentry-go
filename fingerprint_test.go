package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleFingerprint() {
	cl := NewClient()

	cl.Capture(
		// You can specify a fingerprint that extends the default behaviour
		Fingerprint("{{ default }}", "http://example.com/my.url"),

		// Or you can define your own
		Fingerprint("myrpc", "POST", "/foo.bar"),
	)
}

func TestFingerprint(t *testing.T) {
	Convey("Fingerprint", t, func() {
		Convey("Fingerprint()", func() {
			Convey("Should return an Option", func() {
				So(Fingerprint("test"), ShouldImplement, (*Option)(nil))
			})
		})

		Convey("Should use the correct class", func() {
			So(Fingerprint("test").Class(), ShouldEqual, "fingerprint")
		})

		Convey("MarshalJSON", func() {
			So(testOptionsSerialize(Fingerprint("test")), ShouldResemble, []interface{}{"test"})
		})
	})
}
