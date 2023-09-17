/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
)

var resizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "A brief description of your command",
	Long:  `set height or width to 0 will keep aspect`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("input").Changed {
			fmt.Printf("input not specified\n")
			cmd.Usage()
			return
		}
		if !cmd.Flag("output").Changed {
			println("output not specified")
			cmd.Usage()
			return
		}
		img, err := readImage(cmd)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		// 调整图片尺寸
		var resizedImg image.Image
		switch algorithm {
		case "NearestNeighbor":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
		case "Bilinear":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.Bilinear)
		case "Bicubic":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.Bicubic)
		case "MitchellNetravali":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.MitchellNetravali)
		case "Lanczos2":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.Lanczos2)
		case "Lanczos3":
			resizedImg = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
		default:
			fmt.Printf("不支持的算法\n\t%s\n", cmd.Flag("algorithm").Usage)
			return
		}
		outputFile, err := os.Create(cmd.Flag("output").Value.String())
		if err != nil {
			fmt.Println("无法创建输出图片文件:", err)
			return
		}
		defer outputFile.Close()

		// 根据图片格式进行编码
		switch format {
		case "jpeg":
			jpeg.Encode(outputFile, resizedImg, nil)
		case "png":
			png.Encode(outputFile, resizedImg)
		default:
			fmt.Println("不支持的图片格式")
		}
		fmt.Println("图片压缩完成")
	},
}
var width int
var height int
var format string
var algorithm string

func init() {
	imageCmd.AddCommand(resizeCmd)
	resizeCmd.Flags().IntVarP(&width, "width", "w", 0, "output width")
	resizeCmd.Flags().IntVarP(&height, "height", "h", 0, "output height")
	resizeCmd.Flags().IntVarP(&width, "width", "W", 0, "output width")
	resizeCmd.Flags().IntVarP(&height, "height", "H", 0, "output height")
	resizeCmd.Flags().StringVarP(&format, "format", "f", "jpg", "output format (default: jpg)")
	resizeCmd.Flags().StringVarP(&algorithm, "algorithm", "g", "Lanczos3", "algorithm support:\n\t// Nearest-neighbor interpolation\n\tNearestNeighbor InterpolationFunction = iota\n\t// Bilinear interpolation\n\tBilinear\n\t// Bicubic interpolation (with cubic hermite spline)\n\tBicubic\n\t// Mitchell-Netravali interpolation\n\tMitchellNetravali\n\t// Lanczos interpolation (a=2)\n\tLanczos2\n\t// Lanczos interpolation (a=3)\n\tLanczos3 (default)")
}
