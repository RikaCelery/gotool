/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "gotool/cmd"
import _ "gotool/cmd/base64"
import _ "gotool/cmd/hex"
import _ "gotool/cmd/reverse"
import _ "gotool/cmd/image"
import _ "gotool/cmd/digit_converter"
import _ "gotool/cmd/sort"
import _ "gotool/cmd/setgo"

func main() {
	cmd.Execute()
}
