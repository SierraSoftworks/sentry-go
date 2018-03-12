package sentry

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleQueuedEvent() {
	cl := NewClient()

	e := cl.Capture(
		Message("Example Event"),
	)

	// You can then wait on the event to be sent
	e.Wait()

	// Or you can use the WaitChannel if you want support for timeouts
	select {
	case err := <-e.WaitChannel():
		if err != nil {
			fmt.Println("failed to send event: ", err)
		} else {
			// You can also get the EventID for this event
			fmt.Println("sent event: ", e.EventID())
		}
	case <-time.After(time.Second):
		fmt.Println("timed out waiting for send")
	}
}

func ExampleQueuedEventInternal() {
	// If you're implementing your own send queue, you will want to use
	// the QueuedEventInternal to control when events are finished and
	// to access the packet and config related to them.

	cl := NewClient()
	e := cl.Capture()

	if ei, ok := e.(QueuedEventInternal); ok {
		// Get the packet for the event
		p := ei.Packet()

		// Get the config for the event
		cfg := ei.Config()

		// Use the configured transport to send the packet
		err := cfg.Transport().Send(cfg.DSN(), p)

		// Complete the event (with the error, if not nil)
		ei.Complete(err)
	}
}

func TestQueuedEvent(t *testing.T) {
	Convey("QueuedEvent", t, func() {
		id, err := NewEventID()
		So(err, ShouldBeNil)

		cl := NewClient(
			DSN(""),
		)

		cfg, ok := cl.(Config)
		So(ok, ShouldBeTrue)

		p := NewPacket().SetOptions(EventID(id))

		Convey("NewQueuedEvent()", func() {
			e := NewQueuedEvent(cfg, p)
			So(e, ShouldNotBeNil)

			ei, ok := e.(*queuedEvent)
			So(ok, ShouldBeTrue)
			So(ei.cfg, ShouldEqual, cfg)
			So(ei.packet, ShouldEqual, p)
		})

		Convey("EventID()", func() {
			Convey("With a valid packet", func() {
				e := NewQueuedEvent(cfg, p)
				So(e.EventID(), ShouldEqual, id)
			})

			Convey("With an invalid packet", func() {
				e := NewQueuedEvent(cfg, nil)
				So(e.EventID(), ShouldEqual, "")
			})
		})

		Convey("Wait()", func() {
			Convey("When it isn't yet complete", func(c C) {
				e := NewQueuedEvent(cfg, p)
				ch := make(chan struct{})
				defer close(ch)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				go func() {
					ch <- struct{}{}
					e.Wait()
					c.So(e.Error(), ShouldBeNil)
					ch <- struct{}{}
				}()

				// Wait for it to be waiting
				<-ch
				ei.Complete(nil)

				select {
				case <-ch:
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("When it is already complete", func(c C) {
				e := NewQueuedEvent(cfg, p)
				ch := make(chan struct{})
				defer close(ch)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				ei.Complete(nil)

				go func() {
					e.Wait()
					c.So(e.Error(), ShouldBeNil)
					ch <- struct{}{}
				}()

				select {
				case <-ch:
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})
		})

		Convey("WaitChannel()", func() {
			Convey("When it isn't yet complete", func() {
				e := NewQueuedEvent(cfg, p)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				go func() {
					ei.Complete(nil)
				}()

				select {
				case err := <-e.WaitChannel():
					So(err, ShouldBeNil)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("When it is already complete", func() {
				e := NewQueuedEvent(cfg, p)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				ei.Complete(nil)

				select {
				case err := <-e.WaitChannel():
					So(err, ShouldBeNil)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("When it has an error", func() {
				e := NewQueuedEvent(cfg, p)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				err := fmt.Errorf("example error")
				ei.Complete(err)

				select {
				case err := <-e.WaitChannel():
					So(err, ShouldEqual, err)
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})
		})

		Convey("Error()", func() {
			Convey("When it isn't yet complete", func(c C) {
				e := NewQueuedEvent(cfg, p)
				ch := make(chan struct{})
				defer close(ch)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				err := fmt.Errorf("example error")

				go func() {
					ch <- struct{}{}
					c.So(e.Error(), ShouldEqual, err)
					ch <- struct{}{}
				}()

				// Wait for it to be waiting
				<-ch
				ei.Complete(err)

				select {
				case <-ch:
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})

			Convey("When it is already complete", func(c C) {
				e := NewQueuedEvent(cfg, p)
				ch := make(chan struct{})
				defer close(ch)

				ei, ok := e.(QueuedEventInternal)
				So(ok, ShouldBeTrue)

				err := fmt.Errorf("example error")
				ei.Complete(err)

				go func() {
					c.So(e.Error(), ShouldEqual, err)
					ch <- struct{}{}
				}()

				select {
				case <-ch:
				case <-time.After(100 * time.Millisecond):
					So(fmt.Errorf("timeout"), ShouldBeNil)
				}
			})
		})

		Convey("QueuedEventInternal", func() {
			e := NewQueuedEvent(cfg, p)
			So(e, ShouldNotBeNil)

			ei, ok := e.(QueuedEventInternal)
			So(ok, ShouldBeTrue)

			Convey("Complete()", func() {
				Convey("Should set the error", func() {
					err := fmt.Errorf("example error")
					ei.Complete(err)
					So(e.Error(), ShouldEqual, err)
				})

				Convey("When called multiple times", func() {
					ei.Complete(nil)
					ei.Complete(nil)
					So(e.Error(), ShouldBeNil)

					Convey("Should not change the event's details", func() {
						ei.Complete(fmt.Errorf("example error"))
						So(e.Error(), ShouldBeNil)
					})
				})
			})
		})
	})
}
