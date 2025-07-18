package main

import (
	"strings"
)

// cleanInput splits the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	textTrimmedToLower := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(textTrimmedToLower)
}
