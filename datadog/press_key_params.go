package datadog

type PressKeyParams struct {
	Delay int    `json:"delay,omitempty"`
	Value string `json:"value"`
}
