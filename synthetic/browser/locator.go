package browser

import (
	"dynatrace_to_datadog/datadog"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func parseLocators(locators event.Locators) datadog.MultiLocator {
	var ml datadog.MultiLocator

	for i, locator := range locators {
		if i == 0 {
			ml.AB = locator.Value
		} else if i == 1 {
			ml.AT = locator.Value
		} else if i == 2 {
			ml.CL = locator.Value
		} else if i == 3 {
			ml.RO = locator.Value
		} else if i == 4 {
			ml.CTL = locator.Value
		}
	}
	return ml
}
