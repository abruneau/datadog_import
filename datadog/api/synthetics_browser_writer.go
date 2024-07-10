package api

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

type SyntheticsBrowserWriter struct {
	ctx    context.Context
	client *datadogV1.SyntheticsApi
}

func (conf *Config) NewSyntheticsBrowserWriter() *SyntheticsBrowserWriter {
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

	return &SyntheticsBrowserWriter{
		ctx:    ctx,
		client: api}

}

func (writer *SyntheticsBrowserWriter) Write(obj interface {
	MarshalJSON() ([]byte, error)
}, name string) error {
	test := obj.(*datadogV1.SyntheticsBrowserTest)
	_, _, err := writer.client.CreateSyntheticsBrowserTest(writer.ctx, *test)
	return err
}
