package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type DatadogWriter struct {
	ctx    context.Context
	client *datadogV1.SyntheticsApi
}

func (conf *Config) NewDatadogWriter() *DatadogWriter {
	ctx := datadog.NewDefaultContext(context.Background())
	ctx = context.WithValue(ctx, datadog.ContextServerVariables,
		map[string]string{"site": conf.Site})

	keys := make(map[string]datadog.APIKey)
	keys["apiKeyAuth"] = datadog.APIKey{Key: conf.ApiKey}
	keys["appKeyAuth"] = datadog.APIKey{Key: conf.AppKey}
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		keys,
	)

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewSyntheticsApi(apiClient)

	return &DatadogWriter{
		ctx:    ctx,
		client: api}

}

func (writer *DatadogWriter) Write(ctx context.Context, obj interface {
	MarshalJSON() ([]byte, error)
}, name string) error {
	var err error
	var r *http.Response

	switch v := obj.(type) {
	case *datadogV1.SyntheticsBrowserTest:
		_, r, err = writer.client.CreateSyntheticsBrowserTest(writer.ctx, *v)
	case *datadogV1.SyntheticsAPITest:
		_, r, err = writer.client.CreateSyntheticsAPITest(writer.ctx, *v)
	}

	if err != nil {
		msg := fmt.Sprintf("Error %v\n", err)
		if r != nil {
			msg += fmt.Sprintf("Full HTTP response: %v\n", r)
		}
		return fmt.Errorf(msg)
	}
	return nil
}
