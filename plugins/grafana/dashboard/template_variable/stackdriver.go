package templatevariable

import (
	"datadog_import/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func stackdriverConverter(tv types.TemplateVariable) *datadogV1.DashboardTemplateVariable {
	if tv.Name == "alignmentPeriod" {
		return nil
	}
	return defaultConverter(tv)
}
