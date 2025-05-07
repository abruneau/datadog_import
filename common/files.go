package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GetJSONFiles returns all JSON files in a directory and subdirectories.
func GetJSONFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func CreateSubDirectories(rootPath, tree string) (string, error) {
	fStructure := strings.Split(tree, "/")

	fStructure = append([]string{rootPath}, fStructure...)
	outputDirectory := path.Join(fStructure...)

	_, err := os.Stat(outputDirectory)
	if errors.Is(err, os.ErrNotExist) {
		return outputDirectory, os.MkdirAll(outputDirectory, os.ModePerm)
	}
	return outputDirectory, err

}

func WriteFileToDisk(filePath string, data []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, prettyJSON.Bytes(), 0644)
}
