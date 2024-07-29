package common

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// FileWriter writes data to a file
type FileWriter struct {
	filePath string
}

// NewFileWriter creates a new FileWriter
func NewFileWriter(filePath string) (*FileWriter, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &FileWriter{filePath: filePath}, nil
}

// Write writes data to the file
func (fw *FileWriter) Write(obj interface {
	MarshalJSON() ([]byte, error)
}, fileName string) error {
	res, err := obj.MarshalJSON()
	if err != nil {
		return err
	}
	var output string

	fileName = strings.ReplaceAll(fileName, "/", "")

	if strings.HasSuffix(fileName, ".json") {
		output = path.Join(fw.filePath, fileName)
	} else {
		output = path.Join(fw.filePath, fmt.Sprintf("%s.json", fileName))
	}

	return WriteFileToDisk(output, res)
}
