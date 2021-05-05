package utils

import (
	"fmt"
	"os"
)

func ByteReader(file *os.File, len int) []byte {
	bArr := make([]byte, len)
	readBytes, err := file.Read(bArr)
	Check(err)
	if readBytes != len {
		mesg := fmt.Sprintf("Couldn't read bytes, attempted %d, got %d", len, readBytes)
		panic(mesg)
	}

	return bArr
}
