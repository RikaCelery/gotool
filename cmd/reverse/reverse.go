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
	// todo
}
