package dynatrace

import (
	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/dynatrace/api"
	"dynatrace_to_datadog/dynatrace/api/synthetic/monitors/browser"

	dyna "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
)

// APIReader reads data from an API
type APIReader struct {
	client *browser.DefaultService
	tests  dyna.Stubs
	index  int
}

func (conf *Config) NewAPIReader() (*APIReader, error) {
	client := api.NewClient(conf.URL, conf.ApiKey)
	tests, err := client.List()
	if err != nil {
		return nil, err
	}
	return &APIReader{
		client: client,
		tests:  tests,
		index:  0,
	}, nil
}

// Read reads data from the API
func (ar *APIReader) Read() (name string, data []byte, err error) {
	if ar.index >= len(ar.tests) {
		err = common.ErrNoMoreData
		return
	}

	name = ar.tests[ar.index].Name

	data, err = ar.client.Get(ar.tests[ar.index].ID)

	ar.index++
	return
}
