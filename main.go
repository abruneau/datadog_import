/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"dynatrace_to_datadog/cmd"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func generateConfigTemplateWithDocs(filePath string, config interface{}) {
	yamlData, err := marshalWithDocs(config, 0)
	if err != nil {
		log.Fatalf("Error marshaling config with docs to YAML: %v\n", err)
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		log.Fatalf("Error writing config file: %v\n", err)
	}

	log.Printf("Config template with docs written to %s\n", filePath)
}

func marshalWithDocs(v interface{}, indentLevel int) ([]byte, error) {
	var result string
	indent := strings.Repeat("  ", indentLevel)
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Get YAML tag
		yamlTag := fieldType.Tag.Get("yaml")
		if yamlTag == "" {
			yamlTag = fieldType.Name
		}

		// Get documentation tag
		docTag := fieldType.Tag.Get("doc")

		// Handle nested structs recursively
		if field.Kind() == reflect.Struct {
			if docTag != "" {
				result += fmt.Sprintf("\n%s# %s\n", indent, docTag)
			}
			nestedYaml, err := marshalWithDocs(field.Interface(), indentLevel+1)
			if err != nil {
				return nil, err
			}
			result += fmt.Sprintf("%s%s:\n%s\n", indent, yamlTag, string(nestedYaml))
		} else {
			if docTag != "" {
				result += fmt.Sprintf("%s# %s\n", indent, docTag)
			}
			fieldYaml, err := yaml.Marshal(map[string]interface{}{yamlTag: field.Interface()})
			if err != nil {
				return nil, err
			}
			result += indent + string(fieldYaml)
		}
	}

	return []byte(result), nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "generate-config" {
		defaultConfig := cmd.ViperConfig{
			Log: "info",
		}
		generateConfigTemplateWithDocs("config.yaml.example", defaultConfig)
		return
	}
	cmd.Execute()
}
