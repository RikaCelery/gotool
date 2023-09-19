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

var text = `HashcatHashcatHashcatHashcatHa`
var secret = "Fuck The World!"
var mapping = make(map[int32][]string)

var homoCmd = &cobra.Command{
	Use:     "mhomoglyph",
	Aliases: []string{"mhomo"},
	Short:   "uses more Homoglyphs encode bits",
	Long: `mapping:
a:𝑎𝒂𝒶𝓪
B:𝗕𝘉𝘽𝙱𝚩
...`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("input").Changed {
			cmd.Usage()
			return
		}
		initMap(cmd)
		if cmd.Flag("message").Changed {
			final, tail := encode_advanced()
			if cmd.Flag("output").Changed {
				os.Remove(cmd.Flag("output").Value.String())
				file, _ := os.OpenFile(cmd.Flag("output").Value.String(), os.O_WRONLY|os.O_CREATE, 755)
				file.WriteString(final + tail)
			}
		} else {
			decode_advanced()
		}
	},
}

func init() {
	stegoCmd.AddCommand(homoCmd)
	homoCmd.PersistentFlags().String("mapping", "", "custom char mapping file")
	homoCmd.PersistentFlags().StringVar(&text, "input", "", "input text(to decode/encode)")
	homoCmd.PersistentFlags().StringVar(&secret, "message", "", "message to hidden")
	homoCmd.PersistentFlags().StringP("output", "o", "", "output path")
}

func initMap(cmd *cobra.Command) {
	var split []string
	if cmd.Flag("mapping").Changed {
		bytes, err := os.ReadFile(cmd.Flag("mapping").Value.String())
		if err != nil {
			panic(err)
		}
		split = utils.Lines(string(bytes))
	} else {
		split = utils.Lines(defaultMap)
	}
	for _, s := range split {
		if len(s) == 0 {
			continue
		}
		var str = ""
		for _, c := range s[3:] {
			if c == 0xFFFD {
				continue
			}
			str += fmt.Sprintf("%c", c)
		}
		mapping[int32(s[0])] = utils.SplitToUnicodeChars(str)
	}
}
