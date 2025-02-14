package browser

import (
	"datadog_import/internal/plugins/datadog"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func parseKeyStrokesType(evt *event.KeyStrokes) (params datadog.TextParams, variable string) {
	if evt.TextValue != nil {
		params.Value = *evt.TextValue
	} else if evt.Credential != nil {
		params.Value = fmt.Sprintf("{{ %s }}", evt.Credential.Field)
		variable = evt.Credential.Field
	}
	params.Element.MultiLocator = parseLocators(evt.Target.Locators)
	params.Element.UserLocator = getUserLocator(evt.Target.Locators)
	if evt.Wait != nil && evt.Wait.Milliseconds != nil {
		params.Delay = *evt.Wait.Milliseconds
	}

	return params, variable
}

func ParseKeyStrokesStep(evt *event.KeyStrokes) (step datadogV1.SyntheticsStep, variable string, additionalStep *datadogV1.SyntheticsStep) {
	step = *datadogV1.NewSyntheticsStep()
	step.Name = &evt.Description
	step.Type = datadogV1.SYNTHETICSSTEPTYPE_TYPE_TEXT.Ptr()

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

	step.Params, variable = parseKeyStrokesType(evt)

	return step, variable, parseSimulateReturnKey(evt.SimulateReturnKey)
}

func parseSimulateReturnKey(SimulateReturnKey bool) (step *datadogV1.SyntheticsStep) {
	if SimulateReturnKey {
		step = datadogV1.NewSyntheticsStep()
		step.Type = datadogV1.SYNTHETICSSTEPTYPE_PRESS_KEY.Ptr()
		name := "Simulate Return Key"
		step.Name = &name
		falseValue := false
		trueValue := true
		step.AllowFailure = &falseValue
		step.IsCritical = &trueValue
		step.NoScreenshot = &falseValue
		step.Params = datadog.PressKeyParams{
			Value: "Enter",
		}
		return step
	}
	return nil
}
