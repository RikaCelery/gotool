/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package hex

import (
	"github.com/spf13/cobra"
	"gotool/cmd"
)

// hexCmd represents the hex command
var hexCmd = &cobra.Command{
	Use:   "hex",
	Short: "hex encode/decode tool",
	Long:  `hex encode/decode tool`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("hex called")
	// },
}

func init() {
	cmd.RootCmd.AddCommand(hexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hexCmd.PersistentFlags().String("foo", "", "A help for foo")

	hexCmd.PersistentFlags().String(inputFile, "", "inputFile")
	hexCmd.PersistentFlags().String(outputFile, "", "outputFile, empty using stdout")
}

const (
	inputFile  = "input"
	outputFile = "output"
)
