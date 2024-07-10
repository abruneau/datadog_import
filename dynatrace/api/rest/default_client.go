package rest

type ClientFactory func(envURL, apiToken, schemaID string) Client

func DefaultClient(envURL string, apiToken string) Client {
	return &defaultClient{envURL: envURL, apiToken: apiToken}
}

type defaultClient struct {
	envURL   string
	apiToken string
}

func (me *defaultClient) Get(url string, expectedStatusCodes ...int) Request {
	req := &request{client: me, url: url, method: "GET"}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}
