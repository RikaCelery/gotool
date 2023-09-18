/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package setgo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gotool/utils"
)

var tagStegoCmd = &cobra.Command{
	Use:   "tags",
	Short: `ASCII map to Unicode Tags (U+E0000 to U+E007F)`,
	Long:  `Unicode Text Steganography Encoders/Decoders`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("input").Changed {
			cmd.Usage()
			return
		}
		input, err := getInput(cover)
		if err != nil {
			panic(err)
			return
		}
		if cmd.Flag("message").Changed {

			content, err := getInput(message)
			if err != nil {
				panic(err)
				return
			}
			s := encode(input, content)
			if cmd.Flag("output").Changed {
				file, err := os.OpenFile(cmd.Flag("output").Value.String(), os.O_WRONLY|os.O_CREATE, 755)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				file.WriteString(s)
			} else {
				fmt.Println(s)
			}
		} else {
			clean := decode(input)
			if print_clean {
				fmt.Printf("clean text: %s\n", clean)
			}
		}
	},
}

func encode(input string, message string) string {
	var result = ""
	var encoded = make([]string, 0)
	var appendIdx = 0
	for _, i := range message {
		encoded = append(encoded, fmt.Sprintf("%c", i+0xe0000))
		// fmt.Printf("%c U+%04x\n", i, i+0xe0000)
	}
	for _, c := range input {
		if fmt.Sprintf("%c", c) == " " {
			result += " "
			if appendIdx < len(encoded) && !tailTags {
				result += encoded[appendIdx]
				appendIdx++
			}
		} else {
			result += fmt.Sprintf("%c", c)
		}
	}
	var tails = 0
	for appendIdx < len(encoded) {
		result += encoded[appendIdx]
		appendIdx++
		tails++
	}
	// fmt.Printf("encoded success, tail tags: %d\n", tails)
	return result
}

func decode(input string) string {
	var cleanText string
	for _, c := range input {
		if !(c <= 0xe007f && c >= 0xe0000) {
			cleanText += fmt.Sprintf("%c", c)
		} else {
			fmt.Printf("%c", c-0xe0000)
			if detail {
				fmt.Printf(" <=> U+%04X\n", c)
			}
		}
	}
	println()
	return cleanText
}

var cover string
var message string
var tailTags bool
var detail bool
var print_clean bool

func getInput(pathOrRaw string) (string, error) {
	if utils.IsExist(pathOrRaw) {
		bytes, err := os.ReadFile(pathOrRaw)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		return pathOrRaw, nil
	}
}
func init() {
	stegoCmd.AddCommand(tagStegoCmd)
	tagStegoCmd.Flags().StringVar(&cover, "input", "", "input text(to decode/encode)")
	tagStegoCmd.Flags().StringVar(&message, "message", "", "message to hidden")
	tagStegoCmd.Flags().StringP("output", "o", "", "output path")
	tagStegoCmd.Flags().BoolVar(&tailTags, "tail_tags", false, "put all tags to tail of input")
	tagStegoCmd.Flags().BoolVar(&detail, "detail", false, "show mapping info")
	tagStegoCmd.Flags().BoolVar(&print_clean, "print_clean", false, "show clean text")

}
