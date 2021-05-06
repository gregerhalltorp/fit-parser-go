package main

import (
	"fmt"
	"os"

	"./parsers"
	"./utils"
)

// TODO: ?Switch file reading solution to bufio
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fitParse fitFile")
		os.Exit(1)
	}

	fileName := os.Args[1]
	f, err := os.Open(fileName)
	defer f.Close()
	utils.Check(err)

	hsArr := utils.ByteReader(f, 1)
	headerSize := int(hsArr[0])
	header := parsers.ParseHeader(headerSize, f)

	parsers.ParseMessages(headerSize, header.DataSize, f)

	fmt.Println(header)
}
