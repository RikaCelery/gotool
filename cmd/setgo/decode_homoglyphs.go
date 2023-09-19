package setgo

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"

	"gotool/bitlist"
	"gotool/utils"
)

func decode_advanced() {
	var log = ""
	var list = bitlist.NewBitList()

	utils.ForEachUnicode(text, func(r rune) {
		key, i, bList, err := findKeyAndIndex(string(r))
		if err == nil {
			list.Join(bList)
			log += fmt.Sprintf("key:%c,index:%d(%s)\n", key, i, bList.String())
		}
	})

	fmt.Printf("卧槽，你说的是: %s\n", string(list.ToBytes()))
	// println(list.String())
	// println(log)

}
func findKeyAndIndex(s string) (key rune, index int, bin *bitlist.BitList, err error) {
	for k, v := range mapping {
		for i, char := range v {
			if char == s {

				i2 := int(math.Log2(float64(len(v))))
				bitString := fmt.Sprintf("%b", i)
				for len(bitString) < i2 {
					bitString = "0" + bitString
				}
				// r, _ := utf8.DecodeRuneInString(s)
				// fmt.Printf("OK: %s %U\n", s, r)
				return k, i, bitlist.BitStringToBitList(bitString), nil
			}
		}
	}
	r, _ := utf8.DecodeRuneInString(s)
	return 0, 0, nil, errors.New(fmt.Sprintf("error: %s %U\n", s, r))
}
