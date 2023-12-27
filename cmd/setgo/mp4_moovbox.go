package setgo

import (
	"fmt"
	"os"
	"strings"

	"github.com/abema/go-mp4"
	"github.com/joe-xu/mp4parser"
	"github.com/spf13/cobra"
	"gotool/utils"
)

var moovStegoCmd = &cobra.Command{
	Use:   "moov",
	Short: "",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("C:\\Users\\celery\\11.mp4")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		parser := mp4parser.NewParser(file)
		info, err := parser.Parse()
		if err != nil {
			fmt.Println("Error parsing MP4:", err)
			return
		}
		println(info.String())
		println("-----------------")
		// return
		structures, err := mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
			if len(h.Path) > 1 {
				return nil, nil
			}
			var indent = strings.Repeat(" ", len(h.Path)*2)
			// Box Type (e.g. "mdhd", "tfdt", "mdat")
			fmt.Printf("%sdepth %v,type[%v]size:%v\n", indent, len(h.Path), h.BoxInfo.Type.String(), h.BoxInfo.Size)

			if h.BoxInfo.IsSupportedType() {
				// Payload
				box, _, err := h.ReadPayload()
				if err != nil {
					return nil, err
				}
				fmt.Printf("%v\n",box.GetFlags())
				str, err := mp4.Stringify(box, h.BoxInfo.Context)
				if err != nil {
					return nil, err
				}
				fmt.Printf("%spayload %v\n", indent, utils.ATruncate(str, 200))
				// Expands children
				return h.Expand()
			}
			return nil, nil
		})
		if err != nil {
			return
			println(structures[0])
		}
	},
}

// func init() {
// 	stegoCmd.AddCommand(moovStegoCmd)
// }
