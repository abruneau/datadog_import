package datadog

import (
	"context"
	"datadog_import/common"
	"datadog_import/converter"
)

type Config struct {
	Site   string `mapstructure:"site" doc:"The site of the Datadog intake to send data to."`
	ApiKey string `mapstructure:"api_key" doc:"The Datadog API key used to send data to Datadog"`
	AppKey string `mapstructure:"app_key" doc:"The application key used to access Datadog's programatic API"`
	Output string `mapstructure:"output" doc:"Output Directory"`
}

func (conf *Config) GetReader(ctx context.Context) (converter.Reader, error) {
	return nil, nil
}

func (conf *Config) GetTransformer(ctx context.Context) converter.Transformer {
	return nil
}

func (conf *Config) GetWriters(ctx context.Context) ([]converter.Writer, error) {
	var writers []converter.Writer
	if conf.Output != "" {
		writer, err := common.NewFileWriter(conf.Output)
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	if conf.Site != "" && conf.ApiKey != "" && conf.AppKey != "" {
		writer := conf.NewDatadogWriter()
		writers = append(writers, writer)
	}
	return writers, nil
}
