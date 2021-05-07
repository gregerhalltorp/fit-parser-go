package fitTypes

import (
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
	SubFields       []SubField
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
		SubFields:     make([]SubField, 0),
	}

	return f
}

func (f *Field) AddSubField(s SubField) {
	f.SubFields = append(f.SubFields, s)
}
