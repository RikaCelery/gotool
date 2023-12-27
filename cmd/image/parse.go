/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package image

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "jpeg header blocks parser",
	Long:  `jpeg header blocks parser`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := os.Open(cmd.Flag("input").Value.String())
		defer file.Close()
		// 读取 JPEG 文件的内容
		fileInfo, _ := file.Stat()
		fileSize := fileInfo.Size()
		buffer := make([]byte, fileSize)
		_, err := file.Read(buffer)
		if err != nil {
			fmt.Println("无法读取文件:", err)
			return
		}
		if bytes.Equal(buffer[:3], []byte{0xff, 0xd8, 0xff}) {
			// jpeg

			tagMarker := map[byte][]byte{}
			tagStart := false
			dataStart := false
			var tag byte
			for _, b := range buffer {
				if b == 0xff {
					tagStart = true
					continue
				}
				if tagStart {
					if b != 0xff && b != 0x00 {
						tag = b
						// fmt.Printf("%0x\n", b)
						if tagMarker[tag] == nil {
							tagMarker[tag] = make([]byte, 0)
						}
						tagStart = false
						dataStart = true
						continue
					} else {
						tagStart = false
						tagMarker[tag] = append(tagMarker[tag], b)
					}
				}
				if dataStart {
					tagMarker[tag] = append(tagMarker[tag], b)
				}
			}
			for b, i := range tagMarker {
				tag := fmt.Sprintf("ff%x", b)
				// big name map
				switch b {
				case 0xd8:
					tag = "SOI"
				case 0xe0:
					tag = "APP0"
				case 0xe1:
					tag = "APP1"
				case 0xe2:
					tag = "APP2"
				case 0xe3:
					tag = "APP3"
				case 0xe4:
					tag = "APP4"
				case 0xe5:
					tag = "APP5"
				case 0xe6:
					tag = "APP6"
				case 0xe7:
					tag = "APP7"
				case 0xe8:
					tag = "APP8"
				case 0xe9:
					tag = "APP9"
				case 0xea:
					tag = "APP10"
				case 0xeb:
					tag = "APP11"
				case 0xec:
					tag = "APP12"
				case 0xed:
					tag = "APP13"
				case 0xee:
					tag = "APP14"
				case 0xef:
					tag = "APP15"
				case 0xc0:
					tag = "SOF0"
				case 0xc1:
					tag = "SOF1"
				case 0xc2:
					tag = "SOF2"
				case 0xc3:
					tag = "SOF3"
				case 0xc5:
					tag = "SOF5"
				case 0xc6:
					tag = "SOF6"
				case 0xc7:
					tag = "SOF7"
				case 0xc8:
					tag = "JPG"
				case 0xc9:
					tag = "SOF9"
				case 0xca:
					tag = "SOF10"
				case 0xcb:
					tag = "SOF11"
				case 0xcd:
					tag = "SOF13"
				case 0xce:
					tag = "SOF14"
				case 0xcf:
					tag = "SOF15"
				case 0xc4:
					tag = "DHT"
				case 0xcc:
					tag = "EOI"
				case 0xd0:
					tag = "EOI"
				case 0xda:
					tag = "SOS"
				case 0xdb:
					tag = "DQT"
				case 0xdc:
					tag = "DNL"
				case 0xdd:
					tag = "DRI"
				case 0xde:
					tag = "DHP"
				case 0xdf:
					tag = "EXP"
				case 0xf0:
					tag = "JPG0"
				case 0xf1:
					tag = "JPG1"
				case 0xf2:
					tag = "JPG2"
				case 0xf3:
					tag = "JPG3"
				case 0xf4:
					tag = "JPG4"
				case 0xf5:
					tag = "JPG5"
				case 0xf6:
					tag = "JPG6"
				case 0xf7:
					tag = "JPG7"
				case 0xf8:
					tag = "JPG8"
				case 0xf9:
					tag = "JPG9"
				case 0xfa:
					tag = "JPG10"
				case 0xfb:
					tag = "JPG11"
				case 0xfc:
					tag = "JPG12"
				case 0xfd:
					tag = "JPG13"
				case 0xfe:
					tag = "COM"
				case 0x01:
					tag = "TEM"
				case 0xd9:
					tag = "EOI"
				default:
					if b >= 0x02 && b <= 0xbf {
						tag = "RES"
					}
				}
				if cmd.Flag("extract").Changed {
					if cmd.Flag("extract").Value.String() == tag {
						fmt.Printf("[%s](%x), %v\n", tag, b, i)
					}
				} else {
					fmt.Printf("[%s](%x), %v\n", tag, b, len(i))
				}
			}
		} else {
			// png
		}
		return
	},
}

func init() {
	imageCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringP("extract", "e", "", "extract specific tag")
}
