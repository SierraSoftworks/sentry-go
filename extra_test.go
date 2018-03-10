package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleExtra() {
	cl := NewClient(
		// You can define extra fields when you create your client
		Extra(map[string]interface{}{
			"redis": map[string]interface{}{
				"host": "redis",
				"port": 6379,
			},
		}),
	)

	cl.Capture(
		// You can also define extra info when you send the event
		// The extra object will be shallowly merged automatically,
		// so this would send both `redis` and `cache`.
		Extra(map[string]interface{}{
			"cache": map[string]interface{}{
				"key": "user.127.profile",
				"hit": false,
			},
		}),
	)
}

func TestExtra(t *testing.T) {
	Convey("Extra", t, func() {
		Convey("Extra()", func() {
			data := map[string]interface{}{
				"redis": map[string]interface{}{
					"host": "redis",
					"port": 6379,
				},
			}

			Convey("Should return nil if the data is nil", func() {
				So(Extra(nil), ShouldBeNil)
			})

			Convey("Should return an Option", func() {
				So(Extra(data), ShouldImplement, (*Option)(nil))
			})

			Convey("Should use the correct Class()", func() {
				So(Extra(data).Class(), ShouldEqual, "extra")
			})

			Convey("Should implement Merge()", func() {
				So(Extra(data), ShouldImplement, (*MergeableOption)(nil))
			})
		})

		Convey("Merge()", func() {
			data1 := map[string]interface{}{
				"redis": map[string]interface{}{
					"host": "redis",
					"port": 6379,
				},
			}
			e1 := Extra(data1)
			So(e1, ShouldNotBeNil)

			e1m, ok := e1.(MergeableOption)
			So(ok, ShouldBeTrue)

			Convey("Should overwrite if it doesn't recognize the old option", func() {
				So(e1m.Merge(&testOption{}), ShouldEqual, e1)
			})

			Convey("Should merge multiple extra entries", func() {
				data2 := map[string]interface{}{
					"cache": map[string]interface{}{
						"key": "user.127.profile",
						"hit": false,
					},
				}
				e2 := Extra(data2)
				So(e2, ShouldNotBeNil)

				em := e1m.Merge(e2)
				So(em, ShouldNotBeNil)
				So(em, ShouldNotEqual, e1)
				So(em, ShouldNotEqual, e2)

				emm, ok := em.(*extraOption)
				So(ok, ShouldBeTrue)
				So(emm.extra, ShouldContainKey, "cache")
				So(emm.extra, ShouldContainKey, "redis")
			})

			Convey("Should overwrite old entries with new ones", func() {
				data2 := map[string]interface{}{
					"redis": map[string]interface{}{
						"host": "redis-dev",
						"port": 6379,
					},
				}
				e2 := Extra(data2)
				So(e2, ShouldNotBeNil)

				em := e1m.Merge(e2)
				So(em, ShouldNotBeNil)
				So(em, ShouldNotEqual, e1)
				So(em, ShouldNotEqual, e2)

				emm, ok := em.(*extraOption)
				So(ok, ShouldBeTrue)
				So(emm.extra, ShouldContainKey, "redis")
				So(emm.extra["redis"], ShouldResemble, map[string]interface{}{
					"host": "redis",
					"port": 6379,
				})
			})
		})

		Convey("MarshalJSON", func() {
			Convey("Should marshal the fields correctly", func() {
				data := map[string]interface{}{
					"redis": map[string]interface{}{
						"host": "redis",
						// Float mode required since we aren't deserializing into an int
						"port": 6379.,
					},
				}

				serialized := testOptionsSerialize(Extra(data))
				So(serialized, ShouldNotBeNil)
				So(serialized, ShouldHaveSameTypeAs, data)
				So(serialized, ShouldResemble, data)
			})
		})
	})
}
