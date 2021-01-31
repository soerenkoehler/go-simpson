package util

import (
	"bufio"
	"os"
	"regexp"
)

// ReplaceVariable replaces scans TEXT for occurences of ${NAME} and replaces
// them with VALUE.
func ReplaceVariable(text, name, value string) string {
	if search, err := regexp.Compile(`\$\{` + name + `\}`); err == nil {
		return search.ReplaceAllString(text, value)
	}
	return text
}

// FindInFile ... TODO
func FindInFile(path string, query string) []string {
	if search, err := regexp.Compile(query); err == nil {
		if file, err := os.Open(path); err == nil {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if match := search.FindStringSubmatch(scanner.Text()); match != nil {
					return match
				}
			}
		}
	}
	return []string{}
}
