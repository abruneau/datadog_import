package dashboard

import (
	"context"
	"datadog_import/logctx"
	templatevariable "datadog_import/plugins/grafana/dashboard/template_variable"
	"datadog_import/plugins/grafana/dashboard/types"
	"datadog_import/plugins/grafana/dashboard/widgets"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type dashboardConvertor struct {
	ctx               context.Context
	graf              *types.Dashboard
	templateVariables []datadogV1.DashboardTemplateVariable
	datasource        string
	widgets           []datadogV1.Widget
}

func ConvertDashboard(ctx context.Context, graf *types.Dashboard) (*datadogV1.Dashboard, error) {
	convertor := &dashboardConvertor{
		ctx:  ctx,
		graf: graf,
	}

	convertor.init()
	return convertor.build(), nil
}

func (c *dashboardConvertor) init() {
	c.getDefaultDatasource()
	c.extractTemplateVariables()
	c.extractWidgets()
}

func (c *dashboardConvertor) build() *datadogV1.Dashboard {
	dash := datadogV1.NewDashboard(datadogV1.DASHBOARDLAYOUTTYPE_ORDERED, c.graf.Title, c.widgets)
	description := fmt.Sprintf("%s\n\ngenerated with https://github.com/abruneau/datadog_import", c.graf.Description)
	dash.Description.Set(&description)
	dash.TemplateVariables = c.templateVariables
	return dash
}

func (c *dashboardConvertor) extractTemplateVariables() {
	for _, v := range c.graf.Templating.List {
		if _, ok := v.Query.(string); !ok {
			continue
		}
		if v.Type == "datasource" {
			c.datasource = v.Query.(string)
		} else if v.Type == "query" {
			tv := templatevariable.GetTemplateVariable(c.datasource, v)
			if tv != nil {
				c.templateVariables = append(c.templateVariables, *tv)
			}
		}
	}
}

func (c *dashboardConvertor) processPanel(panel types.Panel) *datadogV1.Widget {
	widget, err := widgets.ConvertWidget(c.datasource, panel)
	if err != nil {
		logctx.From(c.ctx).WithField("widget", panel.Title).Warn(err)
		return nil
	}

	for _, p := range panel.Panels {
		childWidget := c.processPanel(p)
		if childWidget == nil {
			continue
		}
		widget.Definition.GroupWidgetDefinition.Widgets = append(widget.Definition.GroupWidgetDefinition.Widgets, *childWidget)
	}

	return &widget
}

func (c *dashboardConvertor) extractWidgets() {
	var lastGroupWidget *datadogV1.Widget
	for _, panel := range c.graf.Panels {
		widget := c.processPanel(panel)

		if widget == nil {
			continue
		}

		if lastGroupWidget != nil && widget.Definition.GroupWidgetDefinition == nil {
			lastGroupWidget.Definition.GroupWidgetDefinition.Widgets = append(lastGroupWidget.Definition.GroupWidgetDefinition.Widgets, *widget)
		} else {
			c.widgets = append(c.widgets, *widget)
		}

		if widget.Definition.GroupWidgetDefinition != nil {
			lastGroupWidget = widget
		}
	}
}

func (c *dashboardConvertor) getDefaultDatasource() {
	for _, r := range c.graf.Requires {
		if r.Type == "datasource" {
			c.datasource = r.ID
			logctx.From(c.ctx).WithField("datasource", c.datasource).Debug("default datasource")
			return
		}

	}
}

// func (c *dashboardConvertor) getDefaultDatasource() {
// 	if len(c.graf.Annotations.List) > 0 {
// 		switch ds := c.graf.Annotations.List[0].Datasource.(type) {
// 		case string:
// 			c.datasource = ds
// 		case map[string]interface{}:
// 			c.datasource = ds["uid"].(string)
// 		}
// 	}
// 	logctx.From(c.ctx).WithField("datasource", c.datasource).Debug("default datasource")
// }
