package browser

import (
	"datadog_import/internal/plugins/datadog"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func getUserLocator(locators event.Locators) datadog.UserLocator {
	var ul datadog.UserLocator

	ul.FailTestOnCannotLocate = false

	if len(locators) > 0 {
		for _, locator := range locators {
			if locator.Type == "css" {

				ul.Values = append(ul.Values, struct {
					Type  string "json:\"type,omitempty\""
					Value string "json:\"value,omitempty\""
				}{"css", locator.Value})
				break
			}
			if locator.Type == "dom" {

				ul.Values = append(ul.Values, struct {
					Type  string "json:\"type,omitempty\""
					Value string "json:\"value,omitempty\""
				}{"xpath", locator.Value})
				break
			}
		}
	}

	return ul
}
