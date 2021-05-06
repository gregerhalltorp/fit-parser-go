package fitTypes

type FieldDefinition struct {
	Num  byte
	Size byte
	Type byte
}

func NewFieldDefinition(num byte, size byte, fType byte) FieldDefinition {
	fd := FieldDefinition{Num: num, Size: size, Type: fType}
	return fd
}
