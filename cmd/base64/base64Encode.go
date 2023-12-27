/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package base64

import (
	"encoding/base32"
	"encoding/base64"
	hex2 "encoding/hex"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var base64EncodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "base64 encoder",
	Long:  `base64 encoder`,
	Run: func(cmd *cobra.Command, args []string) {
		_hex, _ := strconv.ParseBool(cmd.Flag(hex).Value.String())
		_url, _ := strconv.ParseBool(cmd.Flag(url).Value.String())
		_raw, _ := strconv.ParseBool(cmd.Flag(raw).Value.String())
		useCustomEncoder := cmd.Flag(flagEncoder).Changed
		customEncoder := cmd.Flag(flagEncoder).Value.String()
		_base, _ := strconv.ParseInt(cmd.Flag(flagBase).Value.String(), 10, 32)
		// check flags
		{
			if _base != 32 && _base != 64 {
				log.Fatalf("invalid bsse %v\n", cmd.Flag(flagBase).Value.String())
				cmd.Usage()
				return
			}
			if !cmd.Flag(flagInputFile).Changed && !cmd.Flag(flagData).Changed {
				log.Println("no input")
				cmd.Usage()
				return
			}
			if _base == 32 && _raw {
				log.Println("rawCoding is not invalid in base64")
			}
			if _base == 32 && _url {
				log.Println("urlCoding is not invalid in base64")
			}
			if _base == 64 && _hex {
				log.Println("hexCoding is not invalid in base64")
			}
		}
		var inputBytes []byte
		var err error
		if cmd.Flag(flagInputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(flagInputFile).Value.String(), os.O_RDONLY, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer file.Close()
			content, _ := io.ReadAll(file)
			if cmd.Flag(hexInput).Changed {
				inputBytes, err = hex2.DecodeString(string(content))
			} else {
				inputBytes = content
			}
		} else if cmd.Flag(flagData).Changed {
			if cmd.Flag(hexInput).Changed {
				inputBytes, err = hex2.DecodeString(cmd.Flag(flagData).Value.String())
			} else {
				inputBytes = []byte(cmd.Flag(flagData).Value.String())
			}
		}
		if err != nil {
			log.Fatalln(err)
			return
		}

		var encodeResult string
		switch _base {
		case 32:
			if useCustomEncoder {
				encodeResult = base32.NewEncoding(customEncoder).EncodeToString(inputBytes)
			} else if _hex {
				encodeResult = base32.HexEncoding.EncodeToString(inputBytes)
			} else {
				encodeResult = base32.StdEncoding.EncodeToString(inputBytes)
			}
			break
		case 64:
			if useCustomEncoder {
				encodeResult = base64.NewEncoding(customEncoder).EncodeToString(inputBytes)
			} else if _url && _raw {
				encodeResult = base64.RawURLEncoding.EncodeToString(inputBytes)
			} else if _url && !_raw {
				encodeResult = base64.URLEncoding.EncodeToString(inputBytes)
			} else if !_url && _raw {
				encodeResult = base64.RawStdEncoding.EncodeToString(inputBytes)
			} else if !_url && !_raw {
				encodeResult = base64.StdEncoding.EncodeToString(inputBytes)
			}
			break
		}

		if err != nil {
			log.Fatalf("decode error, err=%v", err)
			return
		}
		if cmd.Flag(flagOutputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(flagOutputFile).Value.String(), os.O_WRONLY|os.O_CREATE, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer file.Close()
			if cmd.Flag(flagHexOutput).Changed {
				file.WriteString(hex2.EncodeToString([]byte(encodeResult)))
			} else {
				file.WriteString(encodeResult)
			}
		} else {
			if cmd.Flag(flagHexOutput).Changed {
				println(hex2.EncodeToString([]byte(encodeResult)))
			} else {
				println(encodeResult)
			}
		}

	},
}

func init() {
	base64Cmd.AddCommand(base64EncodeCmd)
	base64EncodeCmd.Flags().Bool(hexInput, false, "process input as hex string")
	base64EncodeCmd.Flags().Bool(hex, false, "hexEncoding")
	base64EncodeCmd.Flags().Bool(url, false, "urlEncoding")
	base64EncodeCmd.Flags().Bool(raw, false, "no padding, (base 64 only)")
}
