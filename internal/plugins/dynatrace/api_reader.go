package dynatrace

import (
	"context"
	"datadog_import/internal/common"
	"datadog_import/internal/logctx"
	"datadog_import/internal/plugins/dynatrace/api"
	"datadog_import/internal/plugins/dynatrace/api/synthetic/monitors/browser"
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
	tests, err = client.List(buildFilter(conf.Filters))
	if err != nil {
		return nil, err
	}
	return &APIReader{
		client: client,
		tests:  tests,
		index:  0,
	}, nil
}

func buildFilter(filters Filters) string {
	var res string
	if filters.ManagementZone != "" {
		res += "managementZone:" + filters.ManagementZone + "&"
	}
	for _, tag := range filters.Tags {
		res += "tag=" + tag + "&"
	}
	if filters.Location != "" {
		res += "location:" + filters.Location + "&"
	}
	if filters.Type != "" && (filters.Type == "BROWSER" || filters.Type == "HTTP") {
		res += "type:" + filters.Type + "&"
	}
	if filters.Enabled != "" && (filters.Enabled == "true" || filters.Enabled == "false") {
		res += "enabled:" + filters.Enabled + "&"
	}
	if filters.CredentialId != "" {
		res += "credentialId:" + filters.CredentialId + "&"
	}
	if filters.CredentialOwner != "" {
		res += "credentialOwner:" + filters.CredentialOwner + "&"
	}
	if filters.AssignedApps != "" {
		res += "assignedApps:" + filters.AssignedApps + "&"
	}

	// Remove the last "&"
	if res != "" {
		res = res[:len(res)-1]
	}

	return res
}

// Read reads data from the API
func (ar *APIReader) Read(ctx context.Context) (id, name string, data []byte, err error) {
	// Read from the list of IDs
	if len(ar.testsIds) > 0 {
		if ar.index >= len(ar.testsIds) {
			err = common.ErrNoMoreData
			return
		}
		id = ar.testsIds[ar.index]
		logctx.From(ctx).Debugf("reading Dyntrace test: %s", id)
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
	logctx.From(ctx).Debugf("reading Dyntrace test: %s", id)
	data, err = ar.client.Get(id)

	ar.index++
	return
}
