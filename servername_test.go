package sentry

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleServerName() {
	cl := NewClient(
		// You can set the logger when you create your client
		ServerName("web01"),
	)

	cl.Capture(
		// You can also specify it when sending an event
		ServerName("web01.prod"),
	)
}

func TestServerName(t *testing.T) {
	Convey("ServerName", t, func() {
		Convey("Should register itself with the default providers", func() {
			opt := testGetOptionsProvider(ServerName(""))
			So(opt, ShouldNotBeNil)
		})

		Convey("ServerName()", func() {
			Convey("Should use the correct Class()", func() {
				So(ServerName("test").Class(), ShouldEqual, "server_name")
			})

			Convey("MarshalJSON", func() {
				Convey("Should marshal to a string", func() {
					Convey("Should marshal to a string", func() {
						b, err := json.Marshal(ServerName("test"))
						So(err, ShouldBeNil)
						So(string(b), ShouldEqual, `"test"`)
					})
				})
			})
		})
	})
}
