/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package image

import (
	"image"
	"os"

	"github.com/spf13/cobra"
	"gotool/cmd"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "image utils",
	Long:  `image utils`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("image called")
	// },
}

var (
	InputFile  string
	OutputFile string
)

func readImage(cmd *cobra.Command) (image.Image, error) {
	inputFile := cmd.Flag("input").Value.String()
	inputImage, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer inputImage.Close()
	img, _, err := image.Decode(inputImage)
	inputImage.Close()
	if err != nil {
		return nil, err
	}
	return img, nil
}
func init() {
	cmd.RootCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "input image")
	imageCmd.PersistentFlags().StringVarP(&OutputFile, "output", "o", "", "input image")

}
