package sentry

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleEventID() {
	id, err := NewEventID()
	if err != nil {
		log.Fatalln(err)
	}

	cl := NewClient()

	ctxCl := cl.With(
		// You could set the event ID for a context specific
		// client if you wanted (but you probably shouldn't).
		EventID(id),
	)

	ctxCl.Capture(
		// The best place to set it is when you are ready to send
		// an event to Sentry.
		EventID(id),
	)
}

func TestEventID(t *testing.T) {
	Convey("EventID", t, func() {
		Convey("NewEventID()", func() {
			Convey("Should return a valid event ID", func() {
				id, err := NewEventID()
				So(err, ShouldBeNil)
				So(id, ShouldHaveLength, 32)

				for _, r := range id {
					So(r, ShouldBeIn, []rune("0123456789abcdef"))
				}
			})
		})

		id, err := NewEventID()
		So(err, ShouldBeNil)

		Convey("EventID()", func() {
			id, err := NewEventID()
			So(err, ShouldBeNil)

			Convey("Should return an Option", func() {
				So(EventID(id), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if the ID is invalid", func() {
				So(EventID("invalid"), ShouldBeNil)
				So(EventID("xx23456789abcdef0123456789abcdef"), ShouldBeNil)
			})
		})

		Convey("Should use the correct class", func() {
			So(EventID(id).Class(), ShouldEqual, "event_id")
		})

		Convey("MarshalJSON", func() {
			So(testOptionsSerialize(EventID(id)), ShouldResemble, id)
		})

		Convey("Packet Extensions", func() {
			Convey("getEventID()", func() {
				p := NewPacket().SetOptions(EventID(id))
				pp := p.(*packet)

				So(pp.getEventID(), ShouldEqual, id)
			})

		})
	})
}
