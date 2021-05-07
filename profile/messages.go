package profile

import (
	"../constants/mesgNum"
	"../constants/types"
	"../fitTypes"
)

func CreateMessage(globalMessageNumber uint16) fitTypes.Message {
	switch globalMessageNumber {
	case mesgNum.FileId:
		return createFileIdMessage()
	}

	panic("asflaksjflasdk")
}

func createFileIdMessage() fitTypes.Message {
	m := fitTypes.NewMessageFromNameNum("FileId", mesgNum.FileId)
	m.SetField(fitTypes.NewField("Type", 0, 0, 1, 0, "", false, types.File))
	m.SetField(fitTypes.NewField("Manufacturer", 1, 132, 1, 0, "", false, types.Manufacturer))
	productField := fitTypes.NewField("Product", 2, 132, 1, 0, "", false, types.Uint16)

	faveroSubField := fitTypes.NewSubFieldFromValues("FaveroProduct", 132, 1, 0, "")
	productField.AddSubField(faveroSubField)
	faveroSubField.AddMap(1, 263)

	garminSubField := fitTypes.NewSubFieldFromValues("GarminProduct", 132, 1, 0, "")
	productField.AddSubField(garminSubField)
	garminSubField.AddMap(1, 1)
	garminSubField.AddMap(1, 15)
	garminSubField.AddMap(1, 13)
	garminSubField.AddMap(1, 89)
	m.SetField(productField)

	m.SetField(fitTypes.NewField("SerialNumber", 3, 140, 1, 0, "", false, types.Uint32z))
	m.SetField(fitTypes.NewField("TimeCreated", 4, 134, 1, 0, "", false, types.DateTime))
	m.SetField(fitTypes.NewField("Number", 5, 132, 1, 0, "", false, types.Uint16))
	m.SetField(fitTypes.NewField("ProductName", 8, 7, 1, 0, "", false, types.String))

	return m
}
