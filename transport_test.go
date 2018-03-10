package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleUseTransport() {
	var myTransport Transport

	cl := NewClient(
		// You can configure the transport to be used on a client level
		UseTransport(myTransport),
	)

	cl.Capture(
		// Or for a specific event when it is sent
		UseTransport(myTransport),
	)
}

func TestTransport(t *testing.T) {
	Convey("Transport", t, func() {
		Convey("UseTransport()", func() {
			Convey("Should return an Option", func() {
				So(UseTransport(newHTTPTransport()), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if the transport is nil", func() {
				So(UseTransport(nil), ShouldBeNil)
			})
		})

		Convey("DefaultTransport()", func() {
			Convey("Should be defined from the start", func() {
				So(DefaultTransport(), ShouldNotBeNil)
			})

			Convey("Should return the correct transport", func() {
				So(DefaultTransport(), ShouldEqual, defaultTransport)
			})
		})

		Convey("SetDefaultTransport()", func() {
			Convey("Should allow you to change the default transport", func() {
				t := newHTTPTransport()
				SetDefaultTransport(t)
				So(DefaultTransport(), ShouldEqual, t)
			})

			Convey("Should set the transport back to the default if nil is provided", func() {
				SetDefaultTransport(nil)
				So(DefaultTransport(), ShouldNotBeNil)
			})
		})
	})
}

func testNewTestTransport() *testTransport {
	return &testTransport{
		ch: make(chan Packet),
	}
}

type testTransport struct {
	ch  chan Packet
	err error
}

func (t *testTransport) Send(dsn string, packet Packet) error {
	t.ch <- packet
	return t.err
}
