package grafana

import (
	"context"
	"datadog_import/common"
	"datadog_import/converter"
	"datadog_import/plugins/grafana/dashboard"
	"datadog_import/plugins/grafana/dashboard/types"
	"encoding/json"
	"fmt"
)

type Config struct {
	Input string `mapstructure:"input" doc:"Input directory containing Grafana dashboards definitions"`
}

func (conf *Config) GetReader(ctx context.Context) (converter.Reader, error) {
	if conf.Input != "" {
		return common.NewFileReader(ctx, conf.Input)
	}
	// return nil, nil
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
