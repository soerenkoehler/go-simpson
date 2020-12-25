package util

import "regexp"

// ReplaceVariable replaces scans TEXT for occurences of ${NAME} and replaces
// them with VALUE.
func ReplaceVariable(text, name, value string) string {
	if search, err := regexp.Compile(`\$\{` + name + `\}`); err == nil {
		return search.ReplaceAllString(text, value)
	}
	return text
}
