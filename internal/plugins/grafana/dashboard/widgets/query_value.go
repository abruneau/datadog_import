package widgets

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func (pc *PanelConvertor) newQueryValueDefinition() (datadogV1.WidgetDefinition, error) {
	var widgetRequest = datadogV1.NewQueryValueWidgetRequest()
	err := pc.newRequest(widgetRequest, true, false)

	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)
	qvDefinition := datadogV1.NewQueryValueWidgetDefinition([]datadogV1.QueryValueWidgetRequest{*widgetRequest}, datadogV1.QUERYVALUEWIDGETDEFINITIONTYPE_QUERY_VALUE)
	qvDefinition.SetTitle(pc.panel.Title)
	qvDefinition.SetTitleSize("16")
	qvDefinition.SetAutoscale(true)
	qvDefinition.SetPrecision(2)

	if pc.panel.Options.GraphMode == "area" {
		qvDefinition.TimeseriesBackground = datadogV1.NewTimeseriesBackground(datadogV1.TIMESERIESBACKGROUNDTYPE_AREA)
	}

	return datadogV1.QueryValueWidgetDefinitionAsWidgetDefinition(qvDefinition), nil
}
