package fitTypes

type SubFieldMap struct {
	RefFieldNum   byte
	RefFieldValue interface{}
}

func (sfm SubFieldMap) GetRefFieldValueInt64() int64 {
	bv, ok := sfm.RefFieldValue.(byte)
	if ok {
		return int64(bv)
	}
	i8v, ok := sfm.RefFieldValue.(int8)
	if ok {
		return int64(i8v)
	}
	u16v, ok := sfm.RefFieldValue.(uint16)
	if ok {
		return int64(u16v)
	}
	i16v, ok := sfm.RefFieldValue.(int16)
	if ok {
		return int64(i16v)
	}
	u32v, ok := sfm.RefFieldValue.(uint32)
	if ok {
		return int64(u32v)
	}
	i32v, ok := sfm.RefFieldValue.(int32)
	if ok {
		return int64(i32v)
	}
	u64v, ok := sfm.RefFieldValue.(uint64)
	if ok {
		return int64(u64v)
	}
	i64v, ok := sfm.RefFieldValue.(int64)
	if ok {
		return int64(i64v)
	}
	f32v, ok := sfm.RefFieldValue.(float32)
	if ok {
		return int64(f32v)
	}
	f64v, ok := sfm.RefFieldValue.(float64)
	if ok {
		return int64(f64v)
	}
	return int64(0)
}

func (sfm SubFieldMap) CanMesgSupport(mesg Message) bool {
	field, found := mesg.GetFieldByNum(sfm.RefFieldNum)
	if found {
		fVal, err := field.GetInt64Value(0)
		if err == nil && fVal == sfm.GetRefFieldValueInt64() {
			return true
		}
	}
	return false
}

type SubField struct {
	Name       string
	Type       byte
	Scale      float32
	Offset     float32
	Units      string
	Maps       []SubFieldMap
	Components []FieldComponent
}

func NewSubFieldFromValues(
	name string,
	fType byte,
	scale float32,
	offset float32,
	units string,
) SubField {
	sf := SubField{
		Name:       name,
		Type:       fType,
		Scale:      scale,
		Offset:     offset,
		Units:      units,
		Maps:       make([]SubFieldMap, 0),
		Components: make([]FieldComponent, 0),
	}

	return sf
}

func (sf *SubField) AddMap(key byte, value interface{}) {
	temp := SubFieldMap{RefFieldNum: key, RefFieldValue: value}
	sf.Maps = append(sf.Maps, temp)
}

func (sf *SubField) AddComponent(c FieldComponent) {
	sf.Components = append(sf.Components, c)
}

func (sf SubField) CanMesgSupport(mesg Message) bool {
	for _, sfm := range sf.Maps {
		if sfm.CanMesgSupport(mesg) {
			return true
		}
	}

	return false
}
