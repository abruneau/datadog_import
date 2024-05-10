package browser

import (
	"dynatrace_to_datadog/datadog"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func getUserLocator(locators event.Locators) datadog.UserLocator {
	var us datadog.UserLocator

	us.FailTestOnCannotLocate = false

	if len(locators) > 0 {
		for _, locator := range locators {
			if locator.Type == "css" {

				us.Values = append(us.Values, struct {
					Type  string "json:\"type,omitempty\""
					Value string "json:\"value,omitempty\""
				}{"css", locator.Value})
				break
			}
		}
	}

	return us
}
