package datadog

type AssertElementPresentParams struct {
	Element struct {
		URL             string      `json:"url,omitempty"`
		TargetOuterHTML string      `json:"targetOuterHTML,omitempty"`
		UserLocator     UserLocator `json:"userLocator,omitempty"`
	} `json:"element,omitempty"`
}

type AssertElementContentParams struct {
	Element struct {
		URL             string      `json:"url,omitempty"`
		TargetOuterHTML string      `json:"targetOuterHTML,omitempty"`
		UserLocator     UserLocator `json:"userLocator,omitempty"`
	} `json:"element,omitempty"`
	Value string `json:"value"`
}
