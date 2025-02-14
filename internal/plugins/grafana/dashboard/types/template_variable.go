package types

type TemplateVariable struct {
	Current struct {
	} `json:"current"`
	Hide        int           `json:"hide"`
	IncludeAll  bool          `json:"includeAll"`
	Label       string        `json:"label"`
	Multi       bool          `json:"multi"`
	Name        string        `json:"name"`
	Options     []interface{} `json:"options"`
	Query       interface{}   `json:"query"`
	QueryValue  string        `json:"queryValue,omitempty"`
	Refresh     int           `json:"refresh"`
	Regex       string        `json:"regex"`
	SkipURLSync bool          `json:"skipUrlSync"`
	Type        string        `json:"type"`
	Datasource  struct {
		UID string `json:"uid"`
	} `json:"datasource,omitempty"`
	Definition     string `json:"definition,omitempty"`
	Sort           int    `json:"sort,omitempty"`
	TagValuesQuery string `json:"tagValuesQuery,omitempty"`
	TagsQuery      string `json:"tagsQuery,omitempty"`
	UseTags        bool   `json:"useTags,omitempty"`
}
