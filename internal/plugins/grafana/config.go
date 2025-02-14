package grafana

import (
	"context"
	"datadog_import/internal/common"
	"datadog_import/internal/converter"
	"datadog_import/internal/plugins/grafana/dashboard"
	"datadog_import/internal/plugins/grafana/dashboard/types"
	"encoding/json"
	"fmt"
)

type Config struct {
	Input   string    `mapstructure:"input" doc:"Input directory containing Grafana dashboards definitions"`
	API     APIConfig `mapstructure:"api" doc:"API configuration"`
	Filters Filters   `mapstructure:"filters" doc:"Filters to apply to the list of dashboards"`
}

func (conf *Config) GetReader(ctx context.Context) (converter.Reader, error) {
	if conf.Input != "" {
		return common.NewFileReader(ctx, conf.Input)
	}
	if conf.API.Host != "" {
		return conf.NewAPIReader()
	}
	if conf.Input == "" && conf.API.Host == "" {
		return nil, nil
	}
	jsonConf, _ := json.MarshalIndent(conf, "", "  ")
	return nil, fmt.Errorf("invalid Grafana configuration:\n%s", jsonConf)
}

func (conf *Config) GetTransformer(ctx context.Context) converter.Transformer {
	return func(ctx context.Context, data []byte) (interface {
		MarshalJSON() ([]byte, error)
	}, error) {
		dash := &types.Dashboard{}
		if err := json.Unmarshal(data, dash); err != nil {
			return nil, err
		}

		return dashboard.ConvertDashboard(ctx, dash)

	}
}

func (conf *Config) GetWriters(ctx context.Context) ([]converter.Writer, error) {
	return nil, nil
}
