/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package digit_converter

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gotool/cmd"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   "digit converter",
	Short: `digit converter`,
	Long: `digit converter
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, sv := range args {
			i, err := strconv.ParseInt(sv, inputDigit, 64)
			if err != nil {
				panic(err)
			}

			fmt.Printf("[%d]%s->[%d]%s\n", inputDigit, sv, outputDigit, strconv.FormatInt(i, outputDigit))
		}
	},
}
var inputDigit int
var outputDigit int

func init() {
	cmd.RootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().IntVar(&inputDigit, "from", 10, "input data")
	base64Cmd.Flags().IntVar(&outputDigit, "to", 10, "input data")

}
