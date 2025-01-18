package chat

import "strings"

const specificTags = "Specific Tags:"

func parseResults(s string) (tags []string) {
	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, specificTags) {
			tagStr := strings.TrimSpace(strings.TrimPrefix(line, specificTags))
			tagStr = strings.Trim(tagStr, "[]")
			foundTags := strings.Split(tagStr, ",")
			for _, foundTag := range foundTags {
				tags = append(tags, strings.TrimSpace(foundTag))
			}
		}
	}

	return
}
