package sentry

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleEnvironment() {
	cl := NewClient(
		// You can configure your environment at the client level
		Environment("development"),
	)

	cl.Capture(
		// ...or at the event level
		Environment("prod"),
	)
}

func TestEnvironment(t *testing.T) {
	Convey("Environment", t, func() {
		Convey("Should register with the default providers", func() {
			Convey("If the ENV environment variable is set", func() {
				os.Setenv("ENV", "testing")
				defer os.Unsetenv("ENV")

				opt := testGetOptionsProvider(&environmentOption{})
				So(opt, ShouldNotBeNil)
				So(opt.(*environmentOption).env, ShouldEqual, "testing")
			})

			Convey("If the ENVIRONMENT environment variable is set", func() {
				os.Setenv("ENVIRONMENT", "testing")
				defer os.Unsetenv("ENVIRONMENT")

				opt := testGetOptionsProvider(&environmentOption{})
				So(opt, ShouldNotBeNil)
				So(opt.(*environmentOption).env, ShouldEqual, "testing")
			})
		})

		Convey("Environment()", func() {
			Convey("Should return an Option", func() {
				So(Environment("testing"), ShouldImplement, (*Option)(nil))
			})
		})

		Convey("Should use the correct class", func() {
			So(Environment("test").Class(), ShouldEqual, "environment")
		})

		Convey("MarshalJSON", func() {
			So(testOptionsSerialize(Environment("test")), ShouldResemble, "test")
		})
	})
}
