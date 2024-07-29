package browser

import (
	"dynatrace_to_datadog/dynatrace/api/rest"
	"dynatrace_to_datadog/dynatrace/api/setting"
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:synthetic:monitors:browser"

type DefaultService struct {
	schemaID string
	client   rest.Client
	options  setting.ServiceOptions
}

func Service(credentials *settings.Credentials) *DefaultService {
	return &DefaultService{
		schemaID: SchemaID,
		client:   rest.DefaultClient(credentials.URL, credentials.Token),
		options: setting.ServiceOptions{
			Get:   settings.Path("/api/v1/synthetic/monitors/%s"),
			List:  settings.Path("/api/v1/synthetic/monitors"),
			Stubs: &monitors.Monitors{},
		},
	}
}

func (me *DefaultService) Get(id string) ([]byte, error) {
	return me.client.Get(me.getURL(id), 200).Raw()
}

func (me *DefaultService) getURL(id string) string {
	if me.options.Get != nil {
		return me.options.Get(id)
	}
	panic("service options must contain a function that provides the GET URL")
}

func (me *DefaultService) listURL() string {
	if me.options.List != nil {
		return me.options.List()
	}
	panic("service options must provide an URL to list records")
}

func (me *DefaultService) List(filter string) (api.Stubs, error) {
	var err error
	var url string
	if filter != "" {
		url = me.listURL() + "?" + filter
	} else {
		url = me.listURL()
	}
	req := me.client.Get(url, 200)
	stubs := me.stubs()
	if err = req.Finish(stubs); err != nil {
		return nil, err
	}

	res := stubs.ToStubs()
	m := map[string]*api.Stub{}
	for _, stub := range res {
		m[stub.ID] = stub
	}
	res = api.Stubs{}
	for _, stub := range m {
		res = append(res, stub)
	}
	return res.ToStubs(), nil
}

func (me *DefaultService) stubs() api.RecordStubs {
	if me.options.Stubs != nil {
		stubsType := reflect.ValueOf(me.options.Stubs).Type()
		if stubsType.Kind() == reflect.Pointer {
			return reflect.New(stubsType.Elem()).Interface().(api.RecordStubs)
		}
		panic("no pointer")
	}
	return &api.StubList{}
}
