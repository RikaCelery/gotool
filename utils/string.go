package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

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
func SplitToUnicodeChars(s string) []string {
	var result []string

	// 循环遍历字符串
	for len(s) > 0 {
		// 从字符串中解码一个 Unicode 字符
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			// 处理解码错误
			fmt.Println("Error decoding rune")
			break
		}

		// 将解码的 Unicode 字符添加到结果切片中
		result = append(result, string(r))

		// 更新字符串，去掉已处理的字符
		s = s[size:]
	}

	return result
}
func Lines(s string) []string {
	return strings.Split(s, "\n")
}
func ForEachUnicode(text string, callback func(r rune)) {
	for len(text) > 0 {
		// 从字符串中解码一个 Unicode 字符
		r, size := utf8.DecodeRuneInString(text)
		if r == utf8.RuneError {
			// 处理解码错误
			fmt.Printf("Error decoding rune %s\n", text[size:])
			break
		}
		callback(r)
		// 更新字符串，去掉已处理的字符
		text = text[size:]
	}
}
