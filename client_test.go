package sentry

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleClient() {
	// You can create a new root client directly and configure
	// it by passing any options you wish
	cl := NewClient(
		DSN(""),
	)

	// You can then create a derivative client with any context-specific
	// options. These are useful if you want to encapsulate context-specific
	// information like the HTTP request that is being handled.
	var r *http.Request
	ctxCl := cl.With(
		HTTPRequest(r).WithHeaders(),
		Logger("http"),
	)

	// You can then use the client to capture an event and send it to Sentry
	err := fmt.Errorf("an error occurred")
	ctxCl.Capture(
		ExceptionForError(err),
	)
}

func ExampleDefaultClient() {
	DefaultClient().Capture(
		Message("This is an example message"),
	)
}

func TestClient(t *testing.T) {
	Convey("Client", t, func() {
		Convey("DefaultClient()", func() {
			So(DefaultClient(), ShouldNotBeNil)
			So(DefaultClient(), ShouldImplement, (*Client)(nil))
			So(DefaultClient(), ShouldEqual, defaultClient)
		})

		Convey("NewClient()", func() {
			Convey("Should return a Client", func() {
				So(NewClient(), ShouldImplement, (*Client)(nil))
			})

			Convey("Should set the client's parent to nil", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.parent, ShouldBeNil)
			})

			Convey("Should set the client's options", func() {
				opt := &testOption{}

				cl := NewClient(opt)
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.options, ShouldResemble, []Option{opt})
			})
		})

		Convey("Capture()", func() {
			tr := testNewTestTransport()
			So(tr, ShouldNotBeNil)

			cl := NewClient(UseTransport(tr))
			So(cl, ShouldNotBeNil)

			e := cl.Capture(Message("test"))
			So(e, ShouldNotBeNil)

			ei, ok := e.(QueuedEventInternal)
			So(ok, ShouldBeTrue)

			select {
			case p := <-tr.ch:
				So(p, ShouldEqual, ei.Packet())

				pi, ok := p.(*packet)
				So(ok, ShouldBeTrue)
				So((*pi)[Message("test").Class()], ShouldResemble, Message("test"))
			case <-time.After(100 * time.Millisecond):
				So(fmt.Errorf("timeout"), ShouldBeNil)
			}
		})

		Convey("UseSendQueue()", func() {
			Convey("Should set the client's queue", func() {
				q := NewSequentialSendQueue(0)
				So(q, ShouldNotBeNil)

				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.queue, ShouldBeNil)

				So(cl.UseSendQueue(q), ShouldEqual, cl)
				So(cll.queue, ShouldEqual, q)
			})
		})

		Convey("With()", func() {
			Convey("Should return a new Client", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				ctxCl := cl.With()
				So(ctxCl, ShouldNotBeNil)
				So(ctxCl, ShouldNotEqual, cl)
				So(ctxCl, ShouldImplement, (*Client)(nil))
			})

			Convey("Should set the client's parent", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				ctxCl := cl.With()
				So(ctxCl, ShouldNotBeNil)

				cll, ok := ctxCl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.parent, ShouldEqual, cl)
			})

			Convey("Should set the client's options", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				opt := &testOption{}
				ctxCl := cl.With(opt)
				So(ctxCl, ShouldNotBeNil)

				cll, ok := ctxCl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.options, ShouldResemble, []Option{opt})
			})
		})

		Convey("fullDefaultOptions()", func() {
			Convey("Should include the default providers' options", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)

				opts := cll.fullDefaultOptions()
				So(opts, ShouldNotBeNil)
				So(opts, ShouldNotBeEmpty)

				i := 0
				for _, provider := range defaultOptionProviders {
					opt := provider()
					if opt == nil {
						continue
					}

					Convey(fmt.Sprintf("%s (%d)", opt.Class(), i), func() {
						So(opts[i], ShouldHaveSameTypeAs, opt)
					})

					i++
				}
			})

			Convey("With a root client", func() {
				opt := &testOption{}

				cl := NewClient(opt)
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)

				So(cll.fullDefaultOptions(), ShouldContain, opt)
			})

			Convey("With a derived client", func() {
				opt1 := &testMergeableOption{data: 1}
				opt2 := &testMergeableOption{data: 2}

				cl := NewClient(opt1)
				So(cl, ShouldNotBeNil)

				dcl := cl.With(opt2)
				So(dcl, ShouldNotBeNil)

				cll, ok := dcl.(*client)
				So(ok, ShouldBeTrue)

				opts := cll.fullDefaultOptions()
				So(opts, ShouldContain, opt1)
				So(opts, ShouldContain, opt2)
			})
		})

		Convey("getQueue", func() {
			Convey("With no custom queue", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.queue, ShouldBeNil)

				So(cll.getQueue(), ShouldEqual, DefaultSendQueue())
			})

			Convey("With a custom queue", func() {
				q := NewSequentialSendQueue(0)
				So(q, ShouldNotBeNil)

				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cl.UseSendQueue(q)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.queue, ShouldEqual, q)

				So(cll.getQueue(), ShouldEqual, q)
			})

			Convey("With no custom queue on a parent", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				dcl := cl.With()
				So(dcl, ShouldNotBeNil)

				cll, ok := dcl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.queue, ShouldBeNil)

				So(cll.getQueue(), ShouldEqual, DefaultSendQueue())
			})

			Convey("With a custom queue on a parent", func() {
				q := NewSequentialSendQueue(0)
				So(q, ShouldNotBeNil)

				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cl.UseSendQueue(q)

				dcl := cl.With()
				So(dcl, ShouldNotBeNil)

				cll, ok := dcl.(*client)
				So(ok, ShouldBeTrue)
				So(cll.queue, ShouldBeNil)

				So(cll.getQueue(), ShouldEqual, q)
			})
		})

		Convey("getConfig()", func() {
			Convey("Should return a config option", func() {
				cl := NewClient()
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)

				So(cll.getConfig(), ShouldHaveSameTypeAs, &configOption{})
			})

			Convey("Should include the most recent DSN", func() {
				cl := NewClient(
					DSN("old"),
					DSN("new"),
				)
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)

				cnf := cll.getConfig()
				So(cnf, ShouldNotBeNil)
				So(cnf.DSN(), ShouldEqual, "new")
			})

			Convey("Should allow options to be overridden by the event", func() {
				cl := NewClient(
					DSN("old"),
					DSN("new"),
				)
				So(cl, ShouldNotBeNil)

				cll, ok := cl.(*client)
				So(ok, ShouldBeTrue)

				cnf := cll.getConfig(DSN("event"))
				So(cnf, ShouldNotBeNil)
				So(cnf.DSN(), ShouldEqual, "event")
			})
		})
	})
}
