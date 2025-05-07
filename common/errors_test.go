package common

import (
	"context"
	"datadog_import/logctx"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	// Create a logger that captures output
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	ctx := logctx.New(context.Background(), logrus.NewEntry(logger))

	t.Run("no error", func(t *testing.T) {
		// This should not panic
		Check(ctx, nil)
	})

	t.Run("with error", func(t *testing.T) {
		// This should panic
		assert.Panics(t, func() {
			Check(ctx, assert.AnError)
		})
	})
}

func TestNoStepToParseError(t *testing.T) {
	err := NoStepToParseError()
	assert.Equal(t, "no step to parse", err.Error())
}

func TestUnknownStepTypeError(t *testing.T) {
	err := UnknownStepTypeError("test_type")
	assert.Equal(t, "unknown step type: test_type", err.Error())
}
