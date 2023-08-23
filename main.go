/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "gotool/cmd"
import _ "gotool/cmd/base64"
import _ "gotool/cmd/hex"
import _ "gotool/cmd/reverse"

func main() {
	cmd.Execute()
}
