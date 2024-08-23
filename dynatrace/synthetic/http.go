package synthetic

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	dynatrace "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings/validation"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/request"
)

type httpTestConverter struct {
	ddTest      *datadogV1.SyntheticsAPITest
	dynaMonitor *dynatrace.SyntheticMonitor
	variables   []datadogV1.SyntheticsConfigVariable
	comment     []string
}

func ConvertHTTPTest(ctx context.Context, monitor *dynatrace.SyntheticMonitor, customTags []string) (*datadogV1.SyntheticsAPITest, error) {
	converter := httpTestConverter{
		dynaMonitor: monitor,
		ddTest:      datadogV1.NewSyntheticsAPITestWithDefaults(),
	}

	converter.setDefaults()
	converter.setTags(customTags)
	converter.parseHTTPRequests()
	converter.setVariables()
	converter.ddTest.Message = strings.Join(converter.comment, "\n")

	return converter.ddTest, nil
}

func (c *httpTestConverter) setDefaults() {
	var status = datadogV1.SYNTHETICSTESTPAUSESTATUS_PAUSED
	c.ddTest.Status = &status
	c.ddTest.Locations = append(c.ddTest.Locations, "aws:eu-central-1")
	var subtype datadogV1.SyntheticsTestDetailsSubType
	if len(c.dynaMonitor.Script.Requests) == 1 {
		subtype = datadogV1.SYNTHETICSTESTDETAILSSUBTYPE_HTTP
	} else {
		subtype = datadogV1.SYNTHETICSTESTDETAILSSUBTYPE_MULTI
	}
	c.ddTest.Subtype = &subtype

	c.ddTest.Name = c.dynaMonitor.Name
	frequency := int64(c.dynaMonitor.FrequencyMin * 60)
	// Value of 'tick_every' should be more than 30
	if frequency < 30 {
		frequency = 30
	}
	c.ddTest.Options.TickEvery = &frequency
}

func (c *httpTestConverter) setTags(customTags []string) {
	c.ddTest.Tags = append(getTags(c.dynaMonitor.Tags), customTags...)
}

func (c *httpTestConverter) parseHTTPRequests() {
	steps := make([]datadogV1.SyntheticsAPIStep, len(c.dynaMonitor.Script.Requests))
	for i, r := range c.dynaMonitor.Script.Requests {
		steps[i] = c.parseHTTPRequest(r)
	}
	if len(steps) == 1 {
		c.ddTest.Config.Assertions = steps[0].SyntheticsAPITestStep.Assertions
		c.ddTest.Config.Request = &steps[0].SyntheticsAPITestStep.Request
		c.ddTest.Config.Request.AllowInsecure = nil
		c.ddTest.Config.Request.FollowRedirects = nil
		c.ddTest.Options.AllowInsecure = steps[0].SyntheticsAPITestStep.Request.AllowInsecure
		c.ddTest.Options.FollowRedirects = steps[0].SyntheticsAPITestStep.Request.FollowRedirects
	} else {
		// Maximum number of elements in parameter 'steps' should be 10
		if len(steps) > 10 {
			c.comment = append(c.comment, "Number of requests exceeds 10. Only the first 10 requests will be converted.")
			c.ddTest.Config.Steps = steps[:10]
		} else {
			c.ddTest.Config.Steps = steps
		}
	}
}

func (c *httpTestConverter) parseHTTPRequest(req *dynatrace.Request) (step datadogV1.SyntheticsAPIStep) {
	subtype := datadogV1.SYNTHETICSAPITESTSTEPSUBTYPE_HTTP
	var name = ""
	if req.Description != nil {
		name = *req.Description
	}
	request := datadogV1.NewSyntheticsTestRequest()
	if req.URL != "" {
		url, vars := parseVariables(req.URL)
		request.Url = &url
		c.variables = append(c.variables, vars...)
	}
	request.Method = &req.Method
	if req.RequestBody != nil {
		b, vars := parseVariables(*req.RequestBody)
		request.Body = &b
		c.variables = append(c.variables, vars...)
	}
	if req.RequestTimeout != nil {
		timeout := float64(*req.RequestTimeout)
		const maxTimeout = 5000
		if timeout > maxTimeout {
			timeout = maxTimeout
		}
		request.Timeout = &timeout
	}
	if req.Configuration != nil {
		var vars []datadogV1.SyntheticsConfigVariable
		request.Headers, vars = parseHTTPHeaders(req.Configuration.RequestHeaders)
		request.AllowInsecure = &req.Configuration.AcceptAnyCertificate
		request.FollowRedirects = &req.Configuration.FollowRedirects
		request.NoSavingResponseBody = req.Configuration.SensitiveData
		c.variables = append(c.variables, vars...)
	}
	assertions := c.parseHTTPValidation(req.Validation)
	if req.PreProcessing != nil {
		c.comment = append(c.comment, fmt.Sprintf("Pre-processing script not supported: %s", *req.PreProcessing))
	}
	if req.PostProcessing != nil {
		c.comment = append(c.comment, fmt.Sprintf("Post-processing script not supported: %s", *req.PostProcessing))
	}
	return datadogV1.SyntheticsAPIStep{SyntheticsAPITestStep: datadogV1.NewSyntheticsAPITestStep(assertions, name, *request, subtype)}
}

var httpRuleTypeMap = map[validation.Type]func(rule validation.Rule) (datadogV1.SyntheticsAssertion, error){
	validation.Types.PatternConstraint:               parseBodyAssertion,
	validation.Types.RegexConstraint:                 parseBodyAssertion,
	validation.Types.HTTPStatusesList:                parseHTTPStatusesList,
	validation.Types.CertificateExpiryDateConstraint: parseCertificateAssertion,
}

func (c *httpTestConverter) parseHTTPValidation(val *validation.Settings) []datadogV1.SyntheticsAssertion {
	var assertions []datadogV1.SyntheticsAssertion = []datadogV1.SyntheticsAssertion{}

	if val != nil {
		for _, rule := range val.Rules {
			a, err := c.parseHTTPAssertion(*rule)
			if err != nil {
				c.comment = append(c.comment, fmt.Sprintf("Error parsing validation rule: %s", err))
			} else {
				assertions = append(assertions, a)
			}
		}
	}
	// Minimum number of elements in parameter 'assertions' should be 1
	if len(assertions) == 0 {
		assertions = append(assertions, datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(datadogV1.SYNTHETICSASSERTIONOPERATOR_IS, 200, datadogV1.SYNTHETICSASSERTIONTYPE_STATUS_CODE)})
		c.comment = append(c.comment, "No validation rules found. Defaulting to HTTP status code 200")
	}
	return assertions
}

func (c *httpTestConverter) parseHTTPAssertion(rule validation.Rule) (datadogV1.SyntheticsAssertion, error) {
	f := httpRuleTypeMap[rule.Type]
	return f(rule)
}

func parseCertificateAssertion(rule validation.Rule) (datadogV1.SyntheticsAssertion, error) {

	return datadogV1.SyntheticsAssertion{}, fmt.Errorf("parseCertificateAssertion is not supported in HTTP tests. It should be implemented in SSL tests")
}

func parseBodyAssertion(rule validation.Rule) (datadogV1.SyntheticsAssertion, error) {
	t := datadogV1.SYNTHETICSASSERTIONTYPE_BODY
	value := rule.Value
	var operator datadogV1.SyntheticsAssertionOperator
	if rule.Type == validation.Types.RegexConstraint {
		operator = datadogV1.SYNTHETICSASSERTIONOPERATOR_MATCHES
	} else {
		operator = datadogV1.SYNTHETICSASSERTIONOPERATOR_IS
	}

	operator = negateIfNeeded(operator, rule.PassIfFound)
	return datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(operator, value, t)}, nil

}

var operatorMap = map[string]datadogV1.SyntheticsAssertionOperator{
	"==": datadogV1.SYNTHETICSASSERTIONOPERATOR_IS,
	"!=": datadogV1.SYNTHETICSASSERTIONOPERATOR_IS_NOT,
	">":  datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN,
	"<":  datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN,
	">=": datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN_OR_EQUAL,
	"<=": datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN_OR_EQUAL,
}

func parseHTTPStatusesList(rule validation.Rule) (datadogV1.SyntheticsAssertion, error) {
	t := datadogV1.SYNTHETICSASSERTIONTYPE_STATUS_CODE
	var operator datadogV1.SyntheticsAssertionOperator
	statusCode, err := strconv.Atoi(rule.Value)
	if err == nil {
		operator = negateIfNeeded(datadogV1.SYNTHETICSASSERTIONOPERATOR_IS, rule.PassIfFound)
		return datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(operator, statusCode, t)}, nil
	}

	statusCode, err = strconv.Atoi(rule.Value[len(rule.Value)-3:])
	if err != nil {
		return datadogV1.SyntheticsAssertion{}, err
	}

	op := strings.TrimSpace(rule.Value[:len(rule.Value)-3])
	operator, ok := operatorMap[op]
	if !ok {
		operator = datadogV1.SYNTHETICSASSERTIONOPERATOR_IS
		return datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(negateIfNeeded(operator, rule.PassIfFound), statusCode, t)}, UnknownSyntheticsAssertionOperatorError(rule.Value, string(operator))
	}

	validOperators := []datadogV1.SyntheticsAssertionOperator{
		datadogV1.SYNTHETICSASSERTIONOPERATOR_IS,
		datadogV1.SYNTHETICSASSERTIONOPERATOR_IS_NOT,
		datadogV1.SYNTHETICSASSERTIONOPERATOR_MATCHES,
		datadogV1.SYNTHETICSASSERTIONOPERATOR_DOES_NOT_MATCH,
	}

	var isValidOperator bool
	for _, validOperator := range validOperators {
		if operator == validOperator {
			isValidOperator = true
			break
		}
	}

	if !isValidOperator {
		return datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(negateIfNeeded(datadogV1.SYNTHETICSASSERTIONOPERATOR_IS, rule.PassIfFound), statusCode, t)}, InvalidSyntheticsAssertionOperatorError(string(operator), string(datadogV1.SYNTHETICSASSERTIONOPERATOR_IS))

	}
	return datadogV1.SyntheticsAssertion{SyntheticsAssertionTarget: datadogV1.NewSyntheticsAssertionTarget(negateIfNeeded(operator, rule.PassIfFound), statusCode, t)}, nil
}

func parseHTTPHeaders(headers request.Headers) (ddHeaders map[string]string, variables []datadogV1.SyntheticsConfigVariable) {
	ddHeaders = make(map[string]string) // Initialize the map
	for _, h := range headers {
		var vars []datadogV1.SyntheticsConfigVariable
		ddHeaders[h.Name], vars = parseVariables(h.Value)
		variables = append(variables, vars...)
	}
	return
}

// parseVariables extracts text between curly brackets in a given string.
// It returns a slice of strings containing the extracted text.
func parseVariables(input string) (formattedString string, variables []datadogV1.SyntheticsConfigVariable) {
	re := regexp.MustCompile(`\{(^\\\s]*?)\}`)
	for _, v := range re.FindAllStringSubmatch(input, -1) {
		variables = append(variables, *datadogV1.NewSyntheticsConfigVariable(v[1], datadogV1.SYNTHETICSCONFIGVARIABLETYPE_TEXT))
	}
	formattedString = re.ReplaceAllString(input, "{{$1}}")
	return
}

func (c *httpTestConverter) setVariables() {
	slices.SortFunc(c.variables, func(a, b datadogV1.SyntheticsConfigVariable) int {
		return strings.Compare(a.Name, b.Name)
	})
	c.ddTest.Config.ConfigVariables = slices.CompactFunc(c.variables, func(i, j datadogV1.SyntheticsConfigVariable) bool {
		return i.Name == j.Name
	})
}

func negateIfNeeded(op datadogV1.SyntheticsAssertionOperator, passIfFound bool) datadogV1.SyntheticsAssertionOperator {
	if passIfFound {
		return op
	}
	switch op {
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_IS:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_IS_NOT
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_IS_NOT:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_IS
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN_OR_EQUAL
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN_OR_EQUAL
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN_OR_EQUAL:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_LESS_THAN_OR_EQUAL:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_MORE_THAN
	case datadogV1.SYNTHETICSASSERTIONOPERATOR_MATCHES:
		return datadogV1.SYNTHETICSASSERTIONOPERATOR_DOES_NOT_MATCH
	}

	return op
}
