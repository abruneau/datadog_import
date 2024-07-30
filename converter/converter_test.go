package converter

import (
	"dynatrace_to_datadog/common"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
)

type mockReader struct {
	index int
}

func (r *mockReader) Read() (string, string, []byte, error) {
	if r.index > 0 {
		r.index--
		return "id", "name", []byte("data"), nil
	} else {
		return "", "", nil, common.ErrNoMoreData
	}
}

type mockObject struct {
	Data string `json:"data"`
}

func (m *mockObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func mockTransform(data []byte) (interface {
	MarshalJSON() ([]byte, error)
}, error) {
	obj := &mockObject{
		Data: string(data),
	}
	return obj, nil
}

type mockWriter struct{}

func (w *mockWriter) Write(interface{ MarshalJSON() ([]byte, error) }, string) error {
	return nil
}

func TestConverter_Convert(t *testing.T) {
	config := &common.Config{
		Log: logrus.New(),
	}
	reader := &mockReader{
		index: 1,
	}
	writers := []Writer{&mockWriter{}}

	converter := &Converter{
		Config:    config,
		Reader:    reader,
		Transform: mockTransform,
		Writers:   writers,
	}

	converter.Convert()

	// Add assertions here to verify the expected behavior
}
