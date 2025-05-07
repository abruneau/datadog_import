package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDirectory(t *testing.T) (string, func()) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "files_test")
	require.NoError(t, err)

	// Create test files and directories
	files := []string{
		"test1.json",
		"test2.json",
		"test3.txt",
		"subdir/test4.json",
		"subdir/test5.txt",
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		// Create parent directory if it doesn't exist
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		require.NoError(t, err)
		// Create the file
		err = os.WriteFile(filePath, []byte(`{"test": "data"}`), 0644)
		require.NoError(t, err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestGetJSONFiles(t *testing.T) {
	tempDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	files, err := GetJSONFiles(tempDir)
	require.NoError(t, err)
	assert.Equal(t, 3, len(files), "Should find 3 JSON files")

	// Verify all returned files are JSON files
	for _, file := range files {
		assert.Equal(t, ".json", filepath.Ext(file))
	}

	// Verify specific files are included
	expectedFiles := []string{
		filepath.Join(tempDir, "test1.json"),
		filepath.Join(tempDir, "test2.json"),
		filepath.Join(tempDir, "subdir", "test4.json"),
	}

	for _, expected := range expectedFiles {
		assert.Contains(t, files, expected)
	}
}

func TestCreateSubDirectories(t *testing.T) {
	tempDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	t.Run("create new directory", func(t *testing.T) {
		dir, err := CreateSubDirectories(tempDir, "new/sub/dir")
		require.NoError(t, err)
		assert.DirExists(t, dir)
	})

	t.Run("create existing directory", func(t *testing.T) {
		// First create the directory
		dir, err := CreateSubDirectories(tempDir, "existing/dir")
		require.NoError(t, err)
		assert.DirExists(t, dir)

		// Try to create it again
		dir2, err := CreateSubDirectories(tempDir, "existing/dir")
		require.NoError(t, err)
		assert.Equal(t, dir, dir2)
	})
}

func TestWriteFileToDisk(t *testing.T) {
	tempDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	t.Run("write valid JSON", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "output.json")
		data := []byte(`{"test": "data", "nested": {"value": 123}}`)

		err := WriteFileToDisk(filePath, data)
		require.NoError(t, err)

		// Read the file back and verify it was pretty-printed
		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "\t")
		assert.Contains(t, string(content), "test")
		assert.Contains(t, string(content), "data")
	})

	t.Run("write invalid JSON", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "invalid.json")
		data := []byte(`{"test": "data", invalid json}`)

		err := WriteFileToDisk(filePath, data)
		assert.Error(t, err)
	})
}
