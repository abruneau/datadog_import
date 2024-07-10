package setting

import (
	"fmt"
	"net/url"
)

func Path(s string) func(args ...string) string {
	return path(s).Complete
}

type path string

func (p path) Complete(args ...string) string {
	if len(args) == 0 {
		return string(p)
	}
	fmtargs := make([]any, len(args))
	for idx, s := range args {
		fmtargs[idx] = url.PathEscape(s)
	}
	return fmt.Sprintf(string(p), fmtargs...)
}
