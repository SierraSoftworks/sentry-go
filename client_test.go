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

		Convey("GetOption()", func() {
			Convey("Should return nil for an unrecognized option", func() {
				cl := NewClient()
				So(cl.GetOption("unknown-option-class"), ShouldBeNil)
			})

			Convey("Should skip over nil options", func() {
				cl := NewClient(nil)
				So(cl.GetOption("unknown-option-class"), ShouldBeNil)
			})

			Convey("Should return the option if it is present", func() {
				opt := &testOption{}
				cl := NewClient(opt)
				So(cl.GetOption("test"), ShouldEqual, opt)
			})

			Convey("Should return the most recent non-mergeable option", func() {
				opt := &testOption{}
				cl := NewClient(&testOption{}, opt)
				So(cl.GetOption("test"), ShouldEqual, opt)
			})

			Convey("Should merge options when supported", func() {
				cl := NewClient(&testMergeableOption{1}, &testMergeableOption{2})
				So(cl.GetOption("test"), ShouldResemble, &testMergeableOption{3})
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

		Convey("Config Interface", func() {
			Convey("DSN()", func() {
				Convey("In a world with no DSNs", func() {
					oldDefaultOptionProviders := defaultOptionProviders
					defer func() {
						defaultOptionProviders = oldDefaultOptionProviders
					}()

					defaultOptionProviders = []func() Option{}
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)
					So(cfg.DSN(), ShouldEqual, "")
				})

				Convey("When someone has implemented their own custom DSN option", func() {
					cl := NewClient(&testCustomClassOption{"sentry-go.dsn"})
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)
					So(cfg.DSN(), ShouldEqual, "")
				})

				Convey("With no custom DSN", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.DSN(), ShouldEqual, "")
				})

				Convey("With a custom DSN", func() {
					cl := NewClient(DSN("test"))
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.DSN(), ShouldEqual, "test")
				})

				Convey("With no custom DSN on a parent", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.DSN(), ShouldEqual, "")
				})

				Convey("With a custom DSN on a parent", func() {
					cl := NewClient(DSN("test"))
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)
					So(cfg.DSN(), ShouldEqual, "test")
				})
			})

			Convey("SendQueue()", func() {
				Convey("In a world with no SendQueues", func() {
					oldDefaultOptionProviders := defaultOptionProviders
					defer func() {
						defaultOptionProviders = oldDefaultOptionProviders
					}()

					defaultOptionProviders = []func() Option{}
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					q := cfg.SendQueue()
					So(q, ShouldNotBeNil)
					So(q, ShouldHaveSameTypeAs, NewSequentialSendQueue(0))
				})

				Convey("When someone has implemented their own custom SendQueue option", func() {
					cl := NewClient(&testCustomClassOption{"sentry-go.sendqueue"})
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					q := cfg.SendQueue()
					So(q, ShouldNotBeNil)
					So(q, ShouldHaveSameTypeAs, NewSequentialSendQueue(0))
				})

				Convey("With no custom queue", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.SendQueue(), ShouldEqual, DefaultClient().GetOption("sentry-go.sendqueue").(*sendQueueOption).queue)
				})

				Convey("With a custom queue", func() {
					q := NewSequentialSendQueue(0)
					So(q, ShouldNotBeNil)

					cl := NewClient(UseSendQueue(q))
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.SendQueue(), ShouldEqual, q)
				})

				Convey("With no custom queue on a parent", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.SendQueue(), ShouldEqual, DefaultClient().GetOption("sentry-go.sendqueue").(*sendQueueOption).queue)
				})

				Convey("With a custom queue on a parent", func() {
					q := NewSequentialSendQueue(0)
					So(q, ShouldNotBeNil)

					cl := NewClient(UseSendQueue(q))
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.SendQueue(), ShouldEqual, q)
				})
			})

			Convey("Transport()", func() {
				Convey("In a world with no Transports", func() {
					oldDefaultOptionProviders := defaultOptionProviders
					defer func() {
						defaultOptionProviders = oldDefaultOptionProviders
					}()

					defaultOptionProviders = []func() Option{}
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					t := cfg.Transport()
					So(t, ShouldNotBeNil)
					So(t, ShouldHaveSameTypeAs, newHTTPTransport())
				})

				Convey("When someone has implemented their own custom Transport option", func() {
					cl := NewClient(&testCustomClassOption{"sentry-go.transport"})
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					t := cfg.Transport()
					So(t, ShouldNotBeNil)
					So(t, ShouldHaveSameTypeAs, newHTTPTransport())
				})

				Convey("With no custom transport", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.Transport(), ShouldEqual, DefaultClient().GetOption("sentry-go.transport").(*transportOption).transport)
				})

				Convey("With a custom transport", func() {
					t := newHTTPTransport()
					So(t, ShouldNotBeNil)

					cl := NewClient(UseTransport(t))
					So(cl, ShouldNotBeNil)

					cfg, ok := cl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.Transport(), ShouldEqual, t)
				})

				Convey("With no custom transport on a parent", func() {
					cl := NewClient()
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.Transport(), ShouldEqual, DefaultClient().GetOption("sentry-go.transport").(*transportOption).transport)
				})

				Convey("With a custom transport on a parent", func() {
					t := newHTTPTransport()
					So(t, ShouldNotBeNil)

					cl := NewClient(UseTransport(t))
					So(cl, ShouldNotBeNil)

					dcl := cl.With()
					So(dcl, ShouldNotBeNil)

					cfg, ok := dcl.(Config)
					So(ok, ShouldBeTrue)

					So(cfg.Transport(), ShouldEqual, t)
				})
			})
		})
	})
}
