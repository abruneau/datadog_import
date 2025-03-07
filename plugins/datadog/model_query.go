package datadog

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type Query struct {
	Metric      string
	Filters     []string
	Aggregator  datadogV1.FormulaAndFunctionMetricAggregation
	GroupBys    []string
	Function    FormulaAndFunctionMetricFunction
	Aggregation string
}

func (q *Query) Build() (string, error) {
	var query, from, by string
	if !q.Aggregator.IsValid() {
		return "", fmt.Errorf("unknown agrregator %s", q.Aggregator)
	}

	if q.Aggregation == "" {
		q.Aggregation = string(q.Aggregator)
	}

	from = "*"
	if len(q.Filters) > 0 {
		for i, v := range q.Filters {
			if i == 0 {
				from = v
			} else {
				if strings.Contains(v, " IN ") {
					from = fmt.Sprintf("%s AND %s", from, v)
				} else {
					from = fmt.Sprintf("%s,%s", from, v)
				}
			}
		}
	}

	by = ""
	if len(q.GroupBys) > 0 {
		by = fmt.Sprintf(" by {%s}", strings.Join(q.GroupBys, ","))
	}

	query = fmt.Sprintf("%s:%s{%s}%s", q.Aggregation, q.Metric, from, by)

	if q.Function != "" {
		if !q.Function.IsValid() {
			return "", fmt.Errorf("unknown function %s", q.Function)
		}
		query = fmt.Sprintf("%s.%s", query, q.Function)
	}

	return query, nil
}
