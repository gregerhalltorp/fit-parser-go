package fitTypes

type Header struct {
	Size            int
	ProtocolVersion byte
	ProfileVersion  uint16
	DataSize        uint32
	DataType        string
	Crc             uint16
}

func NewHeader(
	size int,
	protocolVersion byte,
	profileVersion uint16,
	dataSize uint32,
	dataType string,
	crc uint16) *Header {
	h := Header{
		Size:            size,
		ProtocolVersion: protocolVersion,
		ProfileVersion:  profileVersion,
		DataSize:        dataSize,
		DataType:        dataType,
		Crc:             crc,
	}
	return &h
}
