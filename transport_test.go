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
				t := newHTTPTransport()
				So(UseTransport(t), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if no queue is provided", func() {
				So(UseTransport(nil), ShouldEqual, nil)
			})

			Convey("Should use the correct Class()", func() {
				t := newHTTPTransport()
				So(UseTransport(t).Class(), ShouldEqual, "sentry-go.transport")
			})

			Convey("Should implement Omit() and always return true", func() {
				t := newHTTPTransport()
				o := UseTransport(t)
				So(o, ShouldImplement, (*OmitableOption)(nil))
				So(o.(OmitableOption).Omit(), ShouldBeTrue)
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
