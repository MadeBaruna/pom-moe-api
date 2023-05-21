package utils

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.Replace(text, " ", "-", -1)
	text = regexp.MustCompile(`[^\w-]+`).ReplaceAllString(text, "")
	text = strings.Replace(text, "--", "-", -1)
	text = strings.Trim(text, "-")
	return text
}
