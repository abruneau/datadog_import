package common

import (
	"context"
	"datadog_import/logctx"
	"errors"
	"fmt"
)

func Check(ctx context.Context, e error) {
	if e != nil {
		logctx.From(ctx).Panic(e)
	}
}

func NoStepToParseError() error {
	return fmt.Errorf("no step to parse")
}

func UnknownStepTypeError(t string) error {
	return fmt.Errorf("unknown step type: %s", t)
}

var ErrNoMoreData = errors.New("no more data")
