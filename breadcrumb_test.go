package sentry

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBreadcrumb(t *testing.T) {
	Convey("Breadcrumb", t, func() {
		data := map[string]interface{}{
			"test": true,
		}

		Convey("newBreadcrumb", func() {

			Convey("Should return a Breadcrumb type", func() {
				b := newBreadcrumb("default", data)
				So(b, ShouldNotBeNil)

				So(b, ShouldImplement, (*Breadcrumb)(nil))
			})

			Convey("Should set the timestamp", func() {
				b := newBreadcrumb("default", data)
				So(b, ShouldNotBeNil)
				So(b.Timestamp, ShouldNotEqual, 0)
			})

			Convey("Should set the data", func() {
				b := newBreadcrumb("default", data)
				So(b, ShouldNotBeNil)
				So(b.Data, ShouldEqual, data)
			})

			Convey("Should set the Type correctly", func() {
				Convey("With default type", func() {
					b := newBreadcrumb("default", data)

					So(b, ShouldNotBeNil)
					So(b.Type, ShouldEqual, "")
				})

				Convey("With non-default type", func() {
					b := newBreadcrumb("test", data)

					So(b, ShouldNotBeNil)
					So(b.Type, ShouldEqual, "test")
				})
			})
		})

		Convey("WithMessage()", func() {
			b := newBreadcrumb("default", data)
			So(b, ShouldNotBeNil)

			Convey("Should update the Message field", func() {
				b.WithMessage("test")
				So(b.Message, ShouldEqual, "test")
			})

			Convey("Should be chainable", func() {
				So(b.WithMessage("test"), ShouldEqual, b)
			})
		})

		Convey("WithCategory()", func() {
			b := newBreadcrumb("default", data)
			So(b, ShouldNotBeNil)

			Convey("Should update the Category field", func() {
				b.WithCategory("test")
				So(b.Category, ShouldEqual, "test")
			})

			Convey("Should be chainable", func() {
				So(b.WithCategory("test"), ShouldEqual, b)
			})
		})

		Convey("WithLevel()", func() {
			b := newBreadcrumb("default", data)
			So(b, ShouldNotBeNil)

			Convey("Should update the Level field", func() {
				b.WithLevel(Error)
				So(b.Level, ShouldEqual, Error)
			})

			Convey("Should be chainable", func() {
				So(b.WithLevel(Error), ShouldEqual, b)
			})
		})

		Convey("WithTimestamp()", func() {
			b := newBreadcrumb("default", data)
			So(b, ShouldNotBeNil)

			Convey("Should update the Timestamp field", func() {
				now := time.Now()
				b.WithTimestamp(now)
				So(b.Timestamp, ShouldEqual, now.UTC().Unix())
			})

			Convey("Should be chainable", func() {
				So(b.WithTimestamp(time.Now()), ShouldEqual, b)
			})
		})
	})
}
