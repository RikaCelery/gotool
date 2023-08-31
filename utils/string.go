package utils

import "github.com/mattn/go-runewidth"

func Contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == (item) {
			return true
		}
	}
	return false
}

func Truncate(s string, size int) string {
	return runewidth.Truncate(s, size, "")
}
func TruncateLeft(s string, size int) string {
	return runewidth.TruncateLeft(s, size, "")
}
