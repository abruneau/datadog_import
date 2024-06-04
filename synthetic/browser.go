package synthetic

import (
	"dynatrace_to_datadog/synthetic/browser"
	"fmt"
	"net/http"
	"strings"

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
	test.Options.DeviceIds = getDevice(monitor)

	test.Name = monitor.Name
	frequency := int64(monitor.FrequencyMin * 60)
	test.Options.TickEvery = &frequency
	test.Tags = getTags(monitor.Tags)

	test.Steps, variables, err = getSteps(monitor.Script.Events)
	if err != nil {
		logger.Error(err)
		return nil
	}

	test.Config, err = getBrowserConfig(monitor, variables)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return test
}

func getBrowserConfig(monitor *dynatrace.SyntheticMonitor, variables []string) (conf datadogV1.SyntheticsBrowserTestConfig, err error) {
	conf.Request, err = getRequest(monitor)
	if err != nil {
		return conf, err
	}
	conf.Variables = getBrowserVariables(variables)
	conf.SetCookie = getCookies(monitor)
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
		} else if evt.GetType() == event.Types.Javascript {
			steps = append(steps, *browser.ParseJavascriptStep(evt.(*event.Javascript)))
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

func getRequest(monitor *dynatrace.SyntheticMonitor) (req datadogV1.SyntheticsTestRequest, err error) {
	if len(monitor.Script.Events) == 0 {
		return req, fmt.Errorf("no step found")
	}
	ev := monitor.Script.Events[0].(*event.Navigate)
	req.Url = &ev.URL
	method := "GET"
	req.Method = &method
	req.Headers = getHeaders(monitor)
	return req, nil
}

func getDevice(monitor *dynatrace.SyntheticMonitor) []datadogV1.SyntheticsDeviceID {
	if *monitor.Script.Configuration.Device.Orientation == "portrait" {
		return []datadogV1.SyntheticsDeviceID{
			"mobile_small",
		}
	}
	return []datadogV1.SyntheticsDeviceID{
		"laptop_large",
	}
}

func getHeaders(monitor *dynatrace.SyntheticMonitor) map[string]string {
	var headers = map[string]string{}
	for _, h := range monitor.Script.Configuration.RequestHeaders.Headers {
		headers[h.Name] = h.Value
	}
	return headers
}

func getCookies(monitor *dynatrace.SyntheticMonitor) *string {
	var cookies []string
	for _, c := range monitor.Script.Configuration.Cookies {
		cookie := http.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Domain: c.Domain,
			Path:   *c.Path,
		}
		cookies = append(cookies, cookie.String())
	}
	res := strings.Join(cookies, "/n")
	return &res
}
