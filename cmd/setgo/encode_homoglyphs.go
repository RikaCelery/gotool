package setgo

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"gotool/bitlist"
	"gotool/utils"
)

// todo_done 最后一个byte会丢
func encode_advanced() (final, tail string) {
	if hasConflict(text) {
		fmt.Println("Conflict\n")
		return "", ""
	}

	bitList := bitlist.BytesToBitList([]byte(secret))
	// println(bitList.String())
	fmt.Printf("bits to encode:%d\n", bitList.Length)
	var index = 0
	// 循环遍历字符串
	utils.ForEachUnicode(text, func(r rune) {
		// 获取当前符号和对应的替换表
		table, ok := mapping[r]
		if ok {
			// 计算最多可以编码的bit数
			bitsCount := int(math.Log2(float64(len(table))))
			if index >= bitList.Length {
				tail += fmt.Sprintf("%c", r)
				return
			}
			var replaceIndex byte = 0
			for i := 0; i < bitsCount; i++ {
				if index < bitList.Length {
					bit := bitList.GetBit(index)
					replaceIndex = replaceIndex<<1 | bit
					index++
				}
			}
			// char, _ := utf8.DecodeRuneInString(table[replaceIndex])
			// fmt.Printf(" %c:encode %d bits[%-4s](%d) > %s(%U) %d/%d\n", r, bitsCount, bitList.SubList(index-bitsCount, index).String(), replaceIndex, table[replaceIndex], char, index, bitList.Length)
			final += table[replaceIndex]
		} else {
			final += string(r)
		}
	})

	if index != bitList.Length {
		fmt.Printf("卧槽漏了%dbits信息[%d:%d](%s)\n", (bitList.Length - index), index, bitList.Length, bitList.SubList(index, bitList.Length).String())
		fmt.Printf(" - 你可能需要补充%d个字符\n", (bitList.Length-index)/4)
	}
	if len(tail) != 0 {
		fmt.Printf("没信息的东西:%s\n", tail)
	}
	fmt.Printf("弱智东西:%s\n", final)
	return final, tail
}
func hasConflict(source string) bool {
	var contain bool = false
	for _, v := range mapping {
		for _, s := range v {
			indexRune := strings.Index(source, s)
			if indexRune != -1 {
				i := indexRune - 10
				if i <= 0 {
					i = 0
				}
				j := indexRune + 10
				if j > len(source) {
					j = len(source)
				}
				r, _ := utf8.DecodeRune([]byte(s))
				fmt.Printf("contains %s(U+%04x) at %d\n\t", s, r, indexRune)
				for _, c1 := range source[i:j] {
					if i2 := strings.IndexRune(s[3:], c1); i2 != -1 {
						color.New(color.FgHiYellow).Printf("%c", c1)
					} else {
						fmt.Printf("%c", c1)
					}
				}
				fmt.Printf("\n")
				contain = true
			}
		}
	}
	return contain
}
