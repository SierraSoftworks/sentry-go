package sentry

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrType(t *testing.T) {
	Convey("ErrType", t, func() {
		const errType = ErrType("sentry: this is a test error")

		Convey("IsInstance()", func() {
			So(errType.IsInstance(errType), ShouldBeTrue)
			So(errType.IsInstance(errors.New(errType.Error())), ShouldBeTrue)
		})

		Convey("Error()", func() {
			So(errType.Error(), ShouldEqual, "sentry: this is a test error")
		})
	})
}
