package util

import (
	"bufio"
	"os"
	"regexp"
)

// ReplaceVariables replaces scans TEXT for occurences of ${NAME} and replaces
// them with VALUE.
func ReplaceVariables(text string, replacements map[string]string) string {
	for k, v := range replacements {
		if search, err := regexp.Compile(`\$\{` + k + `\}`); err == nil {
			text = search.ReplaceAllString(text, v)
		}
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
