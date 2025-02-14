package cmd

import (
	"context"
	"datadog_import/internal/common"
	"datadog_import/internal/converter"
	"datadog_import/internal/logctx"
	"datadog_import/internal/plugins"
	"datadog_import/internal/plugins/datadog"
	"datadog_import/internal/plugins/dynatrace"
	"datadog_import/internal/plugins/grafana"

	"github.com/sirupsen/logrus"
)

type ViperConfig struct {
	Log       string           `mapstructure:"log" doc:"Log level"`
	Dynatrace dynatrace.Config `mapstructure:"dynatrace" doc:"Dynatrace Configuration"`
	Grafana   grafana.Config   `mapstructure:"grafana" doc:"Grafana Configuration"`
	Datadog   datadog.Config   `mapstructure:"datadog" doc:"Datadog Configuration"`
}

func (c *ViperConfig) pluginList() []plugins.Plugin {
	return []plugins.Plugin{
		&c.Dynatrace,
		&c.Grafana,
		&c.Datadog,
	}
}

func (c *ViperConfig) SetLogLevel() {
	level, err := logrus.ParseLevel(c.Log)
	if err != nil {
		logctx.Default.Panic(err)
	}
	logrus.SetLevel(level)
}

func (c *ViperConfig) BuildConverter(ctx context.Context) converter.Converter {
	var (
		r  converter.Reader
		t  converter.Transformer
		ws []converter.Writer
	)
	for _, plugin := range c.pluginList() {
		p_reader, err := plugin.GetReader(ctx)
		common.Check(ctx, err)
		if p_reader != nil {
			if r != nil {
				logctx.Default.Panic("multiple readers found. Only one reader is allowed")
			}
			r = p_reader
		}

		if p_transformer := plugin.GetTransformer(ctx); p_reader != nil && p_transformer != nil {
			if t != nil {
				logctx.Default.Panic("multiple transformers found. Only one transformer is allowed")
			}
			t = p_transformer
		}

		p_ws, err := plugin.GetWriters(ctx)
		common.Check(ctx, err)
		ws = append(ws, p_ws...)
	}

	if r == nil {
		logctx.Default.Panic("no reader found")
	}

	if len(ws) == 0 {
		logctx.Default.Panic("no writers found")
	}

	return converter.NewConverter(ctx, r, t, ws)
}
