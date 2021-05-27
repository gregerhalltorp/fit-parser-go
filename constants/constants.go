package constants

const ProtocolVersionMajorShift byte = 4
const ProtocolVersionMajorMask byte = 0x0F << ProtocolVersionMajorShift
const ProtocolMajorVersion byte = 2
const ProtocolMinorVersion byte = 0
const ProtocolVersionCheck byte = ProtocolMajorVersion << ProtocolVersionMajorShift
const HeaderWithCRCSize = 14
const HeaderWithoutCRCSize = 12
const MesgDefinitionMask byte = 0x40
const CompressedHeaderMask byte = 0x80
const MesgHeaderMask byte = 0x00
const BigEndian byte = 0x01
const LocalMesgNumMask byte = 0x0F
const FieldNumInvalid byte = 255
const CompressedTimeMask byte = 0x1F
const SubfieldNameMainField string = ""
const SubfieldIndexActiveSubfield uint16 = 0xFFFE
const SubfieldIndexMainField uint16 = SubfieldIndexActiveSubfield + 1

const Enum byte = 0x00
const SInt8 byte = 0x01
const UInt8 byte = 0x02
const SInt16 byte = 0x03
const UInt16 byte = 0x04
const SInt32 byte = 0x05
const UInt32 byte = 0x06
const String byte = 0x07
const Float32 byte = 0x08
const Float64 byte = 0x09
const UInt8z byte = 0x0A
const UInt16z byte = 0x0B
const UInt32z byte = 0x0C
const Byte byte = 0x0D
const SInt64 byte = 0x0E
const UInt64 byte = 0x0F
const UInt64z byte = 0x10

const BaseTypeNumMask byte = 0x1F
