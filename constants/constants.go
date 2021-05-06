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
