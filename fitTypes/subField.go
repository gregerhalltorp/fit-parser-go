package fitTypes

type SubField struct {
	Name   string
	Type   byte
	Scale  float32
	Offset float32
	Units  string
	Maps   []SubFieldMap
}

type SubFieldMap struct {
	refFieldNum   byte
	refFieldValue interface{}
}

func NewSubFieldFromValues(
	name string,
	fType byte,
	scale float32,
	offset float32,
	units string,
) SubField {
	sf := SubField{
		Name:   name,
		Type:   fType,
		Scale:  scale,
		Offset: offset,
		Units:  units,
		Maps:   make([]SubFieldMap, 0),
	}

	return sf
}

func (sf *SubField) AddMap(key byte, value interface{}) {
	sf.Maps = append(sf.Maps, SubFieldMap{refFieldNum: key, refFieldValue: value})
}
