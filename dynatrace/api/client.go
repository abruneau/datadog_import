package api

import (
	"dynatrace_to_datadog/dynatrace/api/synthetic/monitors/browser"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func NewClient(url, apikey string) *browser.DefaultService {
	credentials := settings.Credentials{
		URL:   url,
		Token: apikey,
	}
	return browser.Service(&credentials)
}
