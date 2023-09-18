/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package reverse

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
		files := args
		// make sure all files exist
		if invalids := checkFiles(files); len(invalids) != 0 {
			fmt.Printf("files invalid:\n%v\n", strings.Join(invalids, "\n"))
			return
		}

		startTime := time.Now().UnixMilli()
		for i := 0; i < len(files); i++ {

			if filepath.Ext(files[i]) == ".reverse" {
				decode(files[i], bit, base, !encodeName, force)
			} else {
				encode(files[i], bit, base, !encodeName, force)
			}

		}
		fmt.Println("everything done, time cost", (time.Now().UnixMilli()-startTime)/1000, "s")
	},
}
var bit int
var base int
var force bool
var encodeName bool

func init() {
	cmd.RootCmd.AddCommand(reverseCmd)
	// todo

	reverseCmd.Flags().IntVar(&bit, "bit", 1280, "block size")
	reverseCmd.Flags().IntVar(&base, "base", 32, "base32(32)/base64(64)")
	reverseCmd.Flags().BoolVar(&force, "force", false, "overwrite file")
	reverseCmd.Flags().BoolVar(&encodeName, "enc_name", false, "do not encode/decode file name")
	reverseCmd.Flags().BoolVar(&deleteOrigin, "delete", false, "delete origin file")
}

var cyan, green, red, yellow = color.New(color.FgHiCyan), color.New(color.FgGreen), color.New(color.FgHiRed), color.New(color.FgHiYellow)
var deleteOrigin bool

func checkFiles(files []string) (invalids []string) {
	for _, file := range files {
		if !isExist(file) {
			invalids = append(invalids, file)

		}
	}
	return invalids
}

const (
	DEBUG           = false
	BUFFER_SIZE     = 1024 * 1024 * 10
	OUT_BUFFER_SIZE = 1024 * 1024 * 10
)

var MagicBytes = []byte{0x7c, 0xc2, 0x32, 0x00, 0xff}

type MetaData struct {
	Bit        int64
	CreateDate int64
	NoFileName bool
	OriginName string
}

func readMeta(file *os.File) (MetaData, error) {
	magic := make([]byte, 5)
	file.Read(magic)
	if bytes.Equal(MagicBytes, magic) {
	} else {
		if DEBUG {
			fmt.Printf("%v\n", magic)
		}
		return MetaData{}, errors.New("magic Miss Match")
	}

	in := make([]byte, binary.MaxVarintLen64)
	file.Read(in)
	bit, _ := binary.Varint(in)
	in = make([]byte, binary.MaxVarintLen64)
	file.Read(in)
	createDate, _ := binary.Varint(in)
	in = make([]byte, 1)
	file.Read(in)
	noFileName := false
	if in[0] == byte(255) {
		noFileName = true
	}
	var name string
	if !noFileName {
		in = make([]byte, binary.MaxVarintLen64)
		file.Read(in)

		length, _ := binary.Varint(in)
		in = make([]byte, length)
		file.Read(in)
		name1, err := decodeFromUTF16(in)
		if err != nil {
			return MetaData{}, errors.New(err.Error())
		}
		name = name1

		in = make([]byte, 522-length)
		file.Read(in)
	}
	return MetaData{
		Bit:        bit,
		CreateDate: createDate,
		NoFileName: noFileName,
		OriginName: name,
	}, nil
}
func decodeFromUTF16(input []byte) (string, error) {
	decoder := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewDecoder()
	encoded, _, err := transform.Bytes(decoder, input)

	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
func encodeToUTF16(s string) ([]byte, int, error) {
	encoder := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewEncoder()
	encoded, _, err := transform.Bytes(encoder, []byte(s))
	if err != nil {
		return nil, 0, err
	}
	return encoded, len(encoded), nil
}
func (d *MetaData) WriteMeta(file *os.File) error {
	file.Write(MagicBytes)

	byteSlice := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(byteSlice, d.Bit)
	file.Write(byteSlice)

	byteSlice = make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(byteSlice, d.CreateDate)
	file.Write(byteSlice)

	byteSlice = make([]byte, 1)
	file.Write(byteSlice)

	// if !d.NoFileName {
	encoded, n, err := encodeToUTF16(d.OriginName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	byteSlice = make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(byteSlice, int64(len(encoded)))
	file.Write(byteSlice)
	if n > 522 {
		return errors.New("name too long")
	}
	// 计算填充的字节数
	paddingLength := 522 - len(encoded)

	// 如果需要填充，则将填充添加到字节切片末尾
	if paddingLength > 0 {
		padBytes := make([]byte, paddingLength)
		encoded = append(encoded, padBytes...)
	}
	file.Write(encoded)
	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func removeExtension(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
func encode(inputFileName string, bit int, base int, keepName bool, overwrite bool) {
	println("encode:", inputFileName)
	inputFile, _ := os.Open(inputFileName)
	var byteOffset int64 = 0

	dir, fileName := filepath.Split(inputFileName)
	ext := filepath.Ext(inputFileName)
	fileName = removeExtension(fileName)
	if !keepName {
		fileName = baseEncode(fileName, base)
	}
	if len(ext) > 0 {
		fileName += ext
	}
	outputFileName := filepath.Join(dir, fileName+".reverse")
	tempFileName := outputFileName + ".reverse.tmp"
	println("======>", outputFileName)

	info, _ := inputFile.Stat()
	if isExist(outputFileName) {
		stat, _ := os.Stat(outputFileName)
		if !overwrite && stat.Size() == info.Size()+558 {
			fmt.Printf("\texist, skip(optput:%v,original:%v/diff:%v)\n", stat.Size(), info.Size(), stat.Size()-info.Size())
			return
		} else {
			yellow.Printf("\toverwrite, (optput:%v,original:%v/diff:%v)\n", stat.Size(), info.Size(), stat.Size()-info.Size())
		}
	}
	_ = os.Remove(tempFileName)
	metaData := MetaData{
		Bit:        int64(bit),
		CreateDate: info.ModTime().UnixMilli(),
		NoFileName: keepName,
		OriginName: info.Name(),
	}

	tempFile, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE, 755)
	if err != nil {
		red.Println("[err decode] failed to open target inputFileName(", tempFileName, "),err=", err.Error())
		return
	}
	err = metaData.WriteMeta(tempFile)
	if err != nil {
		red.Printf("err%v\n", err.Error())
		return
	}

	err = reverseData(inputFile, info.Size()-byteOffset, bit, tempFile)
	inputFile.Close()
	tempFile.Close()
	if err != nil {
		red.Println("error!!", inputFileName+"[", err.Error(), "]")
		_ = os.Remove(tempFileName)
		return
	} else {
		err := os.Rename(tempFileName, outputFileName)
		if err != nil {
			red.Printf("rename failed %v\n", err.Error())
			return
		}
	}
	if deleteOrigin {
		err := os.Remove(inputFile.Name())
		if err != nil {
			fmt.Printf("[error] file %v cannot be removed err=%v", inputFile.Name(), err)
		}
	}
}

func decode(inputFileName string, bit int, base int, keepName bool, overwrite bool) {
	fmt.Printf("+---decode---<< %v\n", inputFileName)

	var inputFile, _ = os.Open(inputFileName)
	var byteOffset int64 = 558
	meta, err := readMeta(inputFile)
	inputInfo, _ := inputFile.Stat()
	dir, fileName := filepath.Split(inputFileName)
	var legacy = false
	if err == nil && !meta.NoFileName {
		fileName = meta.OriginName
		bit = int(meta.Bit)
	} else {
		legacy = true
		yellow.Printf("| error in reading meta use legacy mode(%v)\n", err.Error())
		split := strings.Split(removeExtension(fileName), ".")
		if !keepName {
			s, err2 := baseDecode(split[0], base)
			if err2 != nil {
				s, err2 = baseDecode(split[0], 64)
			}
			if err2 == nil {
				split[0] = s
			}
		}
		fileName = strings.Join(split, ".")
	}
	outputFileName := filepath.Join(dir, fileName)
	tempFileName := outputFileName + ".tmp"
	// fmt.Printf("| tmp:%v\n", tempFileName)
	fmt.Printf("| bit:%v\n", bit)
	fmt.Printf("| name:%v\n", fileName)
	if legacy {
		fmt.Printf("| [legacy mode]\n")
		byteOffset = 0
		_, err := inputFile.Seek(0, 0)
		if err != nil {
			red.Printf("seek error %v", err)
			return
		}
	} else {
		fmt.Printf("| createDate:%v\n", time.UnixMilli(meta.CreateDate).String())
		// fmt.Printf("| noFileName%v\n", meta.NoFileName)

	}
	fmt.Println("+------------>>", outputFileName)
	if isExist(outputFileName) {
		output, _ := os.Stat(outputFileName)
		if !overwrite && (output.Size() == inputInfo.Size()-byteOffset || output.Size() == inputInfo.Size()) {
			yellow.Printf("\texist, skip(output:%v,original:%v/diff:%v)\n", output.Size(), inputInfo.Size(), output.Size()-inputInfo.Size())
			return
		} else {
			yellow.Printf("\toverwrite, (output:%v,original:%v/diff:%v)\n", output.Size(), inputInfo.Size(), output.Size()-inputInfo.Size())
		}
	}

	_ = os.Remove(tempFileName)
	tempFile, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE, 755)
	if err != nil {
		red.Println("[err decode] failed to open target inputFileName(", tempFileName, "),err=", err.Error())
		return
	}
	err = reverseData(inputFile, inputInfo.Size()-byteOffset, bit, tempFile)
	inputFile.Close()
	tempFile.Close()
	if err != nil {
		red.Println("error!!", inputFileName+"[", err.Error(), "]")
		_ = os.Remove(tempFileName)
		return
	} else {
		err := os.Rename(tempFileName, outputFileName)
		if err != nil {
			red.Printf("rename failed %v\n", err.Error())
			return
		}
	}
	if deleteOrigin {
		err := os.Remove(inputFile.Name())
		if err != nil {
			fmt.Printf("[error] file %v cannot be removed err=%v", inputFile.Name(), err)
		}
	}
}

type ProgressInfo struct {
	Time      time.Time
	Processed int64
	Total     int64
}

// 获取终端窗口大小的函数
func getTerminalSize() (width, height int, err error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "mode", "con")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		out, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		lines := strings.Split(string(out), "\n")
		if len(lines) >= 4 {
			// 第三行包含宽度和高度信息，格式如：Columns: 80  Height: 24
			columnsLine := strings.TrimSpace(lines[2])
			heightLine := strings.TrimSpace(lines[3])

			columnsParts := strings.Split(columnsLine, ":")
			heightParts := strings.Split(heightLine, ":")

			if len(columnsParts) == 2 && len(heightParts) == 2 {
				width, _ = strconv.Atoi(strings.TrimSpace(columnsParts[1]))
				height, _ = strconv.Atoi(strings.TrimSpace(heightParts[1]))
				return width, height, nil
			}
		}
	case "linux", "darwin":
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		out, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		sizeParts := strings.Split(strings.TrimSpace(string(out)), " ")
		if len(sizeParts) == 2 {
			width, _ = strconv.Atoi(sizeParts[1])
			height, _ = strconv.Atoi(sizeParts[0])
			return width, height, nil
		}
	}

	return 0, 0, fmt.Errorf("无法获取终端窗口大小")
}

func toMetaFunc(c *color.Color) func(string) string {
	return func(s string) string {
		return c.Sprint(s)
	}
}

func progressRoutine(done <-chan struct{}, progressCh <-chan ProgressInfo, wg *sync.WaitGroup, length int64, name string) {
	wg.Add(1)
	defer wg.Done()

	var window []ProgressInfo
	var currentProgress ProgressInfo
	var p = mpb.New(
		mpb.WithOutput(color.Output),
		mpb.WithAutoRefresh(),
		mpb.WithRefreshRate(100*time.Millisecond),
	)
	var start = time.Now()
	width, _, err := getTerminalSize()
	width /= 5
	if err != nil || width > 40 {
		width = 40
	}
	if runewidth.StringWidth(name) < width {
		width = runewidth.StringWidth(name) + 1
	}
	cost := ""
	bar := p.New(length,
		mpb.BarStyle().Lbound("").Filler("=").Tip("▶").Padding(" ").Rbound(""),
		mpb.PrependDecorators(
			decor.OnComplete(decor.Spinner([]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}, decor.WC{W: 2, C: decor.DidentRight}), "✅  "),
			decor.OnAbort(decor.Name(""), "❌  "),
			decor.OnCompleteMeta(
				decor.OnComplete(decor.Meta(Marquee(name, width, "   "), toMetaFunc(cyan)), name),
				toMetaFunc(green),
			),
			decor.OnComplete(
				decor.Name("["), "",
			),
			decor.OnComplete(
				decor.Any(func(statistics decor.Statistics) string {
					var windowProcessed int64 = 0
					var windowTime = (2.0 * time.Second)
					if len(window) != 0 {
						windowProcessed = window[len(window)-1].Processed - window[0].Processed
						windowTime = window[len(window)-1].Time.Sub(window[0].Time)
					}
					speed := float64(windowProcessed) / windowTime.Seconds() // per second
					speedS := fmt.Sprintf("%.1f", decor.FmtAsSpeed(decor.SizeB1024(math.Round(speed))))
					return speedS
				}),
				"",
			),
			decor.OnComplete(
				decor.Name("] "), "",
			),
			decor.OnComplete(
				decor.CountersKibiByte("%.1f/%.1f"), "",
			),
		),
		mpb.AppendDecorators(decor.Percentage()),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Any(func(statistics decor.Statistics) string {
				if statistics.Completed {
					return ""
				}
				var windowProcessed int64 = 0
				var windowTime = (2.0 * time.Second)
				if len(window) != 0 {
					windowProcessed = window[len(window)-1].Processed - window[0].Processed
					windowTime = window[len(window)-1].Time.Sub(window[0].Time)
				}
				speed := float64(windowProcessed) / windowTime.Seconds() // per second
				remain := float64(currentProgress.Total-currentProgress.Processed) / speed
				reaminS := fmt.Sprintf("[%.1fs]", remain)
				return reaminS
			}),
			decor.Name("["),
			decor.Any(func(statistics decor.Statistics) string {
				if statistics.Completed {
					return cost
				}
				cost = fmt.Sprintf("%.1fs", time.Now().Sub(start).Seconds())
				return cost
			}),
			decor.Name("]"),
		),
		// mpb.BarRemoveOnComplete(),
	)

	update := start
	for {
		select {
		case info := <-progressCh:
			bar.EwmaIncrBy(int(info.Processed-currentProgress.Processed), time.Now().Sub(update))
			currentProgress = info
			now := time.Now()
			update = now
			window = append(window, info)
			// Remove data older than 2 seconds from the window
			for len(window) > 0 && now.Sub(window[0].Time) > 3*time.Second {
				window = window[1:]
			}
		case <-done:
			p.Wait()
			return
		}
	}
}

func reverseData(input *os.File, length int64, bit int, output *os.File) error {
	var processed int64 = 0
	inputBuffer := bytes.Buffer{}
	if bit > BUFFER_SIZE {
		inputBuffer.Grow(bit)
	} else {
		inputBuffer.Grow(BUFFER_SIZE)
	}
	read := 0
	var b []byte
	b = make([]byte, BUFFER_SIZE)
	outBuffer := bytes.Buffer{}
	outBuffer.Grow(OUT_BUFFER_SIZE) // 100MB
	var loopError error = nil

	writeBuffer := func() {
		_, err := output.Write(outBuffer.Bytes())
		if err != nil {
			red.Println(err.Error())
			loopError = err
		}
		outBuffer.Reset()
	}

	done := make(chan struct{})
	progressCh := make(chan ProgressInfo)
	var wg sync.WaitGroup
	// Start the progress routine
	go progressRoutine(done, progressCh, &wg, length, filepath.Base(input.Name()))

	for true {
		// time.Sleep(1 * time.Millisecond)
		if read > bit { // 读了多个block进来
			// 取一个block
			block := inputBuffer.Next(bit)
			read -= len(block)

			// 下面要写oBuffer，如果oBuffer满了，先写入文件
			if outBuffer.Len()+len(block) > outBuffer.Cap() {
				writeBuffer()
			}

			// 将block反转填入oBuffer
			processBlock(block, &outBuffer)
			processed += int64(len(block))

			// progressing
			progressCh <- ProgressInfo{
				Time:      time.Now(),
				Processed: processed,
				Total:     length,
			}
			continue
		}
		n, err := input.Read(b)
		if err != nil {
			if err != io.EOF {
				loopError = err // error
				break
			} else if read <= bit {
				// reach end
				// 上面处理完之后 inputBuffer一定只剩下最后一块
				block := inputBuffer.Next(bit)
				if outBuffer.Len()+len(block) > outBuffer.Cap() {
					writeBuffer()
				}
				processBlock(block, &outBuffer)
				writeBuffer()
				processed += int64(len(block))

				// progressing
				progressCh <- ProgressInfo{
					Time:      time.Now(),
					Processed: processed,
					Total:     length,
				}
				break
			}
		} else {
			// normal
			read += n
			inputBuffer.Write(b[:n])
		}
		// progressing
		progressCh <- ProgressInfo{
			Time:      time.Now(),
			Processed: processed,
			Total:     length,
		}
	}

	// Signal the progress routine to stop and wait for it to finish
	close(done)
	wg.Wait()
	if loopError != nil {
		red.Printf("error%v\n", loopError)
	} else if processed != length {
		loopError = errors.New("file not processed fully")
		red.Printf("error %d/%d \n", processed, length)
	}
	println()
	return loopError
}

func processBlock(block []byte, file *bytes.Buffer) {
	for i, j := 0, len(block)-1; i < j; i, j = i+1, j-1 {
		block[i], block[j] = block[j], block[i] // reverse the slice
	}
	_, err := file.Write(block)

	if err != nil {
		red.Println("error", err.Error())
		return
	}
}

func baseEncode(str string, base int) string {
	ret := str
	if base == 32 {
		encoding := base32.StdEncoding
		ret = encoding.EncodeToString([]byte(str))
	} else if base == 64 {
		encoding := base64.URLEncoding
		ret = encoding.EncodeToString([]byte(str))
	}
	return ret
}
func baseDecode(str string, base int) (string, error) {
	if base == 32 {
		encoding := base32.StdEncoding
		ret, err := encoding.DecodeString(str)
		if err != nil {
			return "", err
		}
		return string(ret), nil
	} else if base == 64 {
		encoding := base64.URLEncoding
		ret, err := encoding.DecodeString(str)
		if err != nil {
			return "", err
		}
		return string(ret), nil
	}
	return str, nil
}
