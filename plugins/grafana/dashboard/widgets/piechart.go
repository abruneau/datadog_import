package widgets

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func (pc *PanelConvertor) newPiechartDefinition() (datadogV1.WidgetDefinition, error) {
	var widgetRequest = datadogV1.NewSunburstWidgetRequest()

	err := pc.newRequest(widgetRequest, true, true)

	if err != nil {
		return datadogV1.WidgetDefinition{}, err
	}
	widgetRequest.SetResponseFormat(datadogV1.FORMULAANDFUNCTIONRESPONSEFORMAT_SCALAR)

	def := datadogV1.NewSunburstWidgetDefinition([]datadogV1.SunburstWidgetRequest{*widgetRequest}, datadogV1.SUNBURSTWIDGETDEFINITIONTYPE_SUNBURST)
	def.SetTitle(pc.panel.Title)
	def.SetTitleSize("16")

	return datadogV1.SunburstWidgetDefinitionAsWidgetDefinition(def), nil
}
