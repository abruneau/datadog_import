package common

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testObject struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (t testObject) MarshalJSON() ([]byte, error) {
	// Use a type alias to avoid infinite recursion
	type testObjectAlias testObject
	return json.Marshal(testObjectAlias(t))
}

type invalidObject struct {
	Invalid func() `json:"invalid"`
}

func (i invalidObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed to marshal")
}

func TestNewFileWriter(t *testing.T) {
	// Create a temporary directory for all tests
	tempDir, err := os.MkdirTemp("", "writer_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("create new directory", func(t *testing.T) {
		testDir := filepath.Join(tempDir, "new_dir")
		writer, err := NewFileWriter(testDir)
		require.NoError(t, err)
		assert.NotNil(t, writer)
		assert.DirExists(t, testDir)
	})

	t.Run("use existing directory", func(t *testing.T) {
		testDir := filepath.Join(tempDir, "existing_dir")
		err := os.MkdirAll(testDir, 0755)
		require.NoError(t, err)

		writer, err := NewFileWriter(testDir)
		require.NoError(t, err)
		assert.NotNil(t, writer)
	})

	t.Run("with invalid path", func(t *testing.T) {
		writer, err := NewFileWriter("/invalid/path/that/should/not/exist")
		assert.Error(t, err)
		assert.Nil(t, writer)
	})
}

func TestFileWriter_Write(t *testing.T) {
	// Create a temporary directory for all tests
	tempDir, err := os.MkdirTemp("", "writer_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	writer, err := NewFileWriter(tempDir)
	require.NoError(t, err)

	ctx := context.Background()
	obj := testObject{
		Name:  "test",
		Value: 123,
	}

	t.Run("write with .json extension", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test.json")
		err := writer.Write(ctx, obj, "test.json")
		require.NoError(t, err)

		content, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), `"name": "test"`)
		assert.Contains(t, string(content), `"value": 123`)
	})

	t.Run("write without extension", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test2.json")
		err := writer.Write(ctx, obj, "test2")
		require.NoError(t, err)

		content, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), `"name": "test"`)
		assert.Contains(t, string(content), `"value": 123`)
	})

	t.Run("write with path separators", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test.json")
		err := writer.Write(ctx, obj, "path/to/test.json")
		require.NoError(t, err)

		content, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), `"name": "test"`)
		assert.Contains(t, string(content), `"value": 123`)
	})

	t.Run("write with invalid object", func(t *testing.T) {
		invalidObj := invalidObject{
			Invalid: func() {},
		}

		err := writer.Write(ctx, invalidObj, "invalid.json")
		assert.Error(t, err)
	})

	t.Run("write to read-only directory", func(t *testing.T) {
		readOnlyDir := filepath.Join(tempDir, "readonly")
		err := os.MkdirAll(readOnlyDir, 0444)
		require.NoError(t, err)

		readOnlyWriter, err := NewFileWriter(readOnlyDir)
		require.NoError(t, err)

		err = readOnlyWriter.Write(ctx, obj, "test.json")
		assert.Error(t, err)
	})
}
