// +build go1.12
package sentry

import (
	"runtime/debug"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestModulesWithGomod(t *testing.T) {
	Convey("Modules", t, func() {
		_, ok := debug.ReadBuildInfo()

		if ok {
			Convey("Should register itself with the default providers", func() {
				opt := testGetOptionsProvider(Modules(map[string]string{"test": "correct"}))
				So(opt, ShouldNotBeNil)
			})
		}
	})
}
