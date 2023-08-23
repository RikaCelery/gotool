/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package base64

import (
	"github.com/spf13/cobra"
	"gotool/cmd"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: `base64 encode/decode tool`,
	Long: `base64 encode/decode tool
support standard/custom codecs
support file input/output
`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
		println("** [V2 Command] **")
		cmd.Usage()
	},
}

func init() {
	cmd.RootCmd.AddCommand(base64Cmd)
	base64Cmd.PersistentFlags().StringP(flagData, "d", "", "input data")
	base64Cmd.PersistentFlags().StringP(flagInputFile, "f", "", "output file (leave empty using --data)")
	base64Cmd.PersistentFlags().StringP(flagOutputFile, "o", "", "output file (leave empty using stdout)")
	base64Cmd.PersistentFlags().IntP(flagBase, "b", 64, "base(32/64)")
	base64Cmd.PersistentFlags().Bool(flagHexOutput, false, "output hex string (upper case)")
	base64Cmd.PersistentFlags().String(flagEncoder, "", "custom encode/decode char sequence (default: STDCodec)")

}
