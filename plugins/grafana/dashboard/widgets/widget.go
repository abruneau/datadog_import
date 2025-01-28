package widgets

import (
	"context"
	"datadog_import/logctx"
	"datadog_import/plugins/grafana/dashboard/types"
	"fmt"
	"math"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type PanelConvertor struct {
	ctx    context.Context
	source string
	panel  types.Panel
}

func NewPanelConvertor(ctx context.Context, source string, panel types.Panel) *PanelConvertor {
	return &PanelConvertor{
		ctx:    ctx,
		source: source,
		panel:  panel,
	}
}

func (pc *PanelConvertor) Convert() (datadogV1.Widget, error) {
	widget := datadogV1.NewWidgetWithDefaults()
	var definition datadogV1.WidgetDefinition
	var err error = nil

	switch pc.panel.Type {
	case "piechart":
		definition, err = pc.newPiechartDefinition()
	case "row":
		definition, err = pc.newGroupDefinition()
	case "stat", "singlestat":
		definition, err = pc.newQueryValueDefinition()
	case "text":
		definition, err = pc.newTextDefinition()
	case "timeseries", "graph", "barchart":
		definition, err = pc.newTimeseriesDefinition()
	default:
		definition, err = pc.newTimeseriesDefinition()
		if err == nil {
			logctx.From(pc.ctx).WithField("widget", pc.panel.Title).Info(fmt.Sprintf("unkown type: %s defaulting to time series", pc.panel.Type))
		}
	}

	if err != nil {
		return *widget, err
	}

	widget.SetDefinition(definition)
	widget.SetLayout(pc.convertGrid())

	return *widget, err
}

// convertGrid converts the Grafana panel grid position to a Datadog widget layout.
// It divides the height, width, x, and y coordinates by 2 and ensures the height is at least 1.
// Returns a new instance of datadogV1.WidgetLayout with the converted values.
func (pc *PanelConvertor) convertGrid() datadogV1.WidgetLayout {
	heigh := int64(math.Floor(float64(pc.panel.GridPos.H) / 2))
	if heigh == 0 {
		heigh = 1
	}
	width := int64(math.Floor(float64(pc.panel.GridPos.W) / 2))
	x := int64(math.Floor(float64(pc.panel.GridPos.X) / 2))
	y := int64(math.Floor(float64(pc.panel.GridPos.Y) / 2))
	return *datadogV1.NewWidgetLayout(heigh, width, x, y)
}
