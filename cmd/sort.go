/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"gotool/utils"
)

var (
	rootFolder = ""
	inputFiles = []string{}
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "sort files into category",
	Long:  `sort files into category`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		if len(args) > 1 {
			inputFiles = args
		}
		if len(inputFiles) != 0 {

			if len(rootFolder) == 0 {
				fmt.Println("root folder can not be empty")
				cmd.Usage()
				return
			}
			// 遍历子项
			for _, path := range inputFiles {
				if utils.IsDir(path) {
					// 忽略子目录
					continue
				}
				if utils.Contains(Ignores, filepath.Base(path)) {
					continue
				}
				err := sortFile(filepath.Join(args[0], path))
				if err != nil {
					if os.IsExist(err) {
						fmt.Printf("File exist %v", err)
					} else if os.IsPermission(err) {
						fmt.Printf("permission denied %v", err)
					} else {
						fmt.Printf("%v", err)
					}
				}
			}
			return
		}
		if len(rootFolder) == 0 {
			rootFolder = args[0]
		}
		// 获取目录下的所有子项（文件和目录）
		items, err := os.ReadDir(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 遍历子项
		for _, item := range items {
			if item.IsDir() {
				// 忽略子目录
				continue
			}
			if utils.Contains(Ignores, item.Name()) {
				continue
			}
			err := sortFile(filepath.Join(args[0], item.Name()))
			if err != nil {
				if os.IsExist(err) {
					fmt.Printf("File exist %v\n", err)
				} else if os.IsPermission(err) {
					fmt.Printf("permission denied %v\n", err)
				} else {
					fmt.Printf("%v\n", err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)
	sortCmd.Example = `  gotool sort path/to/your/dir
  gotool sort path/to/your/dir -D path/destination
  gotool sort --files file1 file2 file3 -D path/destination`
	sortCmd.Flags().StringArrayVarP(&inputFiles, "files", "F", []string{}, "input files")
	sortCmd.Flags().StringVarP(&rootFolder, "dest", "D", "", "sort dest")
}

func sortFile(path string) error {
	historyFile, err := os.OpenFile(filepath.Join(rootFolder, "sort_history"), os.O_CREATE|os.O_APPEND, 777)
	if err != nil {
		return err
	}
	_, err = historyFile.WriteString(time.Now().Format("2006-01-02 15:04:05\n"))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	dir, err := getDestDir(path)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return nil
	}

	if !utils.IsExist(filepath.Join(rootFolder, dir)) {
		err := os.MkdirAll(filepath.Join(rootFolder, dir), 777)
		if err != nil {
			return err
		}
	}
	err = os.Rename(path, filepath.Join(rootFolder, dir, filepath.Base(path)))
	if err != nil {
		return err

	}
	fmt.Printf("%s >> %s\n", path, dir)
	_, err = historyFile.WriteString(fmt.Sprintf("%s >> %s\n", path, dir))
	if err != nil {
		return err
	}
	return nil
}
func getDestDir(path string) (string, error) {
	if anyMatch(Music, path) {
		return "Music", nil
	}
	if anyMatch(Video, path) {
		return "Video", nil
	}
	if anyMatch(Document, path) {
		return "Document", nil
	}
	if anyMatch(Archive, path) {
		return "Archive", nil
	}
	if anyMatch(Book, path) {
		return "Book", nil
	}
	if anyMatch(Code, path) {
		return "Code", nil
	}
	if anyMatch(Photo, path) {
		return "Photo", nil
	}
	if anyMatch(Video, path) {
		return "Video", nil
	}
	if anyMatch(Program, path) {
		return "Program", nil
	}
	if anyMatch(AdobeScripts, path) {
		return "Adobe Scripts", nil
	}
	return "", errors.New("unknown extension " + path)
}

func anyMatch(slice []Match, item string) bool {
	for _, value := range slice {
		if value.GetMatch(item) {
			return true
		}
	}
	return false
}

var Ignores = []string{
	"sort.ps1",
	"sort_history",
	"Document",
}

var Music = extMatches([]string{"mp3", "m4a", "mp2", "flac", "wav", "WAV"})
var Video = extMatches([]string{"mp4", "mov", "flv"})
var Document = extMatches([]string{"txt", "pdf", "ppt", "pptx", "doc", "docx", "xls", "xlsx", "mhtml", "ini", "json", "xml"})
var Photo = extMatches([]string{"jpg", "png", "tif"})
var Archive = append(extMatches([]string{"zip", "7z", "rar", "gz", "xz"}),
	&RegexMatch{regex: regexp.MustCompile(`.+\.7z\.\d+$`)},
	&RegexMatch{regex: regexp.MustCompile(`.+\.zip\.\d+$`)},
	&RegexMatch{regex: regexp.MustCompile(`.+\.rar\.\d+$`)},
	&RegexMatch{regex: regexp.MustCompile(`.+\.gz\.\d+$`)},
	&RegexMatch{regex: regexp.MustCompile(`.+\.xz\.\d+$`)},
)
var Book = extMatches([]string{"epub", "mobi", "azw3", "azw2"})
var Code = extMatches([]string{"c", "cpp", "java", "js", "kts", "py", "sql", "html", "htm"})
var Program = extMatches([]string{"exe", "msi", "jar"})
var AdobeScripts = extMatches([]string{"jsxbin"})

type Match interface {
	GetMatch(name string) bool
}

type RegexMatch struct {
	regex *regexp.Regexp
}

func (m *RegexMatch) GetMatch(name string) bool {
	return m.regex.MatchString(name)
}

type ExtMatch struct {
	ext string
}

func (m *ExtMatch) GetMatch(name string) bool {
	ext := filepath.Ext(name)
	if len(ext) == 0 {
		return false
	}
	return m.ext == ext[1:] // remove dot
}

func extMatches(exts []string) []Match {
	matches := make([]Match, 0)
	for i := range exts {
		matches = append(matches, &ExtMatch{ext: exts[i]})
	}
	return matches
}
