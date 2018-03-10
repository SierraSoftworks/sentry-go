package sentry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleUser() {
	user := UserInfo{
		ID:        "17ba08f7cc89a912bf812918",
		Email:     "test@example.com",
		Username:  "Test User",
		IPAddress: "127.0.0.1",
		Extra: map[string]string{
			"role": "Tester",
		},
	}

	cl := NewClient(
		// You can specify your user when you create your client
		User(&user),
	)

	cl.Capture(
		// Or when you send an event to Sentry
		User(&user),
	)
}

func TestUser(t *testing.T) {
	Convey("User", t, func() {
		user := UserInfo{
			ID:        "17ba08f7cc89a912bf812918",
			Email:     "test@example.com",
			Username:  "Test User",
			IPAddress: "127.0.0.1",
			Extra: map[string]string{
				"role": "Tester",
			},
		}

		fields := map[string]string{
			"id":         "17ba08f7cc89a912bf812918",
			"email":      "test@example.com",
			"username":   "Test User",
			"ip_address": "127.0.0.1",
			"role":       "Tester",
		}

		Convey("User()", func() {
			Convey("Should return an Option", func() {
				So(User(&user), ShouldImplement, (*Option)(nil))
			})

			Convey("Should return nil if the user info is nil", func() {
				So(User(nil), ShouldBeNil)
			})

			Convey("Should use the correct Class()", func() {
				So(User(&user).Class(), ShouldEqual, "user")
			})

			Convey("Should have the correct fields set", func() {
				u := User(&user)
				So(u, ShouldNotBeNil)

				ui, ok := u.(*userOption)
				So(ok, ShouldBeTrue)

				So(ui.fields, ShouldResemble, fields)
			})

			Convey("MarshalJSON", func() {
				u := User(&user)
				So(u, ShouldNotBeNil)

				expected := map[string]interface{}{}
				for k, v := range fields {
					expected[k] = v
				}

				So(testOptionsSerialize(u), ShouldResemble, expected)
			})
		})
	})
}
