package utils

type AccumulatedField struct {
	MessageNumber          int
	DestinationFieldNumber int
	lastValue              int64
	accumulatedValue       int64
}

func NewAccumulatedField(messageNumber int, fieldNumber int) AccumulatedField {
	return AccumulatedField{
		MessageNumber:          messageNumber,
		DestinationFieldNumber: fieldNumber,
		lastValue:              0,
		accumulatedValue:       0,
	}
}

func (af *AccumulatedField) Accumulate(value int64, bits int) int64 {
	mask := (int64(1) << int64(bits)) - 1
	af.accumulatedValue += (value - af.lastValue) & mask
	return af.accumulatedValue
}

func (af *AccumulatedField) Set(value int64) int64 {
	af.accumulatedValue = value
	af.lastValue = value
	return af.accumulatedValue
}

type Accumulator struct {
	accumulatedFields []AccumulatedField
}

func NewAccumulator() Accumulator {
	return Accumulator{accumulatedFields: make([]AccumulatedField, 0)}
}

func (a *Accumulator) Set(mesgNum int, destFieldNum int, value int64) {
	found := false

	for _, af := range a.accumulatedFields {
		if af.MessageNumber == mesgNum && af.DestinationFieldNumber == destFieldNum {
			af.Set(value)
			found = true
			break
		}
	}
	if !found {
		accField := NewAccumulatedField(mesgNum, destFieldNum)
		accField.Set(value)
		a.accumulatedFields = append(a.accumulatedFields, accField)
	}
}

func (a *Accumulator) Accumulate(mesgNum int, destFieldNum int, value int64, bits int) int64 {
	for _, af := range a.accumulatedFields {
		if af.MessageNumber == mesgNum && af.DestinationFieldNumber == destFieldNum {
			return af.Accumulate(value, bits)
		}
	}

	accField := NewAccumulatedField(mesgNum, destFieldNum)
	returnValue := accField.Accumulate(value, bits)
	a.accumulatedFields = append(a.accumulatedFields, accField)
	return returnValue
}
