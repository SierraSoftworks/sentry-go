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

			cfg := &configOption{
				dsn:       &dsn,
				transport: transport,
			}

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

				// First entry will be processed
				e1 := q.Enqueue(cfg, p)
				So(e1, ShouldNotBeNil)

				// Second will be failed
				e2 := q.Enqueue(cfg, p)
				So(e2, ShouldNotBeNil)

				select {
				case <-transport.ch:
					So(e1.Error(), ShouldBeNil)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}

				select {
				case <-transport.ch:
					So(fmt.Errorf("shouldn't send"), ShouldBeNil)
				case <-time.After(10 * time.Millisecond):
					So(e2.Error(), ShouldNotBeNil)
					So(e2.Error().Error(), ShouldContainSubstring, ErrSendQueueFull.Error())
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
				case <-time.After(10 * time.Millisecond):
					So(e.Error(), ShouldNotBeNil)
					So(e.Error().Error(), ShouldContainSubstring, ErrSendQueueShutdown.Error())
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
