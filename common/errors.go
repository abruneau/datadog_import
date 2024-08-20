package common

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

func Check(e error, log *logrus.Logger) {
	if e != nil {
		log.Panic(e)
		// panic(e)
	}
}

func NoStepToParseError() error {
	return fmt.Errorf("no step to parse")
}

func UnknownStepTypeError(t string) error {
	return fmt.Errorf("unknown step type: %s", t)
}

var ErrNoMoreData = errors.New("no more data")
