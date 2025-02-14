package synthetic

import (
	"context"
	"datadog_import/internal/common"
	"datadog_import/internal/logctx"
	"datadog_import/internal/plugins/dynatrace/synthetic/browser"
	"fmt"
	"net/http"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	dynatrace "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/request"
)

func ConvertBrowserTest(ctx context.Context, monitor *dynatrace.SyntheticMonitor, customTags []string) (*datadogV1.SyntheticsBrowserTest, error) {
	var test = datadogV1.NewSyntheticsBrowserTestWithDefaults()
	var err error
	var variables, additional_cookies []string

	// defaults
	var status = datadogV1.SYNTHETICSTESTPAUSESTATUS_PAUSED
	test.Status = &status
	test.Locations = append(test.Locations, "aws:eu-central-1")
	test.Options.DeviceIds = getDevice(monitor)

	test.Name = monitor.Name
	frequency := int64(monitor.FrequencyMin * 60)
	if frequency < 300 {
		frequency = 300
	}
	test.Options.TickEvery = &frequency
	test.Tags = append(getTags(monitor.Tags), customTags...)

	test.Steps, variables, additional_cookies, err = getSteps(ctx, monitor.Script.Events)
	if err != nil {
		return test, err
	}

	test.Config, err = getBrowserConfig(ctx, monitor, variables, additional_cookies)
	if err != nil {
		return test, err
	}

	return test, nil
}

func getBrowserConfig(ctx context.Context, monitor *dynatrace.SyntheticMonitor, variables, additional_cookies []string) (conf datadogV1.SyntheticsBrowserTestConfig, err error) {
	conf.Request, err = getRequest(monitor)
	if err != nil {
		return conf, err
	}
	conf.Variables = getBrowserVariables(variables)
	cookies := strings.Join(append(additional_cookies, getCookies(ctx, monitor)...), "\n")
	conf.SetCookie = &cookies
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

func getSteps(ctx context.Context, events event.Events) (steps []datadogV1.SyntheticsStep, variables []string, cookies []string, err error) {
	foundFirstNavigation := false
	for _, evt := range events {
		if evt.GetType() == event.Types.Click {
			steps = append(steps, browser.ParseClickStep(evt.(*event.Click))...)
		} else if evt.GetType() == event.Types.Tap {
			steps = append(steps, browser.ParseTapStep(evt.(*event.Tap))...)
		} else if evt.GetType() == event.Types.KeyStrokes {
			step, variable, additionalStep := browser.ParseKeyStrokesStep(evt.(*event.KeyStrokes))
			steps = append(steps, step)
			variables = append(variables, variable)
			if additionalStep != nil {
				steps = append(steps, *additionalStep)
			}
		} else if evt.GetType() == event.Types.Navigate {
			if foundFirstNavigation {
				steps = append(steps, browser.ParseNavigateStep(evt.(*event.Navigate)))
			} else {
				foundFirstNavigation = true
			}
		} else if evt.GetType() == event.Types.Javascript {
			steps = append(steps, browser.ParseJavascriptStep(evt.(*event.Javascript)))
		} else if evt.GetType() == event.Types.Cookie {
			for _, c := range evt.(*event.Cookie).Cookies {
				cookies = append(cookies, cookieToString(ctx, c))
			}
		} else {
			err = common.UnknownStepTypeError(string(evt.GetType()))
			return
		}
	}
	return
}

func getRequest(monitor *dynatrace.SyntheticMonitor) (req datadogV1.SyntheticsTestRequest, err error) {
	// Look for the first event of type Navigate
	for _, evt := range monitor.Script.Events {
		if evt.GetType() == event.Types.Navigate {
			req.Url = &evt.(*event.Navigate).URL
			method := "GET"
			req.Method = &method
			req.Headers = getHeaders(monitor)
			return req, nil
		}
	}
	return req, fmt.Errorf("no step found")
}

func getDevice(monitor *dynatrace.SyntheticMonitor) []datadogV1.SyntheticsDeviceID {
	if monitor.Script.Configuration.Device.Mobile != nil && *monitor.Script.Configuration.Device.Mobile {
		return []datadogV1.SyntheticsDeviceID{
			"mobile_small",
		}
	}
	if monitor.Script.Configuration.Device != nil && monitor.Script.Configuration.Device.Orientation != nil {
		// Desktop and laptop devices are not allowed to use the `portrait` orientation
		if *monitor.Script.Configuration.Device.Orientation == "portrait" {
			return []datadogV1.SyntheticsDeviceID{
				"mobile_small",
			}
		}
	}
	return []datadogV1.SyntheticsDeviceID{
		"laptop_large",
	}
}

func getHeaders(monitor *dynatrace.SyntheticMonitor) map[string]string {
	var headers = map[string]string{}
	if monitor.Script.Configuration.RequestHeaders != nil {
		for _, h := range monitor.Script.Configuration.RequestHeaders.Headers {
			headers[h.Name] = h.Value
		}
	}

	return headers
}

func getCookies(ctx context.Context, monitor *dynatrace.SyntheticMonitor) (cookies []string) {
	for _, c := range monitor.Script.Configuration.Cookies {
		cookies = append(cookies, cookieToString(ctx, c))
	}
	return
}

func cookieToString(ctx context.Context, c *request.Cookie) string {
	var path string = ""
	if c.Path != nil {
		path = *c.Path
	}
	cookie := http.Cookie{
		Name:   c.Name,
		Value:  c.Value,
		Domain: c.Domain,
		Path:   path,
	}

	if err := cookie.Valid(); err != nil {
		logctx.From(ctx).Warn(err)
	}

	return cookie.String()
}
