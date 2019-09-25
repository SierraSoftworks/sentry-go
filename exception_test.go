package sentry

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestException(t *testing.T) {
	ex := NewExceptionInfo()
	assert.NotNil(t, ex, "the exception info should not be nil")

	e := Exception(ex)
	assert.NotNil(t, e, "the exception option should not be nil")
	assert.Implements(t, (*Option)(nil), e, "it should implement the Option interface")
	assert.Equal(t, "exception", e.Class(), "it should use the correct option class")

	exx, ok := e.(*exceptionOption)
	assert.True(t, ok, "the exception option should actually be an *exceptionOption")

	t.Run("Merge()", func(t *testing.T) {
		assert.Implements(t, (*MergeableOption)(nil), e, "it should implement the MergeableOption interface")
		
		ex2 := NewExceptionInfo()
		e2 := Exception(ex2)

		mergeable, ok := e2.(MergeableOption)
		assert.True(t, ok, "the exception option should be mergeable")

		e3 := mergeable.Merge(e)
		assert.NotNil(t, e3, "the resulting merged exception option should not be nil")
		assert.IsType(t, e, e3, "the resulting merged exception option should be the same type as the original option")

		exx, ok := e3.(*exceptionOption)
		assert.True(t, ok, "the resulting merged exception option should actually be a *exceptionOption")

		if assert.Len(t, exx.Exceptions, 2, "it should contain both exceptions") {
			assert.Equal(t, ex, exx.Exceptions[0], "the first exception should be the first exception we found")
			assert.Equal(t, ex2, exx.Exceptions[1], "the second exception should be the second exception we found")
		}

		e3 = mergeable.Merge(&testOption{})
		assert.Equal(t, e2, e3, "if the other option is not an exception option then it should be replaced")
	})

	t.Run("Finalize()", func(t *testing.T) {
		assert.Implements(t, (*FinalizeableOption)(nil), e, "it should implement the FinalizeableOption interface")

		assert.Len(t, exx.Exceptions, 1, "one exception should be registered")

		st := exx.Exceptions[0].StackTrace
		assert.NotNil(t, st, "the exception shoudl have a stacktrace")
		st.WithInternalPrefixes("github.com/SierraSoftworks/sentry-go")

		sti, ok := st.(*stackTraceOption)
		assert.True(t, ok, "the stacktrace should actually be a *stackTraceOption")
		assert.NotEmpty(t, sti.Frames, "the stacktrace should include stack frames")

		hasInternal := false
		for _, frame := range sti.Frames {
			if frame.InApp {
				hasInternal = true
			}
		}
		assert.False(t, hasInternal, "the internal stack frames should not have been processed yet")

		exx.Finalize()

		hasInternal = false
		for _, frame := range sti.Frames {
			if frame.InApp {
				hasInternal = true
			}
		}
		assert.True(t, hasInternal, "the internal stack frames should have been identified now")
	})

	t.Run("MarshalJSON()", func(t *testing.T) {
		serialized := testOptionsSerialize(t, Exception(&ExceptionInfo{
			Type:  "TestException",
			Value: "This is a test",
		}))

		assert.Equal(t, map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{
					"type":  "TestException",
					"value": "This is a test",
				},
			},
		}, serialized)
	})
}

func TestExceptionForError(t *testing.T) {
	assert.Nil(t, ExceptionForError(nil), "it should return nil if the error is nil")

	err := fmt.Errorf("example error")
	e := ExceptionForError(err)
	assert.NotNil(t, e, "it should return a non-nil option")
	assert.Implements(t, (*Option)(nil), e, "it should implement the Option interface")

	t.Run("github.com/pkg/errors", func(t *testing.T) {
		err := errors.New("root cause")
		err = errors.Wrap(err, "cause 1")
		err = errors.Wrap(err, "cause 2")
		err = errors.Wrap(err, "example error")

		e := ExceptionForError(err)
		assert.NotNil(t, e, "it should return a non-nil option")
		
		exx, ok := e.(*exceptionOption)
		assert.True(t, ok, "the option should actually be a *exceptionOption")

		// errors.Wrap adds two entries to the cause heirarchy
		// 1 - withMessage{}
		// 2 - withStack{}
		assert.Len(t, exx.Exceptions, 1 + (3*2))
		assert.Equal(t, "root cause", exx.Exceptions[0].Value)
	})
}

func TestExceptionInfo(t *testing.T) {
	t.Run("NewExceptionInfo()", func (t *testing.T) {
		ex := NewExceptionInfo()
		assert.NotNil(t, ex, "it should not return nil")
		assert.Equal(t, "unknown", ex.Type, "it should report an 'unknown' type by default")
		assert.Equal(t, "An unknown error has occurred", ex.Value, "it should report a default error message")
		assert.NotNil(t, ex.StackTrace, "it should contain a stack trace")
	})

	t.Run("ForError()", func(t *testing.T) {
		ex := NewExceptionInfo()
		assert.NotNil(t, ex, "it should not return nil")

		assert.Equal(t, ex, ex.ForError(fmt.Errorf("example error")), "it should return the same exception info object for chaining")
		assert.Equal(t, "example error", ex.Type, "it should load the type from the error")
		assert.Equal(t, "example error", ex.Value, "it should load the message from the error")
		assert.Equal(t, "", ex.Module, "it should load the module from the error")

		t.Run("with no stacktrace", func(t *testing.T) {
			ex := &ExceptionInfo{}
			ex.ForError(fmt.Errorf("example error"))
			assert.NotNil(t, ex.StackTrace, "it should use the location of the current call as the stack trace")
		})

		t.Run("with a fmt.Errorf() error", func(t *testing.T) {
			assert.NotNil(t, ex.StackTrace, "it should use the location of the current call as the stack trace")
		})

		t.Run("with a github.com/pkg/errors error", func(t *testing.T) {
			ex.ForError(errors.New("example error"))
			assert.NotNil(t, ex.StackTrace, "it should use the location of the error as the stack trace")
		})

		t.Run("with a structured error message", func(t *testing.T) {
			ex.ForError(fmt.Errorf("test: example error"))
			assert.Equal(t, "test: example error", ex.Value, "it should load the message from the error")
			assert.Equal(t, "example error", ex.Type, "it should load the type from the error")
			assert.Equal(t, "test", ex.Module, "it should load the module from the error")
		})
	})
}
