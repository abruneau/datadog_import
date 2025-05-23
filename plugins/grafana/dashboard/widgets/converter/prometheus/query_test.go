package prometheus

import (
	"datadog_import/plugins/grafana/dashboard/widgets/shared"
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	expr    string
	refId   string
	request shared.Request
}{
	{
		expr:  "sum(rate(foo{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (tenant) / 1024 / 1024",
		refId: "A",
		request: shared.Request{
			Formulas: []string{"A20 / 1024 / 1024"},
			Queries: []struct {
				Name        string
				Query       string
				Aggregation datadogV1.FormulaAndFunctionMetricAggregation
			}{
				{
					Name:        "A20",
					Query:       "sum:foo{cluster:$cluster.value,namespace:$namespace.value} by {tenant}.as_rate()",
					Aggregation: datadogV1.FORMULAANDFUNCTIONMETRICAGGREGATION_SUM,
				},
			},
		},
	},
	{
		expr:  "clamp_min(foo - foo offset 60s, 0)",
		refId: "A",
		request: shared.Request{
			Formulas: []string{"clamp_min(A20 - hour_before(A30), 0)"},
			Queries: []struct {
				Name        string
				Query       string
				Aggregation datadogV1.FormulaAndFunctionMetricAggregation
			}{
				{
					Name:        "A20",
					Query:       "avg:foo{*}.as_count()",
					Aggregation: datadogV1.FORMULAANDFUNCTIONMETRICAGGREGATION_AVG,
				},
				{
					Name:        "A30",
					Query:       "avg:foo{*}.as_count()",
					Aggregation: datadogV1.FORMULAANDFUNCTIONMETRICAGGREGATION_AVG,
				},
			},
		},
	},
	{
		expr:  "histogram_quantile(.99, sum(rate(foo{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (operation, le))",
		refId: "A",
		request: shared.Request{
			Formulas: []string{"A"},
			Queries: []struct {
				Name        string
				Query       string
				Aggregation datadogV1.FormulaAndFunctionMetricAggregation
			}{
				{
					Name:        "A",
					Query:       "p99:foo{cluster:$cluster.value,namespace:$namespace.value} by {operation,le}.as_rate()",
					Aggregation: datadogV1.FORMULAANDFUNCTIONMETRICAGGREGATION_PERCENTILE,
				},
			},
		},
	},
	{
		expr:  "irate(jvm_classes_unloaded_classes_total{kubernetes_namespace=\"$environment\", app=\"$application\", kubernetes_pod_name=~\"$instance\"}[5m])",
		refId: "A",
		request: shared.Request{
			Formulas: []string{"A"},
			Queries: []struct {
				Name        string
				Query       string
				Aggregation datadogV1.FormulaAndFunctionMetricAggregation
			}{
				{
					Name:        "A",
					Query:       "avg:jvm_classes_unloaded_classes_total{kubernetes_namespace:$environment.value,app:$application.value,kubernetes_pod_name:$instance.value}.as_rate()",
					Aggregation: datadogV1.FORMULAANDFUNCTIONMETRICAGGREGATION_AVG,
				},
			},
		},
	},
}

func TestNewQuery(t *testing.T) {
	for _, test := range tests {
		target := map[string]interface{}{
			"expr":  test.expr,
			"refId": test.refId,
		}

		r, err := NewQuery(target, true)
		assert.Nil(t, err)
		assert.Equal(t, test.request, r)
	}
}
