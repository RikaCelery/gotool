/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package reverse

import (
	"github.com/spf13/cobra"
	"gotool/cmd"
)

// reverseCmd represents the reverse command
var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "A tool reverse bytes blocks of files",
	Long: `A tool reverse bytes blocks of files
  input : <block1><block2><block3>...
  output: <1kcolb><2kcolb><3kcolb>...`,
	Run: func(cmd *cobra.Command, args []string) {
		main() // legacy
	},
}

func init() {
	cmd.RootCmd.AddCommand(reverseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reverseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reverseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
