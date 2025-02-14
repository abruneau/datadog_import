package templatevariable

import (
	"testing"

	"datadog_import/internal/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/stretchr/testify/assert"
)

func TestStackdriverConverter(t *testing.T) {
	tests := []struct {
		name     string
		tv       types.TemplateVariable
		expected *datadogV1.DashboardTemplateVariable
	}{
		{
			name: "alignmentPeriod variable",
			tv: types.TemplateVariable{
				Name: "alignmentPeriod",
			},
			expected: nil,
		},
		{
			name: "other variable",
			tv: types.TemplateVariable{
				Name: "otherVariable",
			},
			expected: &datadogV1.DashboardTemplateVariable{
				Name: "otherVariable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stackdriverConverter(tt.tv)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Name, result.GetName())
			}
		})
	}
}
