package common

import (
	"errors"
	"fmt"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func NoStepToParseError() error {
	return fmt.Errorf("no step to parse")
}

func UnknownStepTypeError(t string) error {
	return fmt.Errorf("unknown step type: %s", t)
}

var ErrNoMoreData = errors.New("no more data")
