package cmd

import "dynatrace_to_datadog/dynatrace"

type ViperConfig struct {
	Log       string           `yaml:"log" doc:"Log level"`
	Dynatrace dynatrace.Config `yaml:"dynatrace" doc:"Dynatrace Configuration"`
	Output    string           `yaml:"output" doc:"Output Directory"`
}
