package sentry

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestException(t *testing.T) {
	Convey("Exception", t, func() {
		Convey("Exception()", func() {
			ex := NewExceptionInfo()
			Convey("Should return an Option", func() {
				So(Exception(ex), ShouldImplement, (*Option)(nil))
			})

			Convey("Should use the correct Class()", func() {
				So(Exception(ex).Class(), ShouldEqual, "exception")
			})

			Convey("Merge()", func() {
				Convey("Should append newer exceptions to the list", func() {
					exNew := NewExceptionInfo()

					exo1 := Exception(ex)
					exo2 := Exception(exNew)

					mergable, ok := exo2.(MergeableOption)
					So(ok, ShouldBeTrue)
					exo3 := mergable.Merge(exo1)
					So(exo3, ShouldNotBeNil)
					So(exo3, ShouldHaveSameTypeAs, exo1)

					exx, ok := exo3.(*exceptionOption)
					So(ok, ShouldBeTrue)
					So(exx.Exceptions, ShouldHaveLength, 2)
					So(exx.Exceptions[0], ShouldEqual, ex)
					So(exx.Exceptions[1], ShouldEqual, exNew)
				})

				Convey("Should overwrite if it doesn't recognize the old option", func() {
					exo := Exception(ex)
					mergable, ok := exo.(MergeableOption)
					So(ok, ShouldBeTrue)

					So(mergable.Merge(&testOption{}), ShouldEqual, exo)
				})
			})

			Convey("Finalize()", func() {
				Convey("Should call finalize on all of its exception's stacktraces", func() {
					err := errors.New("test error")
					So(err, ShouldNotBeNil)

					ex := ExceptionForError(err)
					So(ex, ShouldNotBeNil)

					exx, ok := ex.(*exceptionOption)
					So(ok, ShouldBeTrue)

					So(exx.Exceptions, ShouldHaveLength, 1)

					st := exx.Exceptions[0].StackTrace
					So(st, ShouldNotBeNil)
					st.WithInternalPrefixes("github.com/SierraSoftworks/sentry-go")

					sti, ok := st.(*stackTraceOption)
					So(ok, ShouldBeTrue)
					So(sti.Frames, ShouldNotBeEmpty)

					hasInternal := false
					for _, frame := range sti.Frames {
						if frame.InApp {
							hasInternal = true
						}
					}
					So(hasInternal, ShouldBeFalse)

					exx.Finalize()

					hasInternal = false
					for _, frame := range sti.Frames {
						if frame.InApp {
							hasInternal = true
						}
					}
					So(hasInternal, ShouldBeTrue)
				})
			})
		})

		Convey("ExceptionForError()", func() {
			Convey("Should return an Option", func() {
				err := fmt.Errorf("example error")
				So(ExceptionForError(err), ShouldImplement, (*Option)(nil))
			})

			Convey("With wrapped errors", func() {
				err := errors.New("root cause")
				err = errors.Wrap(err, "cause 1")
				err = errors.Wrap(err, "cause 2")
				err = errors.Wrap(err, "example error")

				ex := ExceptionForError(err)
				So(ex, ShouldNotBeNil)

				exx, ok := ex.(*exceptionOption)
				So(ok, ShouldBeTrue)

				// errors.Wrap adds two entries to the cause heirarchy
				// 1 - withMessage{}
				// 2 - withStack{}
				So(exx.Exceptions, ShouldHaveLength, 1+(3*2))
				So(exx.Exceptions[0].Value, ShouldEqual, "root cause")
			})
		})

		Convey("ExceptionInfo", func() {
			Convey("NewExceptionInfo()", func() {
				ex := NewExceptionInfo()
				So(ex, ShouldNotBeNil)
				So(ex.Type, ShouldEqual, "unknown")
				So(ex.Value, ShouldEqual, "An unknown error has occurred")
				So(ex.StackTrace, ShouldNotBeNil)
			})

			Convey("ForError()", func() {
				ex := NewExceptionInfo()
				So(ex, ShouldNotBeNil)

				Convey("Without an existing StackTrace", func() {
					ex := &ExceptionInfo{}
					err := fmt.Errorf("example error")
					So(ex.ForError(err), ShouldEqual, ex)
					So(ex.Type, ShouldEqual, "example error")
					So(ex.Module, ShouldEqual, "")
					So(ex.Value, ShouldEqual, "example error")
					So(ex.StackTrace, ShouldNotBeNil)
				})

				Convey("With a normal error", func() {
					err := fmt.Errorf("example error")
					So(ex.ForError(err), ShouldEqual, ex)
					So(ex.Type, ShouldEqual, "example error")
					So(ex.Module, ShouldEqual, "")
					So(ex.Value, ShouldEqual, "example error")
					So(ex.StackTrace, ShouldNotBeNil)
				})

				Convey("With a stacktraceable error", func() {
					err := errors.New("example error")
					So(ex.ForError(err), ShouldEqual, ex)
					So(ex.Module, ShouldEqual, "")
					So(ex.Type, ShouldEqual, "example error")
					So(ex.Value, ShouldEqual, "example error")
					So(ex.StackTrace, ShouldNotBeNil)
				})

				Convey("With a well formatted message", func() {
					err := errors.New("test: example error")
					So(ex.ForError(err), ShouldEqual, ex)
					So(ex.Module, ShouldEqual, "test")
					So(ex.Type, ShouldEqual, "example error")
					So(ex.Value, ShouldEqual, "test: example error")
				})
			})
		})

		Convey("MarshalJSON", func() {
			Convey("Should marshal correctly", func() {
				serialized := testOptionsSerialize(Exception(&ExceptionInfo{
					Type:  "TestException",
					Value: "This is a test",
				}))
				So(serialized, ShouldResemble, map[string]interface{}{
					"values": []interface{}{
						map[string]interface{}{
							"type":  "TestException",
							"value": "This is a test",
						},
					},
				})
			})
		})
	})
}
