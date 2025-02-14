package synthetic

import (
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
)

func getTags(input monitors.TagsWithSourceInfo) []string {
	tags := []string{}
	for _, tag := range input {
		var value string
		var key string = tag.Key
		if tag.Value != nil {
			value = *tag.Value
		} else if strings.Contains(tag.Key, ":") {
			key, value, _ = strings.Cut(tag.Key, ":")
		} else {
			continue
		}
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}
	return tags
}
