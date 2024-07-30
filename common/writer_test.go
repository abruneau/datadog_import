package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileWriter(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "file_writer_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test creating a new FileWriter with a non-existing directory
	filePath := filepath.Join(tempDir, "non_existing_dir")
	fw, err := NewFileWriter(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if fw.filePath != filePath {
		t.Errorf("Expected filePath to be '%s', got '%s'", filePath, fw.filePath)
	}

	// Test creating a new FileWriter with an existing directory
	fw, err = NewFileWriter(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if fw.filePath != tempDir {
		t.Errorf("Expected filePath to be '%s', got '%s'", tempDir, fw.filePath)
	}
}
