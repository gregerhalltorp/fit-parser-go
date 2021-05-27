package parsers

import (
	"encoding/binary"
	"fmt"
	"os"

	"../constants"
	"../fitTypes"
	"../utils"
)

func parseDefinitionMessage(headerByte byte, f *os.File) (int, fitTypes.DefinitionMessage) {
	headerArr := utils.ByteReader(f, 5)
	architecture := headerArr[1]
	globalMessageNumber := binary.LittleEndian.Uint16(headerArr[2:])
	fmt.Println("GlobalMessageNumber: ", globalMessageNumber)
	numberOfFields := headerArr[4]
	fieldDefs := make([]fitTypes.FieldDefinition, numberOfFields)
	isBigEndian := architecture == constants.BigEndian
	var j byte
	for j = 0; j < numberOfFields; j++ {
		fieldsArr := utils.ByteReader(f, 3)
		num := fieldsArr[0]
		size := fieldsArr[1]
		fieldType := fieldsArr[2]
		fieldDefinition := fitTypes.NewFieldDefinition(num, size, fieldType)
		fieldDefs[j] = fieldDefinition
	}
	return 5 + 3*int(numberOfFields), fitTypes.NewDefinitionMessage(
		headerByte&constants.LocalMesgNumMask,
		globalMessageNumber,
		architecture,
		headerByte,
		isBigEndian,
		fieldDefs,
	)
}

func parseMessage(headerByte byte, definitionMessage fitTypes.DefinitionMessage, f *os.File) (int, fitTypes.Message) {
	var message fitTypes.Message
	fieldsSize := definitionMessage.GetMessageSize() - 1 // The header is already read
	data := utils.ByteReader(f, fieldsSize)
	message = fitTypes.ReadMessage(definitionMessage, data)
	return 1 + fieldsSize, message
}

func ParseMessages(headerSize int, fileSize uint32, f *os.File) {
	definitionMessages := make(map[byte]fitTypes.DefinitionMessage)
	messages := make([]fitTypes.Message, 0)
	var timestamp uint32 = 0
	var lastTimeOffset int32 = 0
	accumulator := utils.NewAccumulator()
	position := headerSize
	MAX_MSG := 26
	msgRead := 0
	for uint32(position) < fileSize-2 && msgRead < MAX_MSG {
		headerByte := utils.ByteReader(f, 1)
		localMessageNumber := headerByte[0] & constants.LocalMesgNumMask
		if headerByte[0]&constants.MesgDefinitionMask == constants.MesgDefinitionMask {
			fmt.Println("Position: ", position, "Got a definition message, parsing", headerByte[0]&constants.LocalMesgNumMask)
			pos, dm := parseDefinitionMessage(headerByte[0], f)
			fmt.Println("Definition message length:", pos)
			position += pos
			definitionMessages[dm.LocalMessageNumber] = dm
			// c, err := json.Marshal(definitionMessages)
			// if err != nil {
			// 	fmt.Println("ERROR", err)
			// }
			msgRead++
			fmt.Println("Position: ", position, "Msg read:", msgRead) //, "Definition messages: ", string(c))
		} else if headerByte[0]&constants.CompressedHeaderMask == constants.CompressedHeaderMask {
			fmt.Println("Got a compressed header message")
			os.Exit(0)
		} else if headerByte[0]&constants.MesgDefinitionMask == constants.MesgHeaderMask {
			definitionMessage, present := definitionMessages[localMessageNumber]
			if !present {
				mesg := fmt.Sprintf("Missing definition message for message number %d", localMessageNumber)
				panic(mesg)
			}
			fmt.Println("Position: ", position, "Got a normal header message", localMessageNumber)
			pos, mesg := parseMessage(headerByte[0], definitionMessage, f)
			fmt.Println("Message length:", pos)
			timeStampField, tsErr := mesg.GetFieldByName("Timestamp", true)
			if tsErr == nil {
				tsValue := timeStampField.GetValue(0)
				if tsValue != nil {
					timestamp = tsValue.(uint32)
					lastTimeOffset = int32(timestamp) & int32(constants.CompressedTimeMask)
				}
				fmt.Println("GOT A TIMESTAMP VALUE: ", timestamp, lastTimeOffset)
			}
			for _, f := range mesg.Fields {
				if f.IsAccumulated {
					for i := 0; i < len(f.Values); i++ {
						value, err := f.GetInt64Value(i)
						if err != nil {
							fmt.Println(err)
							continue
						}
						for _, fi := range mesg.Fields {
							for _, fc := range fi.Components {
								if fc.FieldNum == f.Num && fc.Accumulate {
									value = int64((float64(value)/float64(f.Scale) - float64(f.Offset) + float64(fc.Offset)) * float64(fc.Scale))
								}
							}
						}
						accumulator.Set(int(mesg.Num), int(f.Num), value)
					}
				}
			}
			mesg.ExpandComponents(&accumulator)

			position += pos
			messages = append(messages, mesg)
			// b, err := json.Marshal(messages)
			// if err != nil {
			// 	fmt.Println("ERROR", err)
			// }
			msgRead++
			fmt.Println("Position:", position, "Msg read:", msgRead) //, "Messages: ", string(b))
		} else {
			fmt.Println("BAD VALUE:", headerByte)
			panic("BAILING!")
		}
	}
}
