package cmd

import (
	"dynatrace_to_datadog/datadog/api"
	"dynatrace_to_datadog/dynatrace"
)

type ViperConfig struct {
	Log       string           `mapstructure:"log" doc:"Log level"`
	Dynatrace dynatrace.Config `mapstructure:"dynatrace" doc:"Dynatrace Configuration"`
	Datadog   api.Config       `mapstructure:"datadog" doc:"Datadog Configuration"`
	Output    string           `mapstructure:"output" doc:"Output Directory"`
}
