package common

import (
	"context"
	"datadog_import/logctx"
	"fmt"
	"os"
	"path"
)

// FileReader reads lines from a file
type FileReader struct {
	files []string
	index int
}

// NewFileReader creates a new FileReader
func NewFileReader(ctx context.Context, filePath string) (*FileReader, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	var files []string

	if fileInfo.IsDir() {
		files, err = GetJSONFiles(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		files = []string{filePath}
	}

	logctx.From(ctx).Debug(fmt.Sprintf("%d files found", len(files)))

	return &FileReader{
		files: files,
		index: 0,
	}, nil
}

// Read reads a line from the file
func (fr *FileReader) Read(ctx context.Context) (id, fileName string, data []byte, err error) {
	if fr.index >= len(fr.files) {
		err = ErrNoMoreData
		return // no more data
	}
	// fileName = fr.files[fr.index]
	_, fileName = path.Split(fr.files[fr.index])
	id = fileName
	logctx.From(ctx).WithField("file", fileName).Debug("reading file")
	data, err = os.ReadFile(fr.files[fr.index])

	fr.index++
	return
}
