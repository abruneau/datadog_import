package synthetic

import (
	"dynatrace_to_datadog/synthetic/browser"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	dynatrace "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
	log "github.com/sirupsen/logrus"
)

var stepTypeMap = map[event.Type]*datadogV1.SyntheticsStepType{
	event.Types.Click:      datadogV1.SYNTHETICSSTEPTYPE_CLICK.Ptr(),
	event.Types.KeyStrokes: datadogV1.SYNTHETICSSTEPTYPE_TYPE_TEXT.Ptr(),
	event.Types.Javascript: datadogV1.SYNTHETICSSTEPTYPE_ASSERT_FROM_JAVASCRIPT.Ptr(),
	event.Types.Navigate:   datadogV1.SYNTHETICSSTEPTYPE_GO_TO_URL.Ptr(),
}

func ConvertBrowserTest(monitor *dynatrace.SyntheticMonitor, logger *log.Entry) *datadogV1.SyntheticsBrowserTest {
	var test = datadogV1.NewSyntheticsBrowserTestWithDefaults()
	var err error
	var variables []string

	// defaults
	var status = datadogV1.SYNTHETICSTESTPAUSESTATUS_PAUSED
	test.Status = &status
	test.Locations = append(test.Locations, "aws:eu-central-1")
	test.Options.DeviceIds = append(test.Options.DeviceIds, "laptop_large")

	test.Name = monitor.Name
	frequency := int64(monitor.FrequencyMin * 60)
	test.Options.TickEvery = &frequency
	test.Tags = getTags(monitor.Tags)

	test.Steps, variables, err = getSteps(monitor.Script.Events)
	if err != nil {
		logger.Error(err)
		return nil
	}

	test.Config, err = getBrowserConfig(monitor.Script.Events, variables)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return test
}

func getBrowserConfig(events event.Events, variables []string) (conf datadogV1.SyntheticsBrowserTestConfig, err error) {
	conf.Request, err = getRequest(events)
	if err != nil {
		return conf, err
	}
	conf.Variables = getBrowserVariables(variables)
	return
}

func getBrowserVariables(variables []string) (syntheticsVariables []datadogV1.SyntheticsBrowserVariable) {
	for _, v := range variables {
		if v != "" {
			syntheticsVariables = append(syntheticsVariables, *datadogV1.NewSyntheticsBrowserVariable(v, datadogV1.SYNTHETICSBROWSERVARIABLETYPE_TEXT))
		}
	}
	return
}

func getSteps(events event.Events) (steps []datadogV1.SyntheticsStep, variables []string, err error) {
	if len(events) < 2 {
		return steps, variables, fmt.Errorf("no step to parse")
	}
	for _, evt := range events[1:] {
		if evt.GetType() == event.Types.Click {
			steps = append(steps, *browser.ParseClickStep(evt.(*event.Click)))
		} else if evt.GetType() == event.Types.KeyStrokes {
			step, variable := browser.ParseKeyStrokesStep(evt.(*event.KeyStrokes))
			steps = append(steps, *step)
			variables = append(variables, variable)
		} else if evt.GetType() == event.Types.Navigate {
			steps = append(steps, *browser.ParseNavigateStep(evt.(*event.Navigate)))
		} else {
			step := datadogV1.NewSyntheticsStep()
			name := evt.GetDescription()
			step.Name = &name
			t, ok := stepTypeMap[evt.GetType()]
			if !ok {
				return steps, variables, fmt.Errorf("unknown step type %s", evt.GetType())
			}
			step.Type = t
			steps = append(steps, *step)
		}

	}

	return
}

func getRequest(events event.Events) (req datadogV1.SyntheticsTestRequest, err error) {
	if len(events) == 0 {
		return req, fmt.Errorf("no step found")
	}
	ev := events[0].(*event.Navigate)
	req.Url = &ev.URL
	method := "GET"
	req.Method = &method
	return req, nil
}
