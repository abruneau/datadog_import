package widgets

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

var displayMap = map[string]datadogV1.WidgetDisplayType{
	"line":   datadogV1.WIDGETDISPLAYTYPE_LINE,
	"bar":    datadogV1.WIDGETDISPLAYTYPE_BARS,
	"points": datadogV1.WIDGETDISPLAYTYPE_LINE,
}

func (pc *PanelConvertor) newTimeseriesDefinition() (datadogV1.WidgetDefinition, error) {
	var widgetRequest = datadogV1.NewTimeseriesWidgetRequest()
	err := pc.newRequest(widgetRequest, false, true)

	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}

	if pc.panel.Type == "barchart" {
		widgetRequest.SetDisplayType(datadogV1.WIDGETDISPLAYTYPE_BARS)
	} else if pc.panel.FieldConfig.Defaults.Custom.DrawStyle != "" {
		widgetRequest.SetDisplayType(displayMap[pc.panel.FieldConfig.Defaults.Custom.DrawStyle])
	}
	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_TIMESERIES)

	tsDefinition := datadogV1.NewTimeseriesWidgetDefinition([]datadogV1.TimeseriesWidgetRequest{*widgetRequest}, datadogV1.TIMESERIESWIDGETDEFINITIONTYPE_TIMESERIES)
	tsDefinition.SetTitle(pc.panel.Title)
	tsDefinition.SetTitleSize("16")

	return datadogV1.TimeseriesWidgetDefinitionAsWidgetDefinition(tsDefinition), nil
}
