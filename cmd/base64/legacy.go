package base64

import (
	"encoding/base64"
	lib_hex "encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	_ "golang.org/x/text/encoding/unicode"
	_ "golang.org/x/text/transform"
)

type Option struct {
	stdout     bool
	hex        bool
	outputFile *os.File
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
func outputStringResult(out []byte, option *Option) {
	if option.stdout {
		if option.hex {
			fmt.Printf("hex:\n%v\n", lib_hex.EncodeToString(out))
		} else {
			fmt.Printf("%v\n", string(out))
		}
	} else {
		_, err := option.outputFile.Write(out)
		if err != nil {
			println(err.Error())
		}
		option.outputFile.Close()
	}
}
func outputBytesResult(out []byte, option *Option) {
	if option.stdout {
		fmt.Printf("hex:\n%v\n", lib_hex.EncodeToString(out))
	} else {
		option.outputFile.Write(out)
		option.outputFile.Close()
	}
}
func main() {
	encodeMode := flag.Bool("encode", false, "encode mode, default will auto detect")
	decodeMode := flag.Bool("decode", false, "encode mode, default will auto detect")
	hexOut := flag.Bool("hex", false, "convert to hex string")
	outputFileName := flag.String("out", "", "output file (empty use stdout)")
	help := flag.Bool("help", false, "print help")

	flag.Parse()
	if len(flag.Args()[1:]) == 0 || *help {
		println("** [Legacy Command] **")
		println("Usage: [Options] input(file path/string)")
		println("Options:")
		flag.PrintDefaults()
		return
	}
	input := flag.Args()[1]

	var outputFile *os.File = nil

	if len(*outputFileName) != 0 {
		var err error
		if !isExist(*outputFileName) {
			file, err := os.OpenFile(*outputFileName, os.O_CREATE, 777)
			if err != nil {
				return
			}
			outputFile = file
		} else {
			outputFile, err = os.Open(*outputFileName)
			if err != nil {
				println(err.Error())
				return
			}
		}
	}
	encoder := base64.StdEncoding
	option := &Option{
		stdout:     len(*outputFileName) == 0,
		hex:        *hexOut,
		outputFile: outputFile,
	}
	if isExist(input) {
		bytes, _ := fileInput(input)
		if utf8.Valid(bytes) && isBase64(string(bytes)) && !*encodeMode || *decodeMode {
			bytes, _ := encoder.DecodeString(string(bytes))
			if utf8.Valid(bytes) {
				outputStringResult(bytes, option)
			} else {
				outputBytesResult(bytes, option)
			}

		} else {
			str := encoder.EncodeToString(bytes)
			outputStringResult([]byte(str), option)
		}
	} else {
		if isBase64(input) && !*encodeMode || *decodeMode {
			bytes, _ := encoder.DecodeString(input)
			if utf8.Valid(bytes) {
				outputStringResult(bytes, option)
			} else {
				outputBytesResult(bytes, option)
			}

		} else {
			str := encoder.EncodeToString([]byte(input))
			outputStringResult([]byte(str), option)
		}
	}
}

func isBase64(input string) bool {
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil
}

func fileInput(s string) ([]byte, error) {
	file, _ := os.Open(s)
	bytes, err := io.ReadAll(file)
	if err != nil {
		return make([]byte, 0), errors.New("error")
	}
	return bytes, nil
}
