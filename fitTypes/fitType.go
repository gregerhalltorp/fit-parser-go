package fitTypes

import (
	"encoding/binary"
	"math"
)

type FitType struct {
	EndianAbility bool
	BaseTypeField byte
	TypeName      string
	InvalidValue  interface{}
	Size          byte
	IsSigned      bool
	IsInteger     bool
}

func NewFitType(
	endianAbility bool,
	baseTypeField byte,
	typeName string,
	invalidValue interface{},
	size byte,
	isSigned bool,
	isInt bool) FitType {
	return FitType{
		EndianAbility: endianAbility,
		BaseTypeField: baseTypeField,
		TypeName:      typeName,
		InvalidValue:  invalidValue,
		Size:          size,
		IsSigned:      isSigned,
		IsInteger:     isInt,
	}
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func GetBaseType(t byte) FitType {
	switch t {
	case 0x00:
		return NewFitType(false, 0x00, "enum", byte(0xFF), 1, false, false)
	case 0x01:
		return NewFitType(false, 0x01, "sint8", int8(0x7F), 1, true, true)
	case 0x02:
		return NewFitType(false, 0x02, "uint8", byte(0xFF), 1, false, true)
	case 0x83:
		return NewFitType(true, 0x83, "sint16", int16(0x7FFF), 2, true, true)
	case 0x84:
		return NewFitType(true, 0x84, "uint16", uint16(0xFFFF), 2, false, true)
	case 0x85:
		return NewFitType(true, 0x85, "sint32", int32(0x7FFFFFFF), 4, true, true)
	case 0x86:
		return NewFitType(true, 0x86, "uint32", uint32(0xFFFFFFFF), 4, false, true)
	case 0x07:
		return NewFitType(false, 0x07, "string", byte(0x00), 1, false, false)
	case 0x88:
		return NewFitType(true, 0x88, "float32", Float32FromBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF}), 4, true, false)
	case 0x89:
		return NewFitType(true, 0x89, "float64", Float64FromBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}), 8, true, false)
	case 0x0A:
		return NewFitType(false, 0x0A, "uint8z", byte(0x00), 1, false, true)
	case 0x8B:
		return NewFitType(true, 0x8B, "uint16z", uint16(0x0000), 2, false, true)
	case 0x8C:
		return NewFitType(true, 0x8C, "uint32z", uint32(0x00000000), 4, false, true)
	case 0x0D:
		return NewFitType(false, 0x0D, "byte", byte(0xFF), 1, false, false)
	case 0x8E:
		return NewFitType(true, 0x8E, "sint64", int64(0x7FFFFFFFFFFFFFFF), 8, true, true)
	case 0x8F:
		return NewFitType(true, 0x8F, "uint64", uint64(0xFFFFFFFFFFFFFFFF), 8, false, true)
	case 0x90:
		return NewFitType(true, 0x90, "uint64z", uint64(0x0000000000000000), 8, false, true)
	default:
		return NewFitType(false, 0x00, "bad value", 0xFF, 1, false, false)
	}
}
