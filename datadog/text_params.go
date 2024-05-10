package datadog

type TextParams struct {
	Delay   int    `json:"delay,omitempty"`
	Value   string `json:"value,omitempty"`
	Element struct {
		URL             string       `json:"url,omitempty"`
		MultiLocator    MultiLocator `json:"multiLocator,omitempty"`
		TargetOuterHTML string       `json:"targetOuterHTML,omitempty"`
		UserLocator     UserLocator  `json:"userLocator,omitempty"`
	} `json:"element,omitempty"`
}
