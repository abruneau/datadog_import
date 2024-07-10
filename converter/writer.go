package converter

// Writer interface to write data to different destinations
type Writer interface {
	Write(obj interface {
		MarshalJSON() ([]byte, error)
	}, name string) error
}
