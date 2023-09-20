/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package setgo

import (
	"fmt"

	"github.com/spf13/cobra"
	"gotool/bitlist"
)

var inputText = "asfhiabuegbghaufecilweohaxbfcuygvfuwes"
var content = ""

// advancedCmd represents the advanced command
var advancedCmd = &cobra.Command{
	Use:   "advanced",
	Short: "uses better looking Homoglyphs",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("input").Changed {
			cmd.Usage()
			return
		}
		if cmd.Flag("message").Changed {
			total := bitlist.NewBitList()
			for _, i := range content {
				total.Join(bitlist.BitStringToBitList(fmt.Sprintf("%07b", i)))
				// fmt.Printf("%c %07b %d\n", i, i, i)
			}
			var final = _encode_advanced(total)
			println(final)
		} else {
			println(_decode_advanced(inputText))
		}
	},
}

func _decode_advanced(s string) (final string) {
	total := bitlist.NewBitList()
	for _, char := range s {
		switch char {
		case ' ':
			total.Append(0, 0, 0)
		case '\u2004':
			total.Append(0, 0, 1)
		case '\u2005':
			total.Append(0, 1, 0)
		case '\u2006':
			total.Append(0, 1, 1)
		case '\u2008':
			total.Append(1, 0, 0)
		case '\u2009':
			total.Append(1, 0, 1)
		case '\u202f':
			total.Append(1, 1, 0)
		case '\u205F':
			total.Append(1, 1, 1)
		default:
			_, ok := mapper[int(char)]
			_, ok2 := mapperReverse[int(char)]
			if ok {
				total.Append(0)
			} else if ok2 {
				total.Append(1)
			}
		}
	}
	// println(total.String())
	for i := 0; i < total.Length; i += 7 {
		if i+7 < total.Length {
			final += fmt.Sprintf("%c", total.SubList(i, i+7).ToBytes()[0])
		} else {
			final += fmt.Sprintf("%c", total.SubList(i, total.Length).ToBytes()[0])
		}
	}
	return final
}
func _encode_advanced(total *bitlist.BitList) (final string) {
	index := 0
	for _, char := range inputText {
		if index >= total.Length {
			final += fmt.Sprintf("%c", char)
			continue
		}
		switch char {
		case ' ':
			list := total.SubList(index, index+3)
			// println(list.String())
			codePoint := mapperSpaces[int(list.ToBytes()[0])]
			index += 3
			final += fmt.Sprintf("%c", codePoint)

		default:
			bit := total.GetBit(index)
			codePoint, ok := mapper[int(char)]
			if ok {
				index++
				// fmt.Printf(" %c encode %d\n", codePoint, bit)
				if bit == 1 {
					final += fmt.Sprintf("%c", codePoint)
				} else {
					final += fmt.Sprintf("%c", char)
				}
			} else {
				final += fmt.Sprintf("%c", char)
			}
		}
	}
	if index != total.Length {
		fmt.Printf("there are %d bits can not be encoded, please consider longer input text\n", len(total.SubList(index, total.Length).String()))
	}
	return final
}

// 3bit
var mapperSpaces = map[int]int{
	0: ' ',
	1: '\u2004',
	2: '\u2005',
	3: '\u2006',
	4: '\u2008',
	5: '\u2009',
	6: '\u202F',
	7: '\u205F',
}
var mapperReverse = map[int]int{

	'\u039A': 'K',
	'\u039C': 'M',
	'\u039D': 'N',
	'\u03BF': 'o',
	'\u039F': 'O',
	'\u0440': 'p',
	'\u03A1': 'P',
	'\u0455': 's',
	'\u0405': 'S',
	'\u03A4': 'T',
	'\u0445': 'x',
	'\u03A7': 'X',
	'\u0443': 'y',
	'\u03A5': 'Y',
	'\u0396': 'Z',
	'\u0408': 'J',
	'\u03F3': 'j',
	'\u0456': 'I',
	'\u041D': 'H',
	'\u04BB': 'h',
	'\u050C': 'G',
	'\u0261': 'g',
	'\u0415': 'E',
	'\u0435': 'e',
	'\u0392': 'B',
	'\u0391': 'A',
	'\u0430': 'a',
}
var mapper = map[int]int{
	'a': '\u0430',
	'A': '\u0391',
	'B': '\u0392',
	'e': '\u0435',
	'E': '\u0415',
	'g': '\u0261',
	'G': '\u050C',
	'h': '\u04BB',
	'H': '\u041D',
	'i': '\u0456',
	'I': '\u0406',
	// Jj
	'j': '\u03F3',
	'J': '\u0408',
	// Kk
	'K': '\u039A',
	// Mm
	'M': '\u039C',
	// Nn
	'N': '\u039D',
	// Oo
	'o': '\u03BF',
	'O': '\u039F',
	// Pp
	'p': '\u0440',
	'P': '\u03A1',
	// Ss
	's': '\u0455',
	'S': '\u0405',
	// Tt
	'T': '\u03A4',
	// Xx
	'x': '\u0445',
	'X': '\u03A7',
	// Yy
	'y': '\u0443',
	'Y': '\u03A5',
	// Zz
	'Z': '\u0396',
}

func init() {
	stegoCmd.AddCommand(advancedCmd)
	// advancedCmd.Flags().String("mapping", "", "custom char mapping")
	advancedCmd.Flags().StringVar(&inputText, "input", "", "cover text")
	advancedCmd.Flags().StringVar(&content, "message", "", "content to hidden(only support ASCII)")
}
