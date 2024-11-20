package datadog

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	dd "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type DatadogWriter struct {
	ctx    context.Context
	client *dd.APIClient
}

func (conf *Config) NewDatadogWriter() *DatadogWriter {
	ctx := dd.NewDefaultContext(context.Background())
	ctx = context.WithValue(ctx, dd.ContextServerVariables,
		map[string]string{"site": conf.Site})

	keys := make(map[string]dd.APIKey)
	keys["apiKeyAuth"] = dd.APIKey{Key: conf.ApiKey}
	keys["appKeyAuth"] = dd.APIKey{Key: conf.AppKey}
	ctx = context.WithValue(
		ctx,
		dd.ContextAPIKeys,
		keys,
	)

	configuration := dd.NewConfiguration()
	apiClient := dd.NewAPIClient(configuration)

	return &DatadogWriter{
		ctx:    ctx,
		client: apiClient}

}

func (writer *DatadogWriter) Write(ctx context.Context, obj interface {
	MarshalJSON() ([]byte, error)
}, name string) error {
	var err error
	var r *http.Response

	switch v := obj.(type) {
	case *datadogV1.SyntheticsBrowserTest:
		api := datadogV1.NewSyntheticsApi(writer.client)
		_, r, err = api.CreateSyntheticsBrowserTest(writer.ctx, *v)
	case *datadogV1.SyntheticsAPITest:
		api := datadogV1.NewSyntheticsApi(writer.client)
		_, r, err = api.CreateSyntheticsAPITest(writer.ctx, *v)
	case *datadogV1.Dashboard:
		api := datadogV1.NewDashboardsApi(writer.client)
		_, r, err = api.CreateDashboard(writer.ctx, *v)
	}

	if err != nil {
		msg := fmt.Sprintf("Error %v\n", err)
		if r != nil {
			msg += fmt.Sprintf("Full HTTP response: %v\n", r)
		}
		return errors.New(msg)
	}
	return nil
}
