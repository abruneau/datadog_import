package converter

// Reader interface to read data from different sources
type Reader interface {
	Read() (string, []byte, error)
}
