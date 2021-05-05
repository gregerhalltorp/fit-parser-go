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

	header := fitTypes.NewHeader(headerSize, protVer, profileVersion, dataSize, dataType)

	if headerSize == constants.HeaderWithCRCSize {
		header.Crc = binary.LittleEndian.Uint16(hArr[11:])
	} else {
		header.Crc = 0
	}

	return header
}
