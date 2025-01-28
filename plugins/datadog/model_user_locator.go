package datadog

type UserLocator struct {
	FailTestOnCannotLocate bool `json:"failTestOnCannotLocate"`
	Values                 []struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"values,omitempty"`
}
