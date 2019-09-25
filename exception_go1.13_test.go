// +build go1.13

package sentry

import (
	"testing"
	"fmt"
	
	"github.com/stretchr/testify/assert"
)

func TestExceptionForErrorWrappingGo113(t *testing.T) {
	t.Run("fmt.Errorf()", func(t *testing.T) {
		err := fmt.Errorf("root cause")
		err = fmt.Errorf("cause 1: %w", err)
		err = fmt.Errorf("cause 2: %w", err)
		err = fmt.Errorf("example error: %w", err)

		e := ExceptionForError(err)
		assert.NotNil(t, e, "it should return a non-nil option")
		
		exx, ok := e.(*exceptionOption)
		assert.True(t, ok, "the option should actually be a *exceptionOption")

		assert.Len(t, exx.Exceptions, 4)
		assert.Equal(t, "root cause", exx.Exceptions[0].Value)
	})
}