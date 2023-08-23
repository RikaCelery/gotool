/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package hex

import (
	hex2 "encoding/hex"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !cmd.Flag(inputFile).Changed {
			println("no inputs")
			os.Exit(1)
			return
		}
		var input []byte
		if len(args) == 0 {
			file, err := os.OpenFile(cmd.Flag(inputFile).Value.String(), os.O_CREATE, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			input, _ = io.ReadAll(file)
		} else {
			input = []byte(args[0])
		}

		if cmd.Flag(outputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(outputFile).Value.String(), os.O_CREATE, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if cmd.Flag(upperCase).Changed {
				file.WriteString(strings.ToUpper(hex2.EncodeToString(input)))
			} else {
				file.WriteString(hex2.EncodeToString(input))
			}
		} else {
			if cmd.Flag(upperCase).Changed {
				println(strings.ToUpper(hex2.EncodeToString(input)))
			} else {
				println(hex2.EncodeToString(input))
			}
		}

	},
}

func init() {
	hexCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().BoolP(upperCase, "u", false, "upperCase")
}

const upperCase = "upperCase"
