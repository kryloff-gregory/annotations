package helpers

import (
	"regexp"
)

var urlRegexp, _ = regexp.Compile("https://youtu.be/[a-zA-Z0-9]+$")

func IsValidURL(link string) bool {
	return urlRegexp.MatchString(link)
}
