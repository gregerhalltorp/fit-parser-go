package fitTypes

type DefinitionMessage struct {
	LocalMessageNumber  byte
	GlobalMessageNumber uint16
	Architecture        byte
	FieldDefinitions    []FieldDefinition
	Header              byte
	IsBigEndian         bool
}

func NewDefinitionMessage(
	localMessageNumber byte,
	globalMessageNumber uint16,
	architecture byte,
	header byte,
	isBigEndian bool,
	fieldDefinitions []FieldDefinition,
) DefinitionMessage {
	dm := DefinitionMessage{
		LocalMessageNumber:  localMessageNumber,
		GlobalMessageNumber: globalMessageNumber,
		Architecture:        architecture,
		Header:              header,
		FieldDefinitions:    fieldDefinitions,
		IsBigEndian:         isBigEndian,
	}

	return dm
}

func (dm DefinitionMessage) GetMessageSize() int {
	tot := 1 // for the header
	for _, field := range dm.FieldDefinitions {
		tot += int(field.Size)
	}
	return tot
}
