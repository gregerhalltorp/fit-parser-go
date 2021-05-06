package parsers

import (
	"encoding/binary"
	"fmt"
	"os"

	"../constants"
	"../fitTypes"
	"../utils"
)

// "./fitTypes"

func getCrc(headerSize int, hArr []byte) uint16 {
	if headerSize == constants.HeaderWithCRCSize {
		return binary.LittleEndian.Uint16(hArr[11:])
	}
	return 0
}

func ParseHeader(headerSize int, f *os.File) *fitTypes.Header {
	hArr := utils.ByteReader(f, headerSize-1)
	protVer := hArr[0]

	if protVer&constants.ProtocolVersionMajorMask > constants.ProtocolVersionCheck {
		mesg := fmt.Sprintf("FIT decode error: Protocol Version %d.X not supported by SDK Protocol Ver%d.%d",
			(protVer&constants.ProtocolVersionMajorMask)>>constants.ProtocolVersionMajorShift,
			constants.ProtocolMajorVersion,
			constants.ProtocolMinorVersion)
		panic(mesg)
	}

	profileVersion := binary.LittleEndian.Uint16(hArr[1:])
	dataSize := binary.LittleEndian.Uint32(hArr[3:])
	dataType := string(hArr[7:11])

	header := fitTypes.NewHeader(headerSize, protVer, profileVersion, dataSize, dataType, getCrc(headerSize, hArr))

	return header
}
