package browser

import (
	"dynatrace_to_datadog/datadog"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func parseTapType(evt *event.Tap) datadog.ClickParams {
	var params datadog.ClickParams

	params.Element.MultiLocator = parseLocators(evt.Target.Locators)
	params.ClickType = clickTypeMap[evt.Button]
	params.Element.UserLocator = getUserLocator(evt.Target.Locators)

	return params
}

func ParseTapStep(evt *event.Tap) (steps []datadogV1.SyntheticsStep) {
	step := datadogV1.NewSyntheticsStep()
	step.Name = &evt.Description
	step.Type = datadogV1.SYNTHETICSSTEPTYPE_CLICK.Ptr()
	falseValue := false
	trueValue := true
	step.AllowFailure = &falseValue
	step.IsCritical = &trueValue
	step.NoScreenshot = &falseValue

	if evt.Wait != nil && evt.Wait.TimeoutInMilliseconds != nil {
		timeout := int64(*evt.Wait.TimeoutInMilliseconds) * 1000
		step.Timeout = &timeout
	}

	step.Params = parseTapType(evt)
	validations := ParseValidation(evt.Validate)

	steps = append(steps, *step)
	steps = append(steps, validations...)
	return
}
