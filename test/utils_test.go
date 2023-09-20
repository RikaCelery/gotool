package test

import (
	"fmt"
	"strings"
	"testing"

	"gotool/utils"
)

func TestLines(t *testing.T) {
	tests := []string{
		"bb",
		"dddd",
		"你好",
		"世界",
		"你好世界",
	}
	result := utils.Lines(strings.Join(tests, "\n"))
	for i, _ := range result {
		if result[i] != tests[i] {
			t.Errorf("error %s %s", result[i], tests[i])
		}
	}
}
func TestContains(t *testing.T) {
	strings := []string{
		"a",
		"bb",
		"ccc",
		"dddd",
		"你好",
		"世界",
		"你好世界",
	}
	testsNot := []string{
		"",
		"s",
		"X",
		"d",
		"cc",
		"你好世界你好世界",
		"aa",
		"b",
		"cdcc",
		"ddadd",
		"你a好",
		"世s界",
		"s你好世界",
	}
	testsTrue := []string{
		"a",
		"bb",
		"ccc",
		"dddd",
		"你好",
		"世界",
		"你好世界",
	}
	for _, test := range testsNot {
		if utils.Contains(strings, test) {
			t.Errorf("should not contains \"%s\"", test)
		}
	}
	for _, test := range testsTrue {
		if !utils.Contains(strings, test) {
			t.Errorf("should contains \"%s\"", test)
		}
	}
}
func TestSplitToUnicodeChars(t *testing.T) {
	result := utils.SplitToUnicodeChars("ɑαа⍺ａ𝐚𝑎𝒂𝒶ℂℭⅭ⊂ⲤꓚＣ𐊢𐌂𐐕𝐂𝐶𝓪𝔞𝔐𝕸ｍ𝕒𝖆𝖺𝗮𝘢𝙖ᴢꮓｚ𝐳𝑧𝒛𝓏𝔃𝕫𝗓𝚊𝛂𝛼𝜶𝝰𝞪")
	expected := []string{"ɑ",
		"α",
		"а",
		"⍺",
		"ａ",
		"𝐚",
		"𝑎",
		"𝒂",
		"𝒶",
		"ℂ",
		"ℭ",
		"Ⅽ",
		"⊂",
		"Ⲥ",
		"ꓚ",
		"Ｃ",
		"𐊢",
		"𐌂",
		"𐐕",
		"𝐂",
		"𝐶",
		"𝓪",
		"𝔞",
		"𝔐",
		"𝕸",
		"ｍ",
		"𝕒",
		"𝖆",
		"𝖺",
		"𝗮",
		"𝘢",
		"𝙖",
		"ᴢ",
		"ꮓ",
		"ｚ",
		"𝐳",
		"𝑧",
		"𝒛",
		"𝓏",
		"𝔃",
		"𝕫",
		"𝗓",
		"𝚊",
		"𝛂",
		"𝛼",
		"𝜶",
		"𝝰",
		"𝞪",
	}
	if len(result) == len(expected) {
		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("Miss match: %s %s", result[i], expected[i])
			}
		}
	} else {
		t.Errorf("Length miss match: %d %d", len(result), len(expected))
	}
}
func TestForEachUnicode(t *testing.T) {
	var result = make([]string, 0)
	utils.ForEachUnicode("ɑαа⍺ａ𝐚𝑎𝒂𝒶ℂℭⅭ⊂ⲤꓚＣ𐊢𐌂𐐕𝐂𝐶𝓪𝔞𝔐𝕸ｍ𝕒𝖆𝖺𝗮𝘢𝙖ᴢꮓｚ𝐳𝑧𝒛𝓏𝔃𝕫𝗓𝚊𝛂𝛼𝜶𝝰𝞪", func(r rune) {
		result = append(result, fmt.Sprintf("%c", r))
	})
	expected := []string{"ɑ",
		"α",
		"а",
		"⍺",
		"ａ",
		"𝐚",
		"𝑎",
		"𝒂",
		"𝒶",
		"ℂ",
		"ℭ",
		"Ⅽ",
		"⊂",
		"Ⲥ",
		"ꓚ",
		"Ｃ",
		"𐊢",
		"𐌂",
		"𐐕",
		"𝐂",
		"𝐶",
		"𝓪",
		"𝔞",
		"𝔐",
		"𝕸",
		"ｍ",
		"𝕒",
		"𝖆",
		"𝖺",
		"𝗮",
		"𝘢",
		"𝙖",
		"ᴢ",
		"ꮓ",
		"ｚ",
		"𝐳",
		"𝑧",
		"𝒛",
		"𝓏",
		"𝔃",
		"𝕫",
		"𝗓",
		"𝚊",
		"𝛂",
		"𝛼",
		"𝜶",
		"𝝰",
		"𝞪",
	}
	if len(result) == len(expected) {
		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("Miss match: %s %s", result[i], expected[i])
			}
		}
	} else {
		t.Errorf("Length miss match: %d %d", len(result), len(expected))
	}
}
