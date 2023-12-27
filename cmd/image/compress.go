/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var quality int

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "image compressor",
	Long:  `image compressor`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("input").Changed && !cmd.Flag("output").Changed {
			fmt.Printf("input or output not specified\n")
			cmd.Usage()
			return
		}
		fmt.Printf("%v >> %v, quality:%v\n", InputFile, OutputFile, quality)

		ext := filepath.Ext(InputFile)
		if !strings.EqualFold(ext, "jpg") && !strings.EqualFold(ext, "jpeg") {
			println("only support jpeg image")
			os.Exit(1)
		}

		err := compressJPEG(quality, InputFile, OutputFile)
		if err != nil {
			fmt.Println("压缩图片时出现错误:", err)
			os.Exit(1)
		}

		fmt.Println("图片压缩完成")
	},
}

func init() {
	imageCmd.AddCommand(compressCmd)
	compressCmd.Flags().IntVarP(&quality, "quality", "q", 70, "input image")
}
func compressJPEG(q int, inputFile string, outputFile string) error {
	// 打开输入图片文件
	inputImage, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("无法打开输入图片文件:", err)
		os.Exit(1)
	}
	defer inputImage.Close()
	img, _, err := image.Decode(inputImage)
	if err != nil {
		return err
	}

	// 创建输出图片文件
	outputImage, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("无法创建输出图片文件:", err)
		os.Exit(1)
	}
	defer outputImage.Close()
	options := &jpeg.Options{Quality: q} // 压缩质量，可以根据需求调整
	return jpeg.Encode(outputImage, img, options)
}
