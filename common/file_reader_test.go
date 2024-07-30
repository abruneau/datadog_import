package common

import (
	"os"
	"path/filepath"
	"testing"
)

var tempDir string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// Create a temporary directory
	var err error
	tempDir, err = os.MkdirTemp("", "file_reader_test")
	if err != nil {
		panic(err)
	}

	// Create temporary files
	file1 := filepath.Join(tempDir, "file1.json")
	file2 := filepath.Join(tempDir, "file2.json")
	err = os.WriteFile(file1, []byte("File 1 content"), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(file2, []byte("File 2 content"), 0644)
	if err != nil {
		panic(err)
	}
}

func shutdown() {
	os.RemoveAll(tempDir)
}

func TestFileReaderReadDirectory(t *testing.T) {
	// Create a new FileReader
	fr, err := NewFileReader(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Test Read method
	id, fileName, data, err := fr.Read()
	if err != nil {
		t.Fatal(err)
	}
	if id != "file1.json" {
		t.Errorf("Expected id to be 'file1.json', got '%s'", id)
	}
	if fileName != "file1.json" {
		t.Errorf("Expected fileName to be 'file1.json', got '%s'", fileName)
	}
	if string(data) != "File 1 content" {
		t.Errorf("Expected data to be 'File 1 content', got '%s'", string(data))
	}

	id, fileName, data, err = fr.Read()
	if err != nil {
		t.Fatal(err)
	}
	if id != "file2.json" {
		t.Errorf("Expected id to be 'file2.json', got '%s'", id)
	}
	if fileName != "file2.json" {
		t.Errorf("Expected fileName to be 'file2.json', got '%s'", fileName)
	}
	if string(data) != "File 2 content" {
		t.Errorf("Expected data to be 'File 2 content', got '%s'", string(data))
	}

	// Test reaching end of files
	id, fileName, data, err = fr.Read()
	if err != ErrNoMoreData {
		t.Errorf("Expected err to be ErrNoMoreData, got '%v'", err)
	}
	if id != "" {
		t.Errorf("Expected id to be empty, got '%s'", id)
	}
	if fileName != "" {
		t.Errorf("Expected fileName to be empty, got '%s'", fileName)
	}
	if data != nil {
		t.Errorf("Expected data to be nil, got '%s'", string(data))
	}
}

func TestFileReaderReadFile(t *testing.T) {
	// Create a new FileReader
	fr, err := NewFileReader(filepath.Join(tempDir, "file1.json"))
	if err != nil {
		t.Fatal(err)
	}

	// Test Read method
	id, fileName, data, err := fr.Read()
	if err != nil {
		t.Fatal(err)
	}
	if id != "file1.json" {
		t.Errorf("Expected id to be 'file1.json', got '%s'", id)
	}
	if fileName != "file1.json" {
		t.Errorf("Expected fileName to be 'file1.json', got '%s'", fileName)
	}
	if string(data) != "File 1 content" {
		t.Errorf("Expected data to be 'File 1 content', got '%s'", string(data))
	}

	// Test reaching end of files
	id, fileName, data, err = fr.Read()
	if err != ErrNoMoreData {
		t.Errorf("Expected err to be ErrNoMoreData, got '%v'", err)
	}
	if id != "" {
		t.Errorf("Expected id to be empty, got '%s'", id)
	}
	if fileName != "" {
		t.Errorf("Expected fileName to be empty, got '%s'", fileName)
	}
	if data != nil {
		t.Errorf("Expected data to be nil, got '%s'", string(data))
	}
}
