package browser

import (
	"datadog_import/plugins/datadog"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
)

func ParseValidation(validations event.Validations) (steps []datadogV1.SyntheticsStep) {
	for _, validation := range validations {
		step := datadogV1.NewSyntheticsStep()
		if validation.Type == event.ValidationTypes.ContentMatch {
			step.Type = datadogV1.SYNTHETICSSTEPTYPE_ASSERT_ELEMENT_CONTENT.Ptr()
			var params datadog.AssertElementContentParams
			if validation.Target == nil {
				// TODO: handle this better.
				continue
			}
			params.Element.UserLocator = getUserLocator(validation.Target.Locators)
			params.Value = validation.Match
			name := "content_match"
			step.Name = &name
			step.Params = params
		}
		if validation.Type == event.ValidationTypes.ElementMatch {
			step.Type = datadogV1.SYNTHETICSSTEPTYPE_ASSERT_ELEMENT_PRESENT.Ptr()
			var params datadog.AssertElementPresentParams
			params.Element.UserLocator = getUserLocator(validation.Target.Locators)
			name := "element_match"
			step.Name = &name
			step.Params = params
		}
		if validation.Type == "text_match" {
			name := "Test text is present on the active page"
			step.Name = &name
			step.Type = datadogV1.SYNTHETICSSTEPTYPE_ASSERT_PAGE_CONTAINS.Ptr()
			step.Params = struct {
				Value string `json:"value"`
			}{Value: validation.Match}
		}
		steps = append(steps, *step)
	}
	return steps
}
