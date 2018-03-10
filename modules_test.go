package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleModules() {
	cl := NewClient(
		// You can specify module versions when creating your
		// client
		Modules(map[string]string{
			"redis": "v1",
			"mgo": "v2",
		}),
	)

	cl.Capture(
		// And override or expand on them when sending an event
		Modules(map[string]string{
			"redis": "v2",
			"sentry-go": "v1",
		}),
	)
}

func TestModules(t *testing.T) {
	Convey("Modules", t, func() {
		Convey("Modules()", func() {
			data := map[string]string{
				"redis": "1.0.0",
			}

			Convey("Should return nil if the data is nil", func() {
				So(Modules(nil), ShouldBeNil)
			})

			Convey("Should return an Option", func() {
				So(Modules(data), ShouldImplement, (*Option)(nil))
			})

			Convey("Should use the correct Class()", func() {
				So(Modules(data).Class(), ShouldEqual, "modules")
			})

			Convey("Should implement Merge()", func() {
				So(Modules(data), ShouldImplement, (*MergeableOption)(nil))
			})
		})

		Convey("Merge()", func() {
			data1 := map[string]string{
				"redis": "1.0.0",
			}
			e1 := Modules(data1)
			So(e1, ShouldNotBeNil)

			e1m, ok := e1.(MergeableOption)
			So(ok, ShouldBeTrue)

			Convey("Should overwrite if it doesn't recognize the old option", func() {
				So(e1m.Merge(&testOption{}), ShouldEqual, e1)
			})

			Convey("Should merge multiple modules entries", func() {
				data2 := map[string]string{
					"pgsql": "5.4.0",
				}

				e2 := Modules(data2)
				So(e2, ShouldNotBeNil)

				em := e1m.Merge(e2)
				So(em, ShouldNotBeNil)
				So(em, ShouldNotEqual, e1)
				So(em, ShouldNotEqual, e2)

				emm, ok := em.(*modulesOption)
				So(ok, ShouldBeTrue)
				So(emm.moduleVersions, ShouldContainKey, "redis")
				So(emm.moduleVersions, ShouldContainKey, "pgsql")
			})

			Convey("Should overwrite old entries with new ones", func() {
				data2 := map[string]string{
					"redis": "0.8.0",
				}
				e2 := Modules(data2)
				So(e2, ShouldNotBeNil)

				em := e1m.Merge(e2)
				So(em, ShouldNotBeNil)
				So(em, ShouldNotEqual, e1)
				So(em, ShouldNotEqual, e2)

				emm, ok := em.(*modulesOption)
				So(ok, ShouldBeTrue)
				So(emm.moduleVersions, ShouldContainKey, "redis")
				So(emm.moduleVersions["redis"], ShouldEqual, "1.0.0")
			})
		})

		Convey("MarshalJSON", func() {
			Convey("Should marshal the fields correctly", func() {
				data := map[string]string{
					"redis": "1.0.0",
				}

				serialized := testOptionsSerialize(Modules(data))
				So(serialized, ShouldNotBeNil)

				expected := map[string]interface{}{
					"redis": "1.0.0",
				}
				So(serialized, ShouldHaveSameTypeAs, expected)
				So(serialized, ShouldResemble, expected)
			})
		})
	})
}
