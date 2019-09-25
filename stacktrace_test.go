package sentry

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAddInternalPrefixes(t *testing.T) {
	assert.Contains(t, defaultInternalPrefixes, "main")
	AddInternalPrefixes("github.com/SierraSoftworks/sentry-go")
	assert.Contains(t, defaultInternalPrefixes, "github.com/SierraSoftworks/sentry-go")
}

func TestStackTrace(t *testing.T) {
	o := StackTrace()
	require.NotNil(t, o, "it should return a non-nil option")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "stacktrace", o.Class(), "it should use the right option class")

	sti, ok := o.(*stackTraceOption)
	require.True(t, ok, "it should actually be a *stackTraceOption")

	assert.NotEmpty(t, sti.Frames, "it should start off with your current stack frames")
	originalFrames := sti.Frames

	err := errors.New("example error")
	assert.Same(t, o, o.ForError(err), "it should reuse the same instance when adding error information")
	assert.NotEmpty(t, sti.Frames, "it should have loaded frame information from the error")
	assert.NotEqual(t, originalFrames, sti.Frames, "the frames should not be the original ones it started with")

	assert.Equal(t, defaultInternalPrefixes, sti.internalPrefixes, "it should start out with the default internal prefixes")

	o.WithInternalPrefixes("github.com/SierraSoftworks")
	assert.Contains(t, sti.internalPrefixes, "github.com/SierraSoftworks", "it should allow you to add new internal prefixes")

	if assert.Implements(t, (*FinalizeableOption)(nil), o, "it should implement the FinalizeableOption interface") {
		for i, frame := range sti.Frames {
			assert.False(t, frame.InApp, "all frames should initially be marked as external (frame index=%d)", i)
		}

		sti.Finalize()

		if assert.NotEmpty(t, sti.Frames, "the frames list should not be empty") {
			assert.True(t, sti.Frames[len(sti.Frames)-1].InApp, "the final frame should be marked as internal")
		}
	}
}
