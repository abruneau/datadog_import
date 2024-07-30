package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetJSONFiles(t *testing.T) {

	files, err := GetJSONFiles(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		filepath.Join(tempDir, "file1.json"),
		filepath.Join(tempDir, "file2.json"),
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}

	for i, file := range files {
		if file != expectedFiles[i] {
			t.Errorf("Expected file %s, got %s", expectedFiles[i], file)
		}
	}
}

func TestCreateSubDirectories(t *testing.T) {
	rootPath, err := os.MkdirTemp("", "TestCreateSubDirectories")
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	tree := "subdir1/subdir2"

	outputDirectory, err := CreateSubDirectories(rootPath, tree)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutputDirectory := filepath.Join(rootPath, tree)

	if outputDirectory != expectedOutputDirectory {
		t.Errorf("Expected output directory %s, got %s", expectedOutputDirectory, outputDirectory)
	}

	// Check if the output directory exists
	_, err = os.Stat(outputDirectory)
	if err != nil {
		t.Errorf("Expected output directory %s to exist, got error: %v", outputDirectory, err)
	}
}

func TestWriteFileToDisk(t *testing.T) {
	rootPath, err := os.MkdirTemp("", "TestWriteFileToDisk")
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(rootPath, "file.json")
	data := []byte(`{"key": "value"}`)

	err = WriteFileToDisk(filePath, data)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the file exists
	_, err = os.Stat(filePath)
	if err != nil {
		t.Errorf("Expected file %s to exist, got error: %v", filePath, err)
	}

}
