package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSendQueue(t *testing.T) {
	Convey("SendQueue", t, func() {
		Convey("DefaultSendQueue()", func() {
			So(DefaultSendQueue(), ShouldNotBeNil)
			So(DefaultSendQueue(), ShouldHaveSameTypeAs, NewSequentialSendQueue(0))
			So(DefaultSendQueue(), ShouldEqual, defaultSendQueue)
		})

		Convey("SetDefaultSendQueue()", func() {
			q := NewSequentialSendQueue(0)
			SetDefaultSendQueue(q)
			So(defaultSendQueue, ShouldEqual, q)

			Convey("When you specify nil it should use the default", func() {
				SetDefaultSendQueue(nil)
				So(defaultSendQueue, ShouldNotBeNil)
			})
		})
	})
}
