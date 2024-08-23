package converter

import "context"

type Transformer func(ctx context.Context, data []byte) (interface {
	MarshalJSON() ([]byte, error)
}, error)
