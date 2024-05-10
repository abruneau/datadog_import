package datadog

type MultiLocator struct {
	AB  string `json:"ab,omitempty"`
	AT  string `json:"at,omitempty"`
	CL  string `json:"cl,omitempty"`
	CO  string `json:"co,omitempty"`
	RO  string `json:"ro,omitempty"`
	CTL string `json:"ctl,omitempty"`
}
