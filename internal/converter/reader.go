package converter

import "context"

// Reader interface to read data from different sources
type Reader interface {
	Read(ctx context.Context) (id, name string, data []byte, err error)
}
