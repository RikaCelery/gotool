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

var base64DecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "base64 decoder",
	Long:  `base64 decoder`,
	Run: func(cmd *cobra.Command, args []string) {
		inputString := ""
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
		if cmd.Flag(flagInputFile).Changed {
			file, err := os.OpenFile(cmd.Flag(flagInputFile).Value.String(), os.O_RDONLY, 777)
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer file.Close()
			strBytes, _ := io.ReadAll(file)
			inputString = string(strBytes)
		} else if cmd.Flag(flagData).Changed {
			inputString = cmd.Flag(flagData).Value.String()
		}

		var decodeResult []byte
		var err error
		switch _base {
		case 32:
			if useCustomEncoder {
				decodeResult, err = base32.NewEncoding(customEncoder).DecodeString(inputString)
			} else if _hex {
				decodeResult, err = base32.HexEncoding.DecodeString(inputString)
			} else {
				decodeResult, err = base32.StdEncoding.DecodeString(inputString)
			}
			break
		case 64:
			if useCustomEncoder {
				decodeResult, err = base64.NewEncoding(customEncoder).DecodeString(inputString)
			} else if _url && _raw {
				decodeResult, err = base64.RawURLEncoding.DecodeString(inputString)
			} else if _url && !_raw {
				decodeResult, err = base64.URLEncoding.DecodeString(inputString)
			} else if !_url && _raw {
				decodeResult, err = base64.RawStdEncoding.DecodeString(inputString)
			} else if !_url && !_raw {
				decodeResult, err = base64.StdEncoding.DecodeString(inputString)
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
				panic(err)
				log.Fatalln(err)
				return
			}
			defer file.Close()
			if cmd.Flag(flagHexOutput).Changed {
				file.WriteString(hex2.EncodeToString(decodeResult))
			} else {
				file.Write(decodeResult)
			}
		} else {
			if cmd.Flag(flagHexOutput).Changed {
				println(hex2.EncodeToString(decodeResult))
			} else {
				println(string(decodeResult))
			}
		}

	},
}

func init() {
	base64Cmd.AddCommand(base64DecodeCmd)
	base64DecodeCmd.Flags().Bool(hex, false, "hexEncoding")
	base64DecodeCmd.Flags().Bool(url, false, "urlEncoding")
	base64DecodeCmd.Flags().Bool(raw, false, "no padding, (base 64 only)")
}
