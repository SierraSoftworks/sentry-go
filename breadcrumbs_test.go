package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleDefaultBreadcrumbs() {
	// We can change the maximum number of breadcrumbs to be stored
	DefaultBreadcrumbs().WithSize(5)

	DefaultBreadcrumbs().NewDefault(nil).WithMessage("This is an example")
	DefaultBreadcrumbs().NewDefault(map[string]interface{}{
		"example": true,
	}).WithMessage("It should give you an idea of how you can use breadcrumbs in your app")

	DefaultBreadcrumbs().
		NewNavigation("introduction", "navigation").
		WithCategory("navigation").
		WithMessage("You can use them to represent navigations from one page to another")

	DefaultBreadcrumbs().
		NewNavigation("navigation", "http").
		WithCategory("navigation").
		WithMessage("Or to represent changes in the state of your application's workflows")

	DefaultBreadcrumbs().
		NewHTTPRequest("GET", "https://example.com/api/v1/awesome", 200, "OK").
		WithLevel(Debug).
		WithMessage("I think we can agree that they're pretty awesome")

	NewClient().Capture(Message("Finally, we send the event with all our breadcrumbs included"))
}

func ExampleBreadcrumbs() {
	rootClient := NewClient()
	DefaultBreadcrumbs().NewDefault(nil).WithMessage("Breadcrumb in the default context")

	breadcrumbs := NewBreadcrumbsList(10)
	contextClient := rootClient.With(Breadcrumbs(breadcrumbs))
	breadcrumbs.NewDefault(nil).WithMessage("Breadcrumb in the private context")

	// Will include only the first breadcrumb
	rootClient.Capture(
		Message("Event in default context"),
		Logger("default"),
	)

	// Will include only the second breadcrumb
	contextClient.Capture(
		Message("Event in private context"),
		Logger("private"),
	)
}

func TestBreadcrumbs(t *testing.T) {
	Convey("Breadcrumbs", t, func() {
		Convey("Should be registered as a default options provider", func() {
			provider := testGetOptionsProvider(&breadcrumbsList{})
			So(provider, ShouldNotBeNil)
		})

		Convey("Should expose a DefaultBreadcrumbs() collection", func() {
			So(DefaultBreadcrumbs(), ShouldNotBeNil)
			So(DefaultBreadcrumbs(), ShouldImplement, (*BreadcrumbsList)(nil))
		})

		Convey("NewBreadcrumbsList()", func() {
			l := NewBreadcrumbsList(3)

			ll, ok := l.(*breadcrumbsList)
			So(ok, ShouldBeTrue)

			Convey("Should implement Option interface", func() {
				So(l, ShouldImplement, (*Option)(nil))
			})

			Convey("Should initialize correctly", func() {
				So(ll.MaxLength, ShouldEqual, 3)
				So(ll.Length, ShouldEqual, 0)
				So(ll.Head, ShouldBeNil)
				So(ll.Tail, ShouldBeNil)
			})

			Convey("NewDefault()", func() {
				data := map[string]interface{}{
					"test": true,
				}

				Convey("Should return a Breadcrumb", func() {
					b := l.NewDefault(data)
					So(b, ShouldImplement, (*Breadcrumb)(nil))
				})

				Convey("Should use a default breadcrumb type", func() {
					b := l.NewDefault(data)
					bb, ok := b.(*breadcrumb)
					So(ok, ShouldBeTrue)
					So(bb.Type, ShouldEqual, "")
					So(bb.Data, ShouldEqual, data)
				})

				Convey("Should allow you to specify nil for default data", func() {
					b := l.NewDefault(nil)
					bb, ok := b.(*breadcrumb)
					So(ok, ShouldBeTrue)
					So(bb.Data, ShouldResemble, map[string]interface{}{})
				})

				Convey("Should insert the breadcrumb into the list at its tail", func() {
					b := l.NewDefault(data)
					So(ll.Length, ShouldEqual, 1)
					So(ll.Tail, ShouldNotBeNil)
					So(ll.Tail.Value, ShouldEqual, b)
				})
			})

			Convey("NewNavigation()", func() {
				Convey("Should return a Breadcrumb", func() {
					b := l.NewNavigation("/from", "/to")
					So(b, ShouldImplement, (*Breadcrumb)(nil))
				})

				Convey("Should use a navigation breadcrumb type", func() {
					b := l.NewNavigation("/from", "/to")
					bb, ok := b.(*breadcrumb)
					So(ok, ShouldBeTrue)
					So(bb.Type, ShouldEqual, "navigation")
					So(bb.Data, ShouldResemble, map[string]interface{}{
						"from": "/from",
						"to":   "/to",
					})
				})

				Convey("Should insert the breadcrumb into the list at its tail", func() {
					b := l.NewNavigation("/from", "/to")
					So(ll.Length, ShouldEqual, 1)
					So(ll.Tail, ShouldNotBeNil)
					So(ll.Tail.Value, ShouldEqual, b)
				})
			})

			Convey("NewHTTPRequest()", func() {
				Convey("Should return a Breadcrumb", func() {
					b := l.NewHTTPRequest("GET", "/test", 200, "OK")
					So(b, ShouldImplement, (*Breadcrumb)(nil))
				})

				Convey("Should use a navigation breadcrumb type", func() {
					b := l.NewHTTPRequest("GET", "/test", 200, "OK")
					bb, ok := b.(*breadcrumb)
					So(ok, ShouldBeTrue)
					So(bb.Type, ShouldEqual, "http")
					So(bb.Data, ShouldResemble, map[string]interface{}{
						"method":      "GET",
						"url":         "/test",
						"status_code": 200,
						"reason":      "OK",
					})
				})

				Convey("Should insert the breadcrumb into the list at its tail", func() {
					b := l.NewHTTPRequest("GET", "/test", 200, "OK")
					So(ll.Length, ShouldEqual, 1)
					So(ll.Tail, ShouldNotBeNil)
					So(ll.Tail.Value, ShouldEqual, b)
				})
			})

			Convey("WithSize()", func() {
				Convey("Should be chainable", func() {
					So(l.WithSize(5), ShouldEqual, l)
				})

				Convey("Should update the max length field", func() {
					So(ll.MaxLength, ShouldEqual, 3)
					l.WithSize(5)
					So(ll.MaxLength, ShouldEqual, 5)
				})

				Convey("Should remove elements which push the length over the limit", func() {
					var b Breadcrumb
					for i := 0; i < 3; i++ {
						b = l.NewDefault(map[string]interface{}{
							"index": i,
						})
						So(b, ShouldNotBeNil)
					}

					So(ll.Length, ShouldEqual, 3)
					l.WithSize(1)
					So(ll.Length, ShouldEqual, 1)
					So(ll.Head, ShouldNotBeNil)
					So(ll.Head.Next, ShouldBeNil)
					So(ll.Head.Value, ShouldEqual, b)

					So(b, ShouldHaveSameTypeAs, &breadcrumb{})
					So(b.(*breadcrumb).Data, ShouldResemble, map[string]interface{}{
						"index": 2,
					})
				})
			})

			Convey("append()", func() {
				Convey("Should evict older entries when we run over the max length", func() {
					var b Breadcrumb
					for i := 0; i < 10; i++ {
						b = l.NewDefault(map[string]interface{}{
							"index": i,
						})
						So(b, ShouldNotBeNil)
					}

					So(ll.Length, ShouldEqual, 3)

					So(b, ShouldHaveSameTypeAs, &breadcrumb{})
					So(b.(*breadcrumb).Data, ShouldResemble, map[string]interface{}{
						"index": 9,
					})
				})
			})

			Convey("list()", func() {
				Convey("Should handle an empty list correctly", func() {
					ol := ll.list()
					So(ol, ShouldNotBeNil)
					So(ol, ShouldHaveLength, 0)
				})

				Convey("Should expose a non-empty list correctly", func() {
					for i := 0; i < 3; i++ {
						So(l.NewDefault(map[string]interface{}{
							"index": i,
						}), ShouldNotBeNil)
					}

					ol := ll.list()
					So(ol, ShouldNotBeNil)
					So(ol, ShouldHaveLength, 3)
					for i, item := range ol {
						So(item, ShouldHaveSameTypeAs, &breadcrumb{})
						So(item.(*breadcrumb).Data, ShouldContainKey, "index")
						So(item.(*breadcrumb).Data["index"], ShouldEqual, i)
					}
				})
			})
		})

		Convey("MarshalJSON", func() {
			l := NewBreadcrumbsList(5)
			l.NewDefault(map[string]interface{}{"test": true})

			data := testOptionsSerialize(Breadcrumbs(l))
			So(data, ShouldNotBeNil)
			So(data, ShouldHaveSameTypeAs, []interface{}{})
		})
	})
}
