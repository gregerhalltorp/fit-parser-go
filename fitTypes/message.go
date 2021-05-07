package fitTypes

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
