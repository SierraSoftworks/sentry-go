package sentry

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrType(t *testing.T) {
	const errType = ErrType("sentry: this is a test error")
	assert.True(t, errType.IsInstance(errType), "it should be an instance of itself")
	assert.True(t, errType.IsInstance(errors.New(errType.Error())), "errors with the same message should be an instance of this error")
	assert.EqualError(t, errType, "sentry: this is a test error", "it should report the correct error message")

	type UnwrappableError interface {
		Unwrap() error
	}

	if assert.Implements(t, (*UnwrappableError)(nil), errType, "it should implement the Unwrap() method") {
		var err error
		err = errType

		assert.Nil(t, err.(UnwrappableError).Unwrap(), "unwrapping the error should return nil")
	}
}
