package widgets

import (
	"datadog_import/plugins/grafana/dashboard/types"
	"datadog_import/plugins/grafana/dashboard/widgets/converter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func newPiechartDefinition(source string, panel types.Panel) (datadogV1.WidgetDefinition, error) {
	request, err := newPiechartRequest(source, panel)
	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	def := datadogV1.NewSunburstWidgetDefinition(request, datadogV1.SUNBURSTWIDGETDEFINITIONTYPE_SUNBURST)
	def.SetTitle(panel.Title)
	def.SetTitleSize("16")

	return datadogV1.SunburstWidgetDefinitionAsWidgetDefinition(def), nil
}

func newPiechartRequest(source string, panel types.Panel) ([]datadogV1.SunburstWidgetRequest, error) {
	var widgetRequest *datadogV1.SunburstWidgetRequest
	var err error

	if source == "" {
		source = panel.Datasource.Type
	}

	con, err := converter.NewConverter(source)
	if err != nil {
		return nil, err
	}

	widgetRequest, err = con.NewSunburstWidgetRequest(panel)

	if err != nil {
		return nil, err
	}

	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)

	return []datadogV1.SunburstWidgetRequest{*widgetRequest}, nil
}
