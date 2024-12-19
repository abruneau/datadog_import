package widgets

import (
	"datadog_import/plugins/grafana/dashboard/widgets/converter"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type baseWidgetRequest interface {
	SetFormulas(v []datadogV1.WidgetFormula)
	SetQueries(v []datadogV1.FormulaAndFunctionQueryDefinition)
}

func (pc *PanelConvertor) newRequest(v any, aggregate bool, groupBy bool) error {
	var err error

	if pc.source == "" {
		pc.source = pc.panel.Datasource.Type
	}

	con, err := converter.NewConverter(pc.source)
	if err != nil {
		return err
	}

	q, f, err := con.Parse(pc.panel, aggregate, groupBy)
	if err != nil {
		return err
	}
	req, ok := (v).(baseWidgetRequest)
	if !ok {
		return fmt.Errorf("type assertion to baseWidgetRequest failed")
	}
	req.SetFormulas(f)
	req.SetQueries(q)

	return err

}
