package fitTypes

type FieldComponent struct {
	FieldNum         byte
	Accumulate       bool
	Bits             int
	Scale            float32
	Offset           float32
	AccumulatedValue int64
	LastValue        int64
}

func NewFieldComponentFromValues(
	fieldNum byte,
	accumulate bool,
	bits int,
	scale float32,
	offset float32) FieldComponent {
	return FieldComponent{
		FieldNum:         fieldNum,
		Accumulate:       accumulate,
		Bits:             bits,
		Scale:            scale,
		Offset:           offset,
		AccumulatedValue: 0,
		LastValue:        0,
	}
}
