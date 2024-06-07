package browser

import (
	"dynatrace_to_datadog/datadog"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func parseJavaScriptType(evt *event.Javascript) datadog.JavascriptParams {
	var params datadog.JavascriptParams

	// Ensure there is q return statement
	params.Code = fmt.Sprintf("%s\nreturn", evt.Javascript)

	return params
}

func ParseJavascriptStep(evt *event.Javascript) datadogV1.SyntheticsStep {
	step := *datadogV1.NewSyntheticsStep()
	step.Name = &evt.Description
	step.Type = datadogV1.SYNTHETICSSTEPTYPE_ASSERT_FROM_JAVASCRIPT.Ptr()
	falseValue := false
	trueValue := true
	step.AllowFailure = &falseValue
	step.IsCritical = &trueValue
	step.NoScreenshot = &falseValue

	if evt.Wait != nil && evt.Wait.TimeoutInMilliseconds != nil {
		timeout := int64(*evt.Wait.TimeoutInMilliseconds) * 1000
		step.Timeout = &timeout
	}
	step.Params = parseJavaScriptType(evt)

	return step
}
