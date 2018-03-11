package sentry

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleAddInternalPrefixes() {
	// This adds the provided prefixes to your list of internal
	// package prefixes used to tag stacktrace frames as in-app.
	AddInternalPrefixes("github.com/SierraSoftworks/sentry-go")
}

func ExampleStackTrace() {
	cl := NewClient()

	cl.Capture(
		// You can specify that a StackTrace should be included when
		// sending your event to Sentry
		StackTrace().
			// You can also gather the stacktrace frames from a specific
			// error if it is created using `pkg/errors`
			ForError(errors.New("example error")).
			// And you can mark frames as "internal" by specifying the
			// internal frame prefixes here.
			WithInternalPrefixes(
				"github.com/SierraSoftworks/sentry-go",
			),
	)
}

func TestStackTrace(t *testing.T) {
	Convey("StackTrace", t, func() {
		Convey("AddInternalPrefixes()", func() {
			AddInternalPrefixes("github.com/SierraSoftworks/sentry-go")
			So(defaultInternalPrefixes, ShouldResemble, []string{"main", "github.com/SierraSoftworks/sentry-go"})
		})

		Convey("StackTrace()", func() {
			Convey("Should return an Option", func() {
				So(StackTrace(), ShouldImplement, (*Option)(nil))
			})

			Convey("Should let you collect stacktrace frames from an error", func() {
				st := StackTrace()
				So(st, ShouldNotBeNil)

				err := errors.New("example error")
				So(st.ForError(err), ShouldEqual, st)
			})

			Convey("Should allow you to set internal package prefixes", func() {
				st := StackTrace()
				So(st, ShouldNotBeNil)

				sti, ok := st.(*stackTraceOption)
				So(ok, ShouldBeTrue)
				So(sti.internalPrefixes, ShouldResemble, defaultInternalPrefixes)

				st.WithInternalPrefixes("github.com/SierraSoftworks/sentry-go")
				st.WithInternalPrefixes("github.com/SierraSoftworks")
				So(sti.internalPrefixes, ShouldContain, "github.com/SierraSoftworks")
				So(sti.internalPrefixes, ShouldContain, "github.com/SierraSoftworks/sentry-go")
			})
		})

		Convey("Should implement Finalize()", func() {
			So(StackTrace(), ShouldImplement, (*FinalizableOption)(nil))
		})

		Convey("Finalize()", func() {
			st := StackTrace().WithInternalPrefixes("github.com/SierraSoftworks/sentry-go")
			So(st, ShouldNotBeNil)

			sti, ok := st.(*stackTraceOption)
			So(ok, ShouldBeTrue)

			sti.Finalize()

			So(len(sti.Frames), ShouldBeGreaterThan, 0)
			So(sti.Frames[len(sti.Frames)-1].InApp, ShouldBeTrue)
		})

		Convey("Should use the correct Class()", func() {
			So(StackTrace().Class(), ShouldEqual, "stacktrace")
		})
	})
}
