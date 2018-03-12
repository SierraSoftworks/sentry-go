package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleUseSendQueue() {
	cl := NewClient(
		// You can override the send queue on your root client
		// All of its derived clients will inherit this queue
		UseSendQueue(NewSequentialSendQueue(10)),
	)

	cl.With(
		// Or you can override it on a derived client
		UseSendQueue(NewSequentialSendQueue(10)),
	)
}

func TestSendQueue(t *testing.T) {
	Convey("SendQueue", t, func() {
		Convey("UseSendQueue()", func() {
			Convey("Should return an Option", func() {
				q := NewSequentialSendQueue(0)
				So(UseSendQueue(q), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if no queue is provided", func() {
				So(UseSendQueue(nil), ShouldEqual, nil)
			})

			Convey("Should use the correct Class()", func() {
				q := NewSequentialSendQueue(0)
				So(UseSendQueue(q).Class(), ShouldEqual, "sentry-go.sendqueue")
			})

			Convey("Should implement Omit() and always return true", func() {
				q := NewSequentialSendQueue(0)
				o := UseSendQueue(q)
				So(o, ShouldImplement, (*OmitableOption)(nil))
				So(o.(OmitableOption).Omit(), ShouldBeTrue)
			})
		})
	})
}
