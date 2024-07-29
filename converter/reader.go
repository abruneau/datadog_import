package converter

// Reader interface to read data from different sources
type Reader interface {
	Read() (id, name string, data []byte, err error)
}
