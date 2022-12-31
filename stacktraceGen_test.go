package sentry

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackTraceGenerator(t *testing.T) {
	t.Run("getStacktraceFramesForError()", func(t *testing.T) {
		t.Run("StackTraceableError", func(t *testing.T) {
			err := errors.New("test error")
			frames := getStacktraceFramesForError(err)
			if assert.NotEmpty(t, frames, "there should be frames from the error") {
				assert.Equal(t, "TestStackTraceGenerator.func1.1", frames[frames.Len()-1].Function, "it should have the right function name as the top-most frame")
			}
		})

		t.Run("error", func(t *testing.T) {
			err := fmt.Errorf("test error")
			frames := getStacktraceFramesForError(err)
			assert.Empty(t, frames, "there should be no frames from a normal error")
		})
	})

	t.Run("getStacktraceFrames()", func(t *testing.T) {
		t.Run("Skip", func(t *testing.T) {
			frames := getStacktraceFrames(999999999)
			assert.Empty(t, frames, "with an extreme skip, there should be no frames")
		})

		t.Run("Current Function", func(t *testing.T) {
			frames := getStacktraceFrames(0)
			if assert.NotEmpty(t, frames, "there should be frames from the current function") {
				assert.Equal(t, "TestStackTraceGenerator.func2.2", frames[frames.Len()-1].Function, "it should have the right function name as the top-most frame")
			}
		})
	})

	t.Run("getStackTraceFrame()", func(t *testing.T) {
		pc, file, line, ok := runtime.Caller(0)
		require.True(t, ok, "we should be able to get the current caller")

		frame := getStacktraceFrame(pc)
		require.NotNil(t, frame, "the frame should not be nil")

		assert.Equal(t, file, frame.AbsoluteFilename, "the filename for the frame should match the caller")
		assert.Equal(t, line, frame.Line, "the line from the frame should match the caller")

		assert.Regexp(t, ".*/sentry-go/stacktraceGen_test.go$", frame.Filename, "it should have the correct filename")
		assert.Equal(t, "TestStackTraceGenerator.func3", frame.Function, "it should have the correct function name")
		assert.Equal(t, "sentry-go/v2", frame.Module, "it should have the correct module name")
		assert.Equal(t, "github.com/SierraSoftworks/sentry-go/v2", frame.Package, "it should have the correct package name")
	})

	t.Run("stackTraceFrame.ClassifyInternal()", func(t *testing.T) {
		frames := getStacktraceFrames(0)
		require.Greater(t, frames.Len(), 3, "the number of frames should be more than 3")

		for i, frame := range frames {
			assert.False(t, frame.InApp, "all frames should initially be marked as external (frame index = %d)", i)
			frame.ClassifyInternal([]string{"github.com/SierraSoftworks/sentry-go"})
		}

		assert.True(t, frames[frames.Len()-1].InApp, "the top-most frame should be marked as internal (this function)")
		assert.False(t, frames[0].InApp, "the bottom-most frame should be marked as external (the test harness main method)")
	})

	t.Run("formatFuncName()", func(t *testing.T) {
		cases := []struct {
			Name string

			FullName     string
			Package      string
			Module       string
			FunctionName string
		}{
			{"Full Name", "github.com/SierraSoftworks/sentry-go.Context", "github.com/SierraSoftworks/sentry-go", "sentry-go", "Context"},
			{"Full Name (v2)", "github.com/SierraSoftworks/sentry-go/v2.Context", "github.com/SierraSoftworks/sentry-go/v2", "sentry-go/v2", "Context"},
			{"Struct Function Name", "github.com/SierraSoftworks/sentry-go.packet.Clone", "github.com/SierraSoftworks/sentry-go", "sentry-go", "packet.Clone"},
			{"Struct Function Name (v2)", "github.com/SierraSoftworks/sentry-go/v2.packet.Clone", "github.com/SierraSoftworks/sentry-go/v2", "sentry-go/v2", "packet.Clone"},
			{"No Package", "sentry-go.Context", "sentry-go", "sentry-go", "Context"},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				pack, module, name := formatFuncName(tc.FullName)
				assert.Equal(t, tc.Package, pack, "the package name should be correct")
				assert.Equal(t, tc.Module, module, "the module name should be correct")
				assert.Equal(t, tc.FunctionName, name, "the function name should be correct")
			})
		}
	})

	t.Run("shortFilename()", func(t *testing.T) {
		t.Run("GOPATH", func(t *testing.T) {
			GOPATH := "/go/src"
			pkg := "github.com/SierraSoftworks/sentry-go"
			file := "stacktraceGen_test.go"
			filename := fmt.Sprintf("%s/%s/%s", GOPATH, pkg, file)

			assert.Equal(t, filename, shortFilename(filename, ""), "should use the original filename if no package is provided")
			assert.Equal(t, filename, shortFilename(filename, "bitblob.com/bender"), "should use the original filename if the package name doesn't match the path")
			assert.Equal(t, fmt.Sprintf("%s/%s", pkg, file), shortFilename(filename, pkg), "should use the $pkg/$file if the package is provided")
		})
	})
}
