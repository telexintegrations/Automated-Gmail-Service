package handlers

import (
	"regexp"
	"strings"
)

// Function to remove HTML tags from a string
func StripHTMLTags(input string) string {
	// Regex to match HTML tags
	re := regexp.MustCompile("<.*?>")
	// Replace HTML tags with an empty string
	cleaned := re.ReplaceAllString(input, "")
	// Trim spaces
	return strings.TrimSpace(cleaned)
}
