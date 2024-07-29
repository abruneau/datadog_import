package dynatrace

import (
	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/converter"
	"dynatrace_to_datadog/dynatrace/synthetic"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

type Config struct {
	URL     string   `mapstructure:"url" doc:"Dynatrace URL"`
	ApiKey  string   `mapstructure:"api_key" doc:"Dynatrace API Key"`
	Input   string   `mapstructure:"input" doc:"Input directory containing Dynatrace synthetics tests definitions"`
	Filters string   `mapstructure:"filters" doc:"Filters to apply to the list of tests separated by & symbol"`
	IdList  []string `mapstructure:"id_list" doc:"List of test IDs to fetch"`
}

func (conf *Config) GetReader() (converter.Reader, error) {
	if conf.Input != "" {
		return common.NewFileReader(conf.Input)
	}
	if conf.ApiKey != "" && conf.URL != "" {
		return conf.NewAPIReader()
	}
	return nil, fmt.Errorf("no reader found")
}

func (conf *Config) GetTransformer() converter.Transformer {
	return func(data []byte) (interface {
		MarshalJSON() ([]byte, error)
	}, error) {
		test := &monitors.SyntheticMonitor{}
		if err := json.Unmarshal(data, test); err != nil {
			return nil, err
		}
		if test.Type == monitors.Types.Browser {
			browserTest := &browser.SyntheticMonitor{}
			if err := json.Unmarshal(data, browserTest); err != nil {
				return nil, err
			}

			return synthetic.ConvertBrowserTest(browserTest)
		} else if test.Type == monitors.Types.HTTP {
			httpTest := &http.SyntheticMonitor{}
			if err := json.Unmarshal(data, httpTest); err != nil {
				return nil, err
			}

			return synthetic.ConvertHTTPTest(httpTest)
		} else {
			return nil, fmt.Errorf("SYnthetic type not supported: %s", test.Type)
		}
	}
}
