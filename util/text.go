package util

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func ReplaceMultiple(text string, replacements map[string]string) string {
	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}
	return text
}

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
