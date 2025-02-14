package browser

import (
	"datadog_import/internal/plugins/datadog"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func parseNavigateType(evt *event.Navigate) datadog.GoToParams {
	var params datadog.GoToParams

	params.Value = evt.URL
	return params
}

func ParseNavigateStep(evt *event.Navigate) datadogV1.SyntheticsStep {
	step := *datadogV1.NewSyntheticsStep()
	step.Name = &evt.Description
	step.Type = datadogV1.SYNTHETICSSTEPTYPE_GO_TO_URL.Ptr()
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

	step.Params = parseNavigateType(evt)

	return step
}
