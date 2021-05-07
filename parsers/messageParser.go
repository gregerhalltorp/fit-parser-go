package parsers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"../constants"
	"../constants/mesgNum"
	"../fitTypes"
	"../profile"
	"../utils"
)

func parseDefinitionMessage(headerByte byte, f *os.File) (int, fitTypes.DefinitionMessage) {
	headerArr := utils.ByteReader(f, 5)
	architecture := headerArr[1]
	globalMessageNumber := binary.LittleEndian.Uint16(headerArr[2:])
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

func parseMessage(headerByte byte, definitionMessage fitTypes.DefinitionMessage, f *os.File) int {
	fieldsSize := definitionMessage.GetMessageSize() - 1 // The header is already read
	fmt.Println("#####", fieldsSize)
	utils.ByteReader(f, fieldsSize)
	if definitionMessage.GlobalMessageNumber == mesgNum.FileId {
		message := profile.CreateMessage(definitionMessage.GlobalMessageNumber)
		fmt.Println("fileId!", message)
		// Set the right endian reader depending on definitionMessage.IsBigEndian
	}
	return 1 + fieldsSize
}

func ParseMessages(headerSize int, fileSize uint32, f *os.File) {
	definitionMessages := make(map[byte]fitTypes.DefinitionMessage)
	position := headerSize
	MAX_MSG := 4
	msgRead := 0
	for uint32(position) < fileSize-2 && msgRead < MAX_MSG {
		headerByte := utils.ByteReader(f, 1)
		localMessageNumber := headerByte[0] & constants.LocalMesgNumMask
		if headerByte[0]&constants.MesgDefinitionMask == constants.MesgDefinitionMask {
			fmt.Println("Got a definition message, parsing", headerByte[0]&constants.LocalMesgNumMask)
			pos, dm := parseDefinitionMessage(headerByte[0], f)
			position = pos
			definitionMessages[dm.LocalMessageNumber] = dm
			fmt.Println(pos, definitionMessages)
			msgRead++
		} else if headerByte[0]&constants.CompressedHeaderMask == constants.CompressedHeaderMask {
			fmt.Println("Got a compressed header message")
			os.Exit(0)
		} else if headerByte[0]&constants.MesgDefinitionMask == constants.MesgHeaderMask {
			definitionMessage, present := definitionMessages[localMessageNumber]
			if !present {
				mesg := fmt.Sprintf("Missing definition message for message number %d", localMessageNumber)
				panic(mesg)
			}
			fmt.Println("Got a normal header message", localMessageNumber)
			pos := parseMessage(headerByte[0], definitionMessage, f)
			position = pos
			b, err := json.Marshal(definitionMessage)
			if err != nil {
				fmt.Println("ERROR", err)
			}
			fmt.Println(string(b))
			msgRead++
		}
	}
}
