package grafana

import (
	"context"
	"crypto/tls"
	"datadog_import/common"
	"datadog_import/logctx"
	"encoding/json"
	"net/url"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/search"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type APIConfig struct {
	Host      string   `mapstructure:"host" doc:"Host is the domain name or IP address of the host that serves the API."`
	BasePath  string   `mapstructure:"base_path" doc:"BasePath is the URL prefix for all API paths, relative to the host root."`
	Schemes   []string `mapstructure:"schemes" doc:"Schemes are the transfer protocols used by the API (http or https)."`
	APIKey    string   `mapstructure:"api_key" doc:"APIKey is an optional API key or service account token."`
	BasicAuth struct {
		Username string `mapstructure:"username" doc:"Username for basic auth credentials."`
		Password string `mapstructure:"password" doc:"Password for basic auth credentials."`
	} `mapstructure:"basic_auth"`
	OrgID int64 `mapstructure:"org_id" doc:"OrgID provides an optional organization ID."`
}

type Filters struct {
	DashboardUIDs []string `mapstructure:"dashboard_uids" doc:"List of dashboard uid’s to search for"`
	Deleted       bool     `mapstructure:"deleted" doc:"Flag indicating if only soft deleted Dashboards should be returned"`
	FolderUIDs    []string `mapstructure:"folder_uids" doc:"List of folder UID’s to search in for dashboards"`
	Limit         int64    `mapstructure:"limit" doc:"Limit the number of returned dashboards (max 5000)"`
	Query         string   `mapstructure:"query" doc:"Search Query"`
	Starred       bool     `mapstructure:"starred" doc:"Flag indicating if only starred Dashboards should be returned"`
	Tag           []string `mapstructure:"tag" doc:"List of tags to search for"`
	Type          string   `mapstructure:"type" doc:"Type of dashboards to search for, dash-folder or dash-db"`
}

type APIReader struct {
	client *goapi.GrafanaHTTPAPI
	dashs  models.HitList
	index  int
}

func (conf *Config) NewAPIReader() (*APIReader, error) {
	client := goapi.NewHTTPClientWithConfig(strfmt.Default, buildConfig(conf.API))
	dashs, err := client.Search.Search(buildFilter(conf.Filters))
	if err != nil {
		return nil, err
	}

	return &APIReader{
		client: client,
		dashs:  dashs.Payload,
	}, nil
}

func buildConfig(conf APIConfig) *goapi.TransportConfig {
	if conf.BasePath == "" {
		conf.BasePath = "/api"
	}
	if len(conf.Schemes) == 0 {
		conf.Schemes = []string{"http"}
	}
	if conf.BasicAuth.Username == "" {
		conf.BasicAuth.Username = "admin"
	}
	if conf.BasicAuth.Password == "" {
		conf.BasicAuth.Password = "admin"
	}
	if conf.OrgID == 0 {
		conf.OrgID = 1
	}

	return &goapi.TransportConfig{
		Host:             conf.Host,
		BasePath:         conf.BasePath,
		Schemes:          conf.Schemes,
		APIKey:           conf.APIKey,
		BasicAuth:        url.UserPassword(conf.BasicAuth.Username, conf.BasicAuth.Password),
		OrgID:            conf.OrgID,
		TLSConfig:        &tls.Config{},
		NumRetries:       3,
		RetryTimeout:     0,
		RetryStatusCodes: []string{"420", "5xx"},
		HTTPHeaders:      map[string]string{},
	}

}

func buildFilter(filters Filters) *search.SearchParams {
	params := search.NewSearchParams()
	params.DashboardUIDs = filters.DashboardUIDs
	params.Deleted = &filters.Deleted
	params.FolderUIDs = filters.FolderUIDs
	params.Limit = &filters.Limit
	params.Query = &filters.Query
	params.Starred = &filters.Starred
	params.Tag = filters.Tag
	params.Type = &filters.Type
	return params
}

// Read reads data from the API
func (ar *APIReader) Read(ctx context.Context) (id, name string, data []byte, err error) {
	if ar.index >= len(ar.dashs) {
		err = common.ErrNoMoreData
		return
	}

	name = ar.dashs[ar.index].Title
	id = ar.dashs[ar.index].UID

	logctx.From(ctx).Debugf("reading Grafana dashboard: %s", id)
	dash, err := ar.client.Dashboards.GetDashboardByUID(id)
	if err != nil {
		return
	}
	data, err = json.Marshal(dash.Payload.Dashboard)
	ar.index++
	return
}
