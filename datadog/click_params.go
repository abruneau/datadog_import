package datadog

type ClickType string

const (
	CLICK_TYPE_PRIMARY    ClickType = "primary"
	CLICK_TYPE_DOUBLE     ClickType = "double"
	CLICK_TYPE_CONTEXTUAL ClickType = "contextual"
)

type ClickParams struct {
	ClickType ClickType `json:"clickType,omitempty"`
	Element   struct {
		URL             string       `json:"url,omitempty"`
		MultiLocator    MultiLocator `json:"multiLocator,omitempty"`
		TargetOuterHTML string       `json:"targetOuterHTML,omitempty"`
		UserLocator     UserLocator  `json:"userLocator,omitempty"`
	} `json:"element,omitempty"`
}
