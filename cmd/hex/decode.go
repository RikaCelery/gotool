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
	Short: "decode from hex",
	Long:  `decode from hex`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !cmd.Flag(inputFile).Changed {
			println("no inputs")
			os.Exit(1)
			return
		}
		var input string
		if len(args) == 0 {
			file, err := os.Open(cmd.Flag(inputFile).Value.String())
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer file.Close()
			all, err := io.ReadAll(file)
			input = string(all)
		} else {
			input = args[0]
		}

		if cmd.Flag(outputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(outputFile).Value.String(), os.O_CREATE|os.O_WRONLY, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer file.Close()
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
