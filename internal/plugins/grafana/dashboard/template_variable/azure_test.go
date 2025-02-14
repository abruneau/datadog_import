package templatevariable

import (
	"testing"

	"datadog_import/internal/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/stretchr/testify/assert"
)

func TestAzureConverter(t *testing.T) {
	tests := []struct {
		name     string
		input    types.TemplateVariable
		expected *datadogV1.DashboardTemplateVariable
	}{
		{
			name: "ResourceName variable",
			input: types.TemplateVariable{
				Name: "ResourceName",
			},
			expected: func() *datadogV1.DashboardTemplateVariable {
				ddtv := datadogV1.NewDashboardTemplateVariable("ResourceName")
				ddtv.SetPrefix("name")
				return ddtv
			}(),
		},
		{
			name: "Other variable",
			input: types.TemplateVariable{
				Name: "OtherName",
			},
			expected: defaultConverter(types.TemplateVariable{Name: "OtherName"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := azureConverter(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
