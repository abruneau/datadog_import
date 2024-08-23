package dynatrace

import (
	"context"
	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/converter"
	"dynatrace_to_datadog/dynatrace/synthetic"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

type Filters struct {
	ManagementZone  string   `mapstructure:"management_zone" doc:"Filters the resulting set of monitors to those which are part of the specified management zone ID."`
	Tags            []string `mapstructure:"tags" doc:"Filters the resulting set of monitors by specified tags."`
	Location        string   `mapstructure:"location" doc:"Filters the resulting set of monitors by specified location ID."`
	Type            string   `mapstructure:"type" doc:"Filters the resulting set of monitors to those of the specified type: BROWSER or HTTP."`
	Enabled         string   `mapstructure:"enabled" doc:"Filters the resulting set of monitors to those which are enabled (true) or disabled (false)"`
	CredentialId    string   `mapstructure:"credential_id" doc:"Filters the resulting set of monitors to those using the specified credential ID."`
	CredentialOwner string   `mapstructure:"credential_owner" doc:"Filters the resulting set of monitors to those using the specified credential owner."`
	AssignedApps    string   `mapstructure:"assigned_apps" doc:"Filters the resulting set of monitors to those assigned to the specified application IDs."`
}

type Config struct {
	URL        string   `mapstructure:"url" doc:"Dynatrace URL"`
	ApiKey     string   `mapstructure:"api_key" doc:"Dynatrace API Key"`
	Input      string   `mapstructure:"input" doc:"Input directory containing Dynatrace synthetics tests definitions"`
	Filters    Filters  `mapstructure:"filters" doc:"Filters to apply to the list of tests"`
	IdList     []string `mapstructure:"id_list" doc:"List of test IDs to fetch"`
	CustomTags []string `mapstructure:"custom_tags" doc:"List of custom tags to add to the tests"`
}

func (conf *Config) GetReader() (converter.Reader, error) {
	if conf.Input != "" {
		return common.NewFileReader(conf.Input)
	}
	if conf.ApiKey != "" && conf.URL != "" {
		return conf.NewAPIReader()
	}
	jsonConf, _ := json.MarshalIndent(conf, "", "  ")
	return nil, fmt.Errorf("invalid Dynatrace configuration:\n%s", jsonConf)
}

func (conf *Config) GetTransformer() converter.Transformer {
	return func(ctx context.Context, data []byte) (interface {
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

			return synthetic.ConvertBrowserTest(ctx, browserTest, conf.CustomTags)
		} else if test.Type == monitors.Types.HTTP {
			httpTest := &http.SyntheticMonitor{}
			if err := json.Unmarshal(data, httpTest); err != nil {
				return nil, err
			}

			return synthetic.ConvertHTTPTest(ctx, httpTest, conf.CustomTags)
		} else {
			return nil, fmt.Errorf("SYnthetic type not supported: %s", test.Type)
		}
	}
}
