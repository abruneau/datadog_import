package synthetic

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
)

func getTags(input monitors.TagsWithSourceInfo) []string {
	tags := []string{}
	for _, tag := range input {
		tags = append(tags, fmt.Sprintf("%s:%s", tag.Key, *tag.Value))
	}
	return tags
}
