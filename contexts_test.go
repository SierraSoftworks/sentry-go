package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContexts(t *testing.T) {
	Convey("Contexts", t, func() {
		Convey("RuntimeContext()", func() {

		})

		Convey("OSContext()", func() {

		})

		Convey("DeviceContext()", func() {

		})

		Convey("Context()", func() {
			c := Context("test", "data")
			So(c, ShouldNotBeNil)
			So(c, ShouldHaveSameTypeAs, &contextOption{})

			cc := c.(*contextOption)
			So(cc.contexts, ShouldContainKey, "test")
			So(cc.contexts["test"], ShouldEqual, "data")
		})

		Convey("ContextOption", func() {
			c := Context("test", "data")
			So(c, ShouldNotBeNil)

			Convey("Should have the correct Class()", func() {
				So(c.Class(), ShouldEqual, "contexts")
			})

			Convey("Should implement MergableOption interface", func() {
				So(c, ShouldImplement, (*MergeableOption)(nil))
			})

			Convey("Merge()", func() {
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)

				Convey("Should overwrite if it cannot identify the old type", func() {
					out := cc.Merge(&testOption{})
					So(out, ShouldEqual, c)
				})

				Convey("Should overwriting old fields", func() {
					old := Context("test", "oldData")
					out := cc.Merge(old)
					So(out, ShouldNotBeNil)
					So(out, ShouldHaveSameTypeAs, &contextOption{})

					So(out.(*contextOption).contexts, ShouldResemble, map[string]interface{}{
						"test": "data",
					})
				})

				Convey("Should add new fields", func() {
					old := Context("old", "data")
					out := cc.Merge(old)
					So(out, ShouldNotBeNil)
					So(out, ShouldHaveSameTypeAs, &contextOption{})

					So(out.(*contextOption).contexts, ShouldResemble, map[string]interface{}{
						"old":  "data",
						"test": "data",
					})
				})
			})

			Convey("MarshalJSON", func() {
				c := Context("test", "data")
				So(testOptionsSerialize(c), ShouldResemble, map[string]interface{}{
					"test": "data",
				})
			})
		})
	})
}
