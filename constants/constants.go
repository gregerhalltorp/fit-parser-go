package constants

const ProtocolVersionMajorShift byte = 4
const ProtocolVersionMajorMask byte = 0x0F << ProtocolVersionMajorShift
const ProtocolMajorVersion byte = 2
const ProtocolMinorVersion byte = 0
const ProtocolVersionCheck byte = ProtocolMajorVersion << ProtocolVersionMajorShift
const HeaderWithCRCSize = 14
const HeaderWithoutCRCSize = 12
const MesgDefinitionMask byte = 0x40
