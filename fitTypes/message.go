package fitTypes

import (
	"errors"
	"fmt"

	"../constants"
	"../utils"
)

type Message struct {
	Name   string
	Num    uint16
	Fields []Field
}

func NewMessageFromNameNum(name string, num uint16) Message {
	m := Message{Name: name, Num: num, Fields: make([]Field, 0)}
	return m
}

func (m *Message) SetField(field Field) {
	if m.Fields == nil {
		m.Fields = make([]Field, 1)
		m.Fields[0] = field
	} else {
		m.Fields = append(m.Fields, field)
	}
}

func (m Message) GetFieldByNum(fieldNum byte) (Field, bool) {
	for _, f := range m.Fields {
		if f.Num == fieldNum {
			return f, true
		}
	}
	return Field{}, false
}

func (m Message) GetFieldByName(name string, checkSubFields bool) (Field, error) {
	for _, f := range m.Fields {
		if f.Name == name {
			return f, nil
		}
	}
	return Field{}, errors.New("Couldn't find field")
}

func (m Message) GetActiveSubFieldIndex(fieldNum byte) uint16 {
	theField, found := m.GetFieldByNum(fieldNum)
	if !found {
		return constants.SubfieldIndexMainField
	}
	for i, sf := range theField.SubFields {
		if sf.CanMesgSupport(m) {
			return uint16(i)
		}
	}
	return constants.SubfieldIndexMainField
}

func (m *Message) ExpandComponents(accumulator *utils.Accumulator) {
	for i := 0; i < len(m.Fields); i++ {
		var componentList []FieldComponent
		activeSubField := m.GetActiveSubFieldIndex(m.Fields[i].Num)
		if activeSubField == constants.SubfieldIndexMainField {
			componentList = m.Fields[i].Components
		} else {
			subField, found := m.Fields[i].GetSubfield(int(activeSubField))
			if !found {
				componentList = m.Fields[i].Components
			}
			componentList = subField.Components
		}
		fmt.Println(componentList)
		m.ExpandComponentsInList(componentList, m.Fields[i], 0, accumulator)
	}
}

func (m Message) ExpandComponentsInList(
	componentList []FieldComponent,
	currentField Field,
	offset int,
	accumulator *utils.Accumulator) []FieldComponentExpansion {
	expansions := make([]FieldComponentExpansion, 0)
	if len(componentList) > 0 {
		for _, fc := range componentList {
			if fc.FieldNum != constants.FieldNumInvalid {
				newField, ok := CreateMessage(m.Num).GetFieldByNum(fc.FieldNum)
				newField.IsExpandedField = true
				f, ok := m.GetFieldByNum(newField.Num)

				fmt.Println(newField, f, ok)
			}
		}
	}
	return expansions
}

type FieldComponentExpansion struct {
	offset int
	field  Field
}

func (fce FieldComponentExpansion) GetOffset() int {
	return fce.offset
}
func (fce FieldComponentExpansion) GetField() Field {
	return fce.field
}
