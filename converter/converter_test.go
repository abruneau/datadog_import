package converter

import (
	"context"
	"datadog_import/common"
	"datadog_import/logctx"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// mockReader implements the Reader interface for testing
type mockReader struct {
	data    [][]byte
	current int
}

func (m *mockReader) Read(ctx context.Context) (id, name string, data []byte, err error) {
	if m.current >= len(m.data) {
		return "", "", nil, common.ErrNoMoreData
	}
	data = m.data[m.current]
	m.current++
	return "test-id", "test-name", data, nil
}

// mockTransformer implements the Transformer interface for testing
type mockTransformer struct {
	transform func(data []byte) (interface {
		MarshalJSON() ([]byte, error)
	}, error)
}

func (m *mockTransformer) Transform(ctx context.Context, data []byte) (interface {
	MarshalJSON() ([]byte, error)
}, error) {
	return m.transform(data)
}

// mockWriter implements the Writer interface for testing
type mockWriter struct {
	written []interface{}
}

func (m *mockWriter) Write(ctx context.Context, obj interface {
	MarshalJSON() ([]byte, error)
}, name string) error {
	m.written = append(m.written, obj)
	return nil
}

type testObject struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (t testObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func TestConverter_Convert(t *testing.T) {
	// Create a test context with a logger
	ctx := logctx.New(context.Background(), logrus.NewEntry(logrus.New()))

	tests := []struct {
		name          string
		readerData    [][]byte
		transformFunc func([]byte) (interface {
			MarshalJSON() ([]byte, error)
		}, error)
		expectedWrites int
		expectError    bool
	}{
		{
			name: "successful conversion",
			readerData: [][]byte{
				[]byte(`{"name": "test1", "value": 1}`),
				[]byte(`{"name": "test2", "value": 2}`),
			},
			transformFunc: func(data []byte) (interface {
				MarshalJSON() ([]byte, error)
			}, error) {
				var obj testObject
				if err := json.Unmarshal(data, &obj); err != nil {
					return nil, err
				}
				return obj, nil
			},
			expectedWrites: 2,
			expectError:    false,
		},
		{
			name: "transform error",
			readerData: [][]byte{
				[]byte(`{"name": "test1", "value": 1}`),
				[]byte(`invalid json`),
			},
			transformFunc: func(data []byte) (interface {
				MarshalJSON() ([]byte, error)
			}, error) {
				var obj testObject
				if err := json.Unmarshal(data, &obj); err != nil {
					return nil, err
				}
				return obj, nil
			},
			expectedWrites: 1,
			expectError:    false,
		},
		{
			name:       "empty reader",
			readerData: [][]byte{},
			transformFunc: func(data []byte) (interface {
				MarshalJSON() ([]byte, error)
			}, error) {
				var obj testObject
				if err := json.Unmarshal(data, &obj); err != nil {
					return nil, err
				}
				return obj, nil
			},
			expectedWrites: 0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock components
			reader := &mockReader{data: tt.readerData}
			transformer := &mockTransformer{transform: tt.transformFunc}
			writer := &mockWriter{}

			// Create converter
			converter := NewConverter(ctx, reader, transformer.Transform, []Writer{writer})

			// Run conversion
			converter.Convert()

			// Verify results
			if tt.expectError {
				assert.NotEmpty(t, writer.written)
			} else {
				assert.Equal(t, tt.expectedWrites, len(writer.written))
			}
		})
	}
}

func TestNewConverter(t *testing.T) {
	ctx := context.Background()
	reader := &mockReader{}
	transformer := &mockTransformer{}
	writer := &mockWriter{}

	// Test successful creation
	conv := NewConverter(ctx, reader, transformer.Transform, []Writer{writer})
	assert.NotNil(t, conv)

	// Test with nil reader - should not panic
	conv = NewConverter(ctx, nil, transformer.Transform, []Writer{writer})
	assert.NotNil(t, conv)

	// Test with nil writers - should not panic
	conv = NewConverter(ctx, reader, transformer.Transform, nil)
	assert.NotNil(t, conv)

	// Test with empty writers slice - should not panic
	conv = NewConverter(ctx, reader, transformer.Transform, []Writer{})
	assert.NotNil(t, conv)
}
