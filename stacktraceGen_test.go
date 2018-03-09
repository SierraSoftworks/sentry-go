package sentry

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pkg/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStackTraceGenerator(t *testing.T) {
	Convey("StackTrace Generator", t, func() {
		Convey("ForError", func() {
			Convey("With .StackTrace()", func() {
				err := errors.New("test error")
				frames := getStacktraceFramesForError(err)
				So(frames.Len(), ShouldBeGreaterThan, 0)
				So(frames[frames.Len()-1].Function, ShouldEqual, "TestStackTraceGenerator.func1.1.1")
			})

			Convey("Without .StackTrace()", func() {
				err := fmt.Errorf("test error")
				frames := getStacktraceFramesForError(err)
				So(frames.Len(), ShouldEqual, 0)
			})
		})

		Convey("ForCurrentContext", func() {
			Convey("For the current function", func() {
				frames := getStacktraceFrames(0)
				So(frames.Len(), ShouldBeGreaterThan, 3)
				So(frames[frames.Len()-1].Function, ShouldEqual, "TestStackTraceGenerator.func1.2.1")
			})

			Convey("With an extreme skip", func() {
				frames := getStacktraceFrames(999999999)
				So(frames.Len(), ShouldEqual, 0)
			})
		})

		Convey("getStacktraceFrame()", func() {
			pc, file, line, ok := runtime.Caller(0)
			So(ok, ShouldBeTrue)

			frame := getStacktraceFrame(pc)
			So(frame, ShouldNotBeNil)
			So(frame.AbsoluteFilename, ShouldEqual, file)
			So(frame.Line, ShouldEqual, line)

			So(frame.Filename, ShouldEqual, "github.com/SierraSoftworks/sentry-go/stacktraceGen_test.go")
			So(frame.Function, ShouldStartWith, "TestStackTraceGenerator.func1.3")
			So(frame.Module, ShouldEqual, "sentry-go")
			So(frame.Package, ShouldEqual, "github.com/SierraSoftworks/sentry-go")
		})

		Convey("stacktraceFrame", func() {
			Convey("ClassifyInternal", func() {
				frames := getStacktraceFrames(0)
				So(frames.Len(), ShouldBeGreaterThan, 3)

				for _, frame := range frames {
					frame.ClassifyInternal([]string{"github.com/SierraSoftworks/sentry-go"})
				}

				So(frames[frames.Len()-1].InApp, ShouldBeTrue)
				So(frames[0].InApp, ShouldBeFalse)
			})
		})

		Convey("formatFuncName()", func() {
			Convey("With a full package name", func() {
				pack, module, name := formatFuncName("github.com/SierraSoftworks/sentry-go.Context")
				So(pack, ShouldEqual, "github.com/SierraSoftworks/sentry-go")
				So(module, ShouldEqual, "sentry-go")
				So(name, ShouldEqual, "Context")
			})

			Convey("With no package", func() {
				pack, module, name := formatFuncName("sentry-go.Context")
				So(pack, ShouldEqual, "sentry-go")
				So(module, ShouldEqual, "sentry-go")
				So(name, ShouldEqual, "Context")
			})
		})

		Convey("shortFilename()", func() {
			GOPATH := "/go/src"
			pkg := "github.com/SierraSoftworks/sentry-go"
			file := "stacktraceGen_test.go"
			filename := fmt.Sprintf("%s/%s/%s", GOPATH, pkg, file)

			Convey("With no package", func() {
				So(shortFilename(filename, ""), ShouldEqual, filename)
			})

			Convey("With a valid package path", func() {
				So(shortFilename(filename, pkg), ShouldEqual, fmt.Sprintf("%s/%s", pkg, file))
			})

			Convey("With an invalid package path", func() {
				So(shortFilename(filename, "cithub.com/SierraSoftworks/bender"), ShouldEqual, filename)
			})
		})
	})
}
