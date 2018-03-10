package sentry

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleTimestamp() {
	cl := NewClient()

	cl.Capture(
		// You can specify the timestamp when sending an event to Sentry
		Timestamp(time.Now()),
	)
}

func TestTimestamp(t *testing.T) {
	Convey("Timestamp", t, func() {
		Convey("Should register itself with the default providers", func() {
			opt := testGetOptionsProvider(Timestamp(time.Now()))
			So(opt, ShouldNotBeNil)
		})

		Convey("Timestamp()", func() {
			Convey("Should use the correct Class()", func() {
				So(Timestamp(time.Now()).Class(), ShouldEqual, "timestamp")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					Convey("Should marshal to a string", func() {
						t := time.Now()
						b, err := json.Marshal(Timestamp(t))
						So(err, ShouldBeNil)
						So(string(b), ShouldEqual, fmt.Sprintf(`"%s"`, t.UTC().Format("2006-01-02T15:04:05")))
					})
				})
			})
		})
	})
}
