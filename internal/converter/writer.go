package converter

import "context"

// Writer interface to write data to different destinations
type Writer interface {
	Write(ctx context.Context, obj interface {
		MarshalJSON() ([]byte, error)
	}, name string) error
}
