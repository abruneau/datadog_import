package grafana

import (
	"crypto/tls"
	"net/url"
	"reflect"
	"testing"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
)

func TestBuildConfig(t *testing.T) {
	tests := []struct {
		name string
		conf APIConfig
		want *goapi.TransportConfig
	}{
		{
			name: "Default values",
			conf: APIConfig{},
			want: &goapi.TransportConfig{
				Host:             "",
				BasePath:         "/api",
				Schemes:          []string{"http"},
				BasicAuth:        url.UserPassword("admin", "admin"),
				OrgID:            1,
				TLSConfig:        &tls.Config{},
				NumRetries:       3,
				RetryTimeout:     0,
				RetryStatusCodes: []string{"420", "5xx"},
				HTTPHeaders:      map[string]string{},
			},
		},
		{
			name: "Custom values",
			conf: APIConfig{
				Host:     "example.com",
				BasePath: "/custom",
				Schemes:  []string{"https"},
				BasicAuth: struct {
					Username string `mapstructure:"username" doc:"Username for basic auth credentials."`
					Password string `mapstructure:"password" doc:"Password for basic auth credentials."`
				}{
					Username: "user",
					Password: "pass",
				},
				OrgID: 2,
			},
			want: &goapi.TransportConfig{
				Host:             "example.com",
				BasePath:         "/custom",
				Schemes:          []string{"https"},
				BasicAuth:        url.UserPassword("user", "pass"),
				OrgID:            2,
				TLSConfig:        &tls.Config{},
				NumRetries:       3,
				RetryTimeout:     0,
				RetryStatusCodes: []string{"420", "5xx"},
				HTTPHeaders:      map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildConfig(tt.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
