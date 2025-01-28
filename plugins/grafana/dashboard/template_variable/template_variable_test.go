package templatevariable

import (
	"testing"

	"datadog_import/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplateVariable(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		tv       types.TemplateVariable
		expected *datadogV1.DashboardTemplateVariable
	}{
		{
			name:   "default converter",
			source: "unknown",
			tv:     types.TemplateVariable{Name: "TestVariable"},
			expected: func() *datadogV1.DashboardTemplateVariable {
				tvName := "test_variable"
				ddtv := datadogV1.NewDashboardTemplateVariable("TestVariable")
				ddtv.SetPrefix(tvName)
				return ddtv
			}(),
		},
		{
			name:   "stackdriver converter",
			source: "stackdriver",
			tv:     types.TemplateVariable{Name: "StackVariable"},
			expected: func() *datadogV1.DashboardTemplateVariable {
				// Assuming stackdriverConverter is defined and works correctly
				return stackdriverConverter(types.TemplateVariable{Name: "StackVariable"})
			}(),
		},
		{
			name:   "azure converter",
			source: "grafana-azure-monitor-datasource",
			tv:     types.TemplateVariable{Name: "AzureVariable"},
			expected: func() *datadogV1.DashboardTemplateVariable {
				// Assuming azureConverter is defined and works correctly
				return azureConverter(types.TemplateVariable{Name: "AzureVariable"})
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTemplateVariable(tt.source, tt.tv)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestDefaultConverter(t *testing.T) {
	tests := []struct {
		name     string
		tv       types.TemplateVariable
		expected *datadogV1.DashboardTemplateVariable
	}{
		{
			name: "simple variable name",
			tv:   types.TemplateVariable{Name: "TestVariable"},
			expected: func() *datadogV1.DashboardTemplateVariable {
				tvName := "test_variable"
				ddtv := datadogV1.NewDashboardTemplateVariable("TestVariable")
				ddtv.SetPrefix(tvName)
				return ddtv
			}(),
		},
		{
			name: "variable name with spaces",
			tv:   types.TemplateVariable{Name: "Test Variable"},
			expected: func() *datadogV1.DashboardTemplateVariable {
				tvName := "test_variable"
				ddtv := datadogV1.NewDashboardTemplateVariable("Test Variable")
				ddtv.SetPrefix(tvName)
				return ddtv
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultConverter(tt.tv)
			assert.Equal(t, tt.expected, result)
		})
	}
}
