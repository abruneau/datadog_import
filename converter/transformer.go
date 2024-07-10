package converter

type Transformer func(data []byte) (interface {
	MarshalJSON() ([]byte, error)
}, error)
