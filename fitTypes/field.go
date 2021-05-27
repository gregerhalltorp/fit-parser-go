package fitTypes

import (
	"errors"

	"../constants"
	"../constants/types"
)

type Field struct {
	Name            string
	Num             byte
	FieldType       byte
	Scale           float32
	Offset          float32
	Units           string
	IsAccumulated   bool
	IsExpandedField bool
	ProfileType     types.Type
	Values          []interface{}
	SubFields       []SubField
	Components      []FieldComponent
}

func NewField(
	name string,
	num byte,
	fType byte,
	scale float32,
	offset float32,
	units string,
	isAccumulated bool,
	profileType types.Type) Field {
	f := Field{
		Name:          name,
		Num:           num,
		FieldType:     fType,
		Scale:         scale,
		Offset:        offset,
		Units:         units,
		IsAccumulated: isAccumulated,
		ProfileType:   profileType,
		Values:        make([]interface{}, 0),
		SubFields:     make([]SubField, 0),
		Components:    make([]FieldComponent, 0),
	}

	return f
}

func GetUnknownField(num byte, fType byte) Field {
	return Field{
		Name:            "unknown",
		Num:             num,
		FieldType:       fType,
		Scale:           1,
		Offset:          0,
		Units:           "",
		IsAccumulated:   false,
		ProfileType:     types.Enum,
		IsExpandedField: false,
	}
}

func (f *Field) AddSubField(s SubField) {
	f.SubFields = append(f.SubFields, s)
}

func (f *Field) AddValue(v interface{}) {
	f.Values = append(f.Values, v)
}

func (f *Field) AddComponent(c FieldComponent) {
	f.Components = append(f.Components, c)
}

func (f Field) GetInt64Value(index int) (int64, error) {
	if index >= len(f.Values) || index < 0 {
		return int64(0), errors.New("index out of range")
	}

	value := f.Values[index]
	shiftedType := f.FieldType & constants.BaseTypeNumMask
	switch shiftedType {
	case constants.Enum, constants.Byte, constants.UInt8, constants.UInt8z:
		return int64(value.(byte)), nil
	case constants.SInt8:
		return int64(value.(int8)), nil
	case constants.UInt16, constants.UInt16z:
		return int64(value.(uint16)), nil
	case constants.SInt16:
		return int64(value.(int16)), nil
	case constants.UInt32, constants.UInt32z:
		return int64(value.(uint32)), nil
	case constants.SInt32:
		return int64(value.(int32)), nil
	case constants.UInt64, constants.UInt64z:
		return int64(value.(uint64)), nil
	case constants.SInt64:
		return value.(int64), nil
	case constants.Float32:
		return int64(value.(float32)), nil
	case constants.Float64:
		return int64(value.(float64)), nil
	case constants.String:
		return 0, errors.New("Cannot get string as int64")
	}

	return 0, errors.New("Type not found")
}

func (f Field) GetValue(index int) interface{} {
	if index >= len(f.Values) || index < 0 {
		return nil
	}
	value := f.Values[index]

	shiftedType := f.FieldType & constants.BaseTypeNumMask
	switch shiftedType {
	case constants.Enum:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(byte))
		}
		return value.(byte)
	case constants.Byte, constants.UInt8, constants.UInt8z:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(byte))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(byte))/f.Scale - f.Offset
		}
		return value.(byte)
	case constants.SInt8:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(int8))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(int8))/f.Scale - f.Offset
		}
		return value.(int8)
	case constants.UInt16, constants.UInt16z:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(uint16))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(uint16))/f.Scale - f.Offset
		}
		return value.(uint16)
	case constants.SInt16:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(int16))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(int16))/f.Scale - f.Offset
		}
		return value.(int16)
	case constants.UInt32, constants.UInt32z:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(uint32))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(uint32))/f.Scale - f.Offset
		}
		return value.(uint32)
	case constants.SInt32:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(int32))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(int32))/f.Scale - f.Offset
		}
		return value.(int32)
	case constants.UInt64, constants.UInt64z:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(uint64))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(uint64))/f.Scale - f.Offset
		}
		return value.(uint64)
	case constants.SInt64:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(int64))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(int64))/f.Scale - f.Offset
		}
		return value.(int64)
	case constants.Float32:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(float32))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(float32))/f.Scale - f.Offset
		}
		return value.(float32)
	case constants.Float64:
		if value == GetBaseType(shiftedType).InvalidValue && f.Scale != 1.0 {
			return float32(value.(float64))
		}
		if f.Scale != 1.0 || f.Offset != 0.0 {
			return float32(value.(float64))/f.Scale - f.Offset
		}
		return value.(float64)
	case constants.String:
		return value.(string)
	default:
		return nil
	}
}

func (f Field) GetBitsValue(offset int, bits int, componentType byte) (int64, bool) {
	bitsInValue := 0
	index := 0

	if int(GetBaseType(componentType&constants.BaseTypeNumMask).Size*8) < bits {
		bits = int(GetBaseType(componentType&constants.BaseTypeNumMask).Size * 8)
	}

	if len(f.Values) == 0 {
		return -1, false
	}

	for bitsInValue < bits {
		if index == len(f.Values) {
			return -1, false
		}
	}
}

func (f Field) GetSubfield(subFieldIndex int) (SubField, bool) {
	if subFieldIndex >= 0 && subFieldIndex < len(f.SubFields) {
		return f.SubFields[subFieldIndex], true
	}

	return SubField{}, false
}
