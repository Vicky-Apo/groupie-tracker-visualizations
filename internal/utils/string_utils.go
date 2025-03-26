package utils

import "strings"

// ReplaceSpaces replaces spaces with dashes in the given string.
func ReplaceSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "-")
}

// CleanDate removes all asterisk (*) characters from a date string.
func CleanDate(date string) string {
	return strings.ReplaceAll(date, "*", "")
}
