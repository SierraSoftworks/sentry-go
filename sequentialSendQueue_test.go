package sentry

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSequentialSendQueue(t *testing.T) {
	Convey("SequentialSendQueue", t, func() {
		Convey("NewSequentialSendQueue()", func() {
			q := NewSequentialSendQueue(10)
			So(q, ShouldNotBeNil)
			So(q, ShouldImplement, (*SendQueue)(nil))
			defer q.Shutdown(true)

			So(q, ShouldHaveSameTypeAs, &sequentialSendQueue{})
		})

		Convey("Send()", func() {
			dsn := "http://user:pass@example.com/sentry/1"
			transport := testNewTestTransport()
			So(transport, ShouldNotBeNil)

			cl := NewClient(
				DSN(dsn),
				UseTransport(transport),
			)

			cfg, ok := cl.(Config)
			So(ok, ShouldBeTrue)

			Convey("Normal Operation", func() {
				q := NewSequentialSendQueue(10)
				So(q, ShouldNotBeNil)
				defer q.Shutdown(true)

				p := NewPacket()
				So(p, ShouldNotBeNil)

				e := q.Enqueue(cfg, p)
				So(e, ShouldNotBeNil)

				select {
				case p2 := <-transport.ch:
					So(p2, ShouldEqual, p)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}

				select {
				case err, ok := <-e.WaitChannel():
					So(err, ShouldBeNil)
					So(ok, ShouldBeFalse)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("Buffer Overflow", func() {
				q := NewSequentialSendQueue(0)
				So(q, ShouldNotBeNil)
				defer q.Shutdown(true)

				p := NewPacket()
				So(p, ShouldNotBeNil)

				time.Sleep(1 * time.Millisecond)

				// First entry will be processed
				e1 := q.Enqueue(cfg, p)
				So(e1, ShouldNotBeNil)

				// Second will be failed
				e2 := q.Enqueue(cfg, p)
				So(e2, ShouldNotBeNil)

				select {
				case <-transport.ch:
					So(e1.Error(), ShouldBeNil)
				case err := <-e1.WaitChannel():
					So(err, ShouldBeNil)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}

				select {
				case <-transport.ch:
					So(fmt.Errorf("shouldn't send"), ShouldBeNil)
				case err := <-e2.WaitChannel():
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, ErrSendQueueFull.Error())
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("When Shutdown", func() {
				q := NewSequentialSendQueue(10)
				So(q, ShouldNotBeNil)
				q.Shutdown(true)

				p := NewPacket()
				So(p, ShouldNotBeNil)

				e := q.Enqueue(cfg, p)
				So(e, ShouldNotBeNil)

				select {
				case <-transport.ch:
					So(fmt.Errorf("shouldn't send"), ShouldBeNil)
				case err := <-e.WaitChannel():
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, ErrSendQueueShutdown.Error())
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})
		})

		Convey("Shutdown()", func() {
			Convey("Should be safe to call repeatedly", func() {
				q := NewSequentialSendQueue(0)
				q.Shutdown(true)
				q.Shutdown(true)
			})
		})
	})
}
