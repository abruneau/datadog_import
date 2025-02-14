package templatevariable

import (
	"datadog_import/internal/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func azureConverter(tv types.TemplateVariable) *datadogV1.DashboardTemplateVariable {
	if tv.Name == "ResourceName" {
		ddtv := datadogV1.NewDashboardTemplateVariable(tv.Name)
		ddtv.SetPrefix("name")
		return ddtv
	}
	return defaultConverter(tv)
}
