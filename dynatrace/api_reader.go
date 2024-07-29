package dynatrace

import (
	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/dynatrace/api"
	"dynatrace_to_datadog/dynatrace/api/synthetic/monitors/browser"
	"encoding/json"

	dyna "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
)

// APIReader reads data from an API
type APIReader struct {
	client   *browser.DefaultService
	tests    dyna.Stubs
	testsIds []string
	index    int
}

func (conf *Config) NewAPIReader() (*APIReader, error) {
	client := api.NewClient(conf.URL, conf.ApiKey)
	var tests dyna.Stubs
	var err error
	if len(conf.IdList) > 0 {
		return &APIReader{
			client:   client,
			testsIds: conf.IdList,
			index:    0,
		}, nil
	}
	tests, err = client.List(conf.Filters)
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
func (ar *APIReader) Read() (id, name string, data []byte, err error) {

	// Read from the list of IDs
	if len(ar.testsIds) > 0 {
		if ar.index >= len(ar.testsIds) {
			err = common.ErrNoMoreData
			return
		}
		id = ar.testsIds[ar.index]

		data, err = ar.client.Get(id)
		test := &monitors.SyntheticMonitor{}
		if err := json.Unmarshal(data, test); err != nil {
			name = id
		} else {
			name = test.Name
		}
		ar.index++
		return
	}

	if ar.index >= len(ar.tests) {
		err = common.ErrNoMoreData
		return
	}

	name = ar.tests[ar.index].Name
	id = ar.tests[ar.index].ID

	data, err = ar.client.Get(id)

	ar.index++
	return
}
