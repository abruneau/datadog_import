package browser

import (
	"dynatrace_to_datadog/datadog"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

var clickTypeMap = map[int]datadog.ClickType{
	0: datadog.CLICK_TYPE_PRIMARY,
	1: datadog.CLICK_TYPE_CONTEXTUAL,
}

func parseClickType(evt *event.Click) datadog.ClickParams {
	var params datadog.ClickParams

	if evt.Target != nil {
		params.Element.MultiLocator = parseLocators(evt.Target.Locators)
		params.Element.UserLocator = getUserLocator(evt.Target.Locators)
	}

	params.ClickType = clickTypeMap[evt.Button]

	return params
}

func ParseClickStep(evt *event.Click) (steps []datadogV1.SyntheticsStep) {
	step := datadogV1.NewSyntheticsStep()
	step.Name = &evt.Description
	step.Type = datadogV1.SYNTHETICSSTEPTYPE_CLICK.Ptr()
	falseValue := false
	trueValue := true
	step.AllowFailure = &falseValue
	step.IsCritical = &trueValue
	step.NoScreenshot = &falseValue

	if evt.Wait != nil && evt.Wait.TimeoutInMilliseconds != nil {
		timeout := int64(*evt.Wait.TimeoutInMilliseconds)
		if timeout > 5000 {
			timeout = 5000
		}
		step.Timeout = &timeout
	}

	step.Params = parseClickType(evt)
	validations := ParseValidation(evt.Validate)

	steps = append(steps, *step)
	steps = append(steps, validations...)
	return
}
