/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package hex

import (
	hex2 "encoding/hex"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
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
		var input string
		if len(args) == 0 {
			file, err := os.OpenFile(cmd.Flag(inputFile).Value.String(), os.O_CREATE, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			all, err := io.ReadAll(file)
			input = string(all)
		} else {
			input = args[0]
		}

		if cmd.Flag(outputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(outputFile).Value.String(), os.O_CREATE, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			decodeString, err := hex2.DecodeString(input)
			if err != nil {
				log.Fatalln(err)
				return
			}
			file.Write(decodeString)
		} else {
			decodeString, err := hex2.DecodeString(input)
			if err != nil {
				log.Fatalln(err)
				return
			}
			println(decodeString)
		}

	},
}

func init() {
	hexCmd.AddCommand(decodeCmd)
}
