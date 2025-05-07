package common

import (
	"context"
	"datadog_import/logctx"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFiles(t *testing.T) (string, func()) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "file_reader_test")
	require.NoError(t, err)

	// Create test files
	files := []string{
		"test1.json",
		"test2.json",
		"test3.json",
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		err := os.WriteFile(filePath, []byte(`{"test": "data"}`), 0644)
		require.NoError(t, err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestNewFileReader(t *testing.T) {
	t.Run("with directory", func(t *testing.T) {
		tempDir, cleanup := setupTestFiles(t)
		defer cleanup()

		ctx := context.Background()
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		ctx = logctx.New(ctx, logrus.NewEntry(logger))

		reader, err := NewFileReader(ctx, tempDir)
		require.NoError(t, err)
		assert.NotNil(t, reader)
		assert.Equal(t, 3, len(reader.files))
		assert.Equal(t, 0, reader.index)
	})

	t.Run("with single file", func(t *testing.T) {
		tempDir, cleanup := setupTestFiles(t)
		defer cleanup()

		singleFile := filepath.Join(tempDir, "test1.json")
		ctx := context.Background()
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		ctx = logctx.New(ctx, logrus.NewEntry(logger))

		reader, err := NewFileReader(ctx, singleFile)
		require.NoError(t, err)
		assert.NotNil(t, reader)
		assert.Equal(t, 1, len(reader.files))
		assert.Equal(t, 0, reader.index)
	})

	t.Run("with non-existent path", func(t *testing.T) {
		ctx := context.Background()
		reader, err := NewFileReader(ctx, "non_existent_path")
		assert.Error(t, err)
		assert.Nil(t, reader)
	})
}

func TestFileReader_Read(t *testing.T) {
	tempDir, cleanup := setupTestFiles(t)
	defer cleanup()

	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	ctx = logctx.New(ctx, logrus.NewEntry(logger))

	reader, err := NewFileReader(ctx, tempDir)
	require.NoError(t, err)

	// Test reading all files
	for i := 0; i < 3; i++ {
		id, fileName, data, err := reader.Read(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.NotEmpty(t, fileName)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), `{"test": "data"}`)
	}

	// Test reading beyond available files
	id, fileName, data, err := reader.Read(ctx)
	assert.Equal(t, ErrNoMoreData, err)
	assert.Empty(t, id)
	assert.Empty(t, fileName)
	assert.Empty(t, data)
}
