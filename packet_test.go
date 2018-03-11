package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExamplePacket() {
	// Create a new packet object which can be sent to
	// Sentry by one of the transports or send queues.
	p := NewPacket().SetOptions(
		DSN(""),
		Message("Custom packet creation"),
	)

	// Create a clone of this packet if you want to use
	// it as a template
	p.Clone().SetOptions(
		Message("Overridden message which doesn't affect the original"),
	)
}

func TestPacket(t *testing.T) {
	Convey("Packet", t, func() {
		Convey("NewPacket()", func() {
			p := NewPacket()
			So(p, ShouldNotBeNil)
			So(p, ShouldImplement, (*Packet)(nil))
		})

		Convey("Clone()", func() {
			p := NewPacket()
			p2 := p.Clone()

			So(p, ShouldNotEqual, p2)
			So(p, ShouldResemble, p2)
		})

		Convey("SetOptions()", func() {
			p := NewPacket()
			So(p.SetOptions(), ShouldResemble, p)

			pp, ok := p.(*packet)
			So(ok, ShouldBeTrue)
			So(pp, ShouldNotBeNil)

			pi := *pp

			Convey("Should ignore nil options", func() {
				Convey("When only nil options are provided", func() {
					p2 := p.Clone()
					So(p.SetOptions(nil), ShouldResemble, p2)
				})

				Convey("When both nil and non-nil options are provided", func() {
					p2 := p.Clone()
					opt := &testOption{}
					So(p.SetOptions(nil, opt), ShouldResemble, p2.SetOptions(opt))
				})
			})

			Convey("Should set normal option fields", func() {
				opt := &testOption{}
				p.SetOptions(opt)
				So(pi, ShouldContainKey, "test")
				So(pi["test"], ShouldEqual, opt)
			})

			Convey("Should obey the Omit() function", func() {
				Convey("If it returns false", func() {
					opt := &testOmitableOption{
						omit: false,
					}

					p.SetOptions(opt)
					So(pi, ShouldContainKey, "test")
					So(pi["test"], ShouldEqual, opt)
				})

				Convey("If it returns true", func() {
					opt := &testOmitableOption{
						omit: true,
					}

					p.SetOptions(opt)
					So(pi, ShouldNotContainKey, "test")
				})
			})

			Convey("Should use the Finalize() function", func() {
				opt := &testFinalizableOption{}
				So(opt.finalized, ShouldBeFalse)

				p.SetOptions(opt)
				So(opt.finalized, ShouldBeTrue)
			})

			Convey("Should handle existing keys", func() {
				Convey("Should replace by default", func() {
					opt1 := &testOption{}
					opt2 := &testOption{}

					p.SetOptions(opt1)
					So(pi, ShouldContainKey, "test")
					So(pi["test"], ShouldEqual, opt1)

					p.SetOptions(opt2)
					So(pi, ShouldContainKey, "test")
					So(pi["test"], ShouldEqual, opt2)
				})

				Convey("Should merge when Merge() is present", func() {
					opt1 := &testMergeableOption{data: 1}
					opt2 := &testMergeableOption{data: 2}

					p.SetOptions(opt1)
					So(pi, ShouldContainKey, "test")
					So(pi["test"], ShouldEqual, opt1)

					p.SetOptions(opt2)
					So(pi, ShouldContainKey, "test")
					So(opt1.data, ShouldEqual, 1)
					So(opt2.data, ShouldEqual, 2)
					So(pi["test"], ShouldResemble, &testMergeableOption{data: 3})
				})
			})
		})

		Convey("MarshalJSON", func() {
			p := NewPacket()

			Convey("With basic options", func() {
				opt := &testOption{}
				p.SetOptions(opt)

				b, err := json.Marshal(p)
				So(err, ShouldBeNil)

				var data map[string]interface{}
				So(json.Unmarshal(b, &data), ShouldBeNil)

				So(data, ShouldContainKey, "test")
				So(data["test"], ShouldResemble, map[string]interface{}{})
			})

			Convey("With custom MarshalJSON implementations", func() {
				opt := &testSerializableOption{data: "testing"}
				p.SetOptions(opt)

				b, err := json.Marshal(p)
				So(err, ShouldBeNil)

				var data map[string]interface{}
				So(json.Unmarshal(b, &data), ShouldBeNil)

				So(data, ShouldContainKey, "test")
				So(data["test"], ShouldResemble, "testing")
			})
		})
	})
}
