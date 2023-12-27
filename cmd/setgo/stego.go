/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package setgo

import (
	"fmt"

	"github.com/spf13/cobra"
	"gotool/cmd"
)

var stegoCmd = &cobra.Command{
	Use:   "setgo",
	Short: `Unicode Text Steganography Encoders/Decoders`,
	Long:  `Unicode Text Steganography Encoders/Decoders`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("[child commands]\n")
		for _, command := range cmd.Commands() {
			fmt.Printf("  %s: %s\n", command.Use, command.Short)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(stegoCmd)
	stegoCmd.AddCommand(moovStegoCmd)
}
