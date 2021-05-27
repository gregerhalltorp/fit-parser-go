package fitTypes

import (
	"fmt"

	"../constants"
	"../constants/mesgNum"
	"../constants/types"
	"../utils"
)

func ReadMessage(defnMesg DefinitionMessage, data []byte) Message {
	mesg, typeFields := MessageCreator(defnMesg.GlobalMessageNumber)
	converter := utils.NewByteConverter(data, defnMesg.IsBigEndian)
	for _, fieldDefinition := range defnMesg.FieldDefinitions {
		// read := true
		var foundField Field

		found := false
		for _, field := range typeFields {
			if field.Num == fieldDefinition.Num {
				foundField = field
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Didn't find ", fieldDefinition.Num)
			foundField = GetUnknownField(fieldDefinition.Num, fieldDefinition.Type)
		}
		if foundField.FieldType != fieldDefinition.Type {
			fmt.Println("@@@@@@ Type mismatch!!!")
		}
		if foundField.FieldType&constants.BaseTypeNumMask == constants.String {
			value := converter.GetString(int(fieldDefinition.Size))
			foundField.AddValue(value)
		} else {
			baseType := GetBaseType(foundField.FieldType)
			numElements := fieldDefinition.Size / baseType.Size
			for i := 0; i < int(numElements); i++ {
				invalid := true
				var value interface{}

				switch foundField.FieldType & constants.BaseTypeNumMask {
				case constants.Enum, constants.Byte, constants.UInt8, constants.UInt8z:
					value = converter.GetByte()
				case constants.SInt8:
					value = converter.GetSByte()
				case constants.SInt16:
					value = converter.GetSInt16()
				case constants.UInt16, constants.UInt16z:
					value = converter.GetUInt16()
				case constants.SInt32:
					value = converter.GetSInt32()
				case constants.UInt32, constants.UInt32z:
					value = converter.GetUInt32()
				case constants.SInt64:
					value = converter.GetSInt64()
				case constants.UInt64, constants.UInt64z:
					value = converter.GetUInt64()
				case constants.Float32:
					value = converter.GetFloat32()
				case constants.Float64:
					value = converter.GetFloat64()
				default:
					value = converter.GetBytes(fieldDefinition.Size)
				}
				if value != baseType.InvalidValue {
					invalid = false
				}
				if !invalid {
					foundField.AddValue(value)
				}
			}
		}
		mesg.SetField(foundField)
	}
	return mesg
}

func MessageCreator(globalMesgNum uint16) (Message, []Field) {
	switch globalMesgNum {
	case mesgNum.FileId:
		return NewMessageFromNameNum("FileId", mesgNum.FileId), getFileIdMessageFields()
	case mesgNum.FileCreator:
		return NewMessageFromNameNum("FileCreator", mesgNum.FileCreator), getFileCreatorMessageFields()
	case mesgNum.Event:
		return NewMessageFromNameNum("Event", mesgNum.Event), getEventMessageFields()
	case mesgNum.DeviceInfo:
		return NewMessageFromNameNum("DeviceInfo", mesgNum.DeviceInfo), getDeviceInfoFields()
	case mesgNum.Sport:
		return NewMessageFromNameNum("Sport", mesgNum.Sport), getSportFields()
	case mesgNum.Record:
		return NewMessageFromNameNum("Record", mesgNum.Record), getRecordFields()
	// case mesgNum.SdmProfile:
	// 	mesg = NewMessageFromNameNum("SdmProfile", mesgNum.SdmProfile)
	// 	typeFields = getSdmProfileFields()
	default:
		fmt.Println("UNKNOWN MESSAGE, GLOBALMESSAGENUMBER", globalMesgNum)
		return NewMessageFromNameNum("unknown", mesgNum.Invalid), make([]Field, 0)
	}
}

func CreateMessage(messageNumber uint16) Message {
	mesg, typeFields := MessageCreator(messageNumber)
	mesg.Fields = typeFields
	return mesg
}

func getFileIdMessageFields() []Field {
	f := make([]Field, 7)
	f[0] = NewField("Type", 0, 0, 1, 0, "", false, types.File)
	f[1] = NewField("Manufacturer", 1, 132, 1, 0, "", false, types.Manufacturer)
	productField := NewField("Product", 2, 132, 1, 0, "", false, types.Uint16)

	faveroSubField := NewSubFieldFromValues("FaveroProduct", 132, 1, 0, "")
	productField.AddSubField(faveroSubField)
	faveroSubField.AddMap(1, 263)

	garminSubField := NewSubFieldFromValues("GarminProduct", 132, 1, 0, "")
	productField.AddSubField(garminSubField)
	garminSubField.AddMap(1, 1)
	garminSubField.AddMap(1, 15)
	garminSubField.AddMap(1, 13)
	garminSubField.AddMap(1, 89)
	f[2] = productField

	f[3] = NewField("SerialNumber", 3, 140, 1, 0, "", false, types.Uint32z)
	f[4] = NewField("TimeCreated", 4, 134, 1, 0, "", false, types.DateTime)
	f[5] = NewField("Number", 5, 132, 1, 0, "", false, types.Uint16)
	f[6] = NewField("ProductName", 8, 7, 1, 0, "", false, types.String)

	return f
}

func getFileCreatorMessageFields() []Field {
	f := make([]Field, 2)
	f[0] = NewField("SoftwareVersion", 0, 132, 1, 0, "", false, types.Uint16)
	f[1] = NewField("HardwareVersion", 1, 2, 1, 0, "", false, types.Uint8)

	return f
}

func getEventMessageFields() []Field {
	f := make([]Field, 15)
	f[0] = NewField("Timestamp", 253, 134, 1, 0, "s", false, types.DateTime)
	f[1] = NewField("Event", 0, 0, 1, 0, "", false, types.Event)
	f[2] = NewField("EventType", 1, 0, 1, 0, "", false, types.EventType)

	data16Field := NewField("Data16", 2, 132, 1, 0, "", false, types.Uint16)
	data16Field.AddComponent(NewFieldComponentFromValues(3, false, 16, 1, 0))
	f[3] = data16Field

	dataField := NewField("Data", 3, 134, 1, 0, "", false, types.Uint32)
	timerTriggerSubField := NewSubFieldFromValues("TimerTrigger", 0, 1, 0, "")
	timerTriggerSubField.AddMap(0, 0)
	dataField.AddSubField(timerTriggerSubField)
	coursePointIndexSubField := NewSubFieldFromValues("CoursePointIndex", 132, 1, 0, "")
	coursePointIndexSubField.AddMap(0, 10)
	dataField.AddSubField(coursePointIndexSubField)
	batteryLevelSubField := NewSubFieldFromValues("BatteryLevel", 132, 1000, 0, "V")
	batteryLevelSubField.AddMap(0, 11)
	dataField.AddSubField(batteryLevelSubField)
	virtualPartnerSpeedSubField := NewSubFieldFromValues("VirtualPartnerSpeed", 132, 1000, 0, "m/s")
	virtualPartnerSpeedSubField.AddMap(0, 12)
	dataField.AddSubField(virtualPartnerSpeedSubField)
	hrHighAlertSubField := NewSubFieldFromValues("HrHighAlert", 2, 1, 0, "bpm")
	hrHighAlertSubField.AddMap(0, 13)
	dataField.AddSubField(hrHighAlertSubField)
	hrLowAlertSubField := NewSubFieldFromValues("HrLowAlert", 2, 1, 0, "bpm")
	hrLowAlertSubField.AddMap(0, 14)
	dataField.AddSubField(hrLowAlertSubField)
	speedHighAlertSubField := NewSubFieldFromValues("SpeedHighAlert", 134, 1000, 0, "m/s")
	speedHighAlertSubField.AddMap(0, 15)
	dataField.AddSubField(speedHighAlertSubField)
	speedLowAlertSubField := NewSubFieldFromValues("SpeedLowAlert", 134, 1000, 0, "m/s")
	speedLowAlertSubField.AddMap(0, 16)
	dataField.AddSubField(speedLowAlertSubField)
	cadHighAlertSubField := NewSubFieldFromValues("CadHighAlert", 132, 1, 0, "rpm")
	cadHighAlertSubField.AddMap(0, 17)
	dataField.AddSubField(cadHighAlertSubField)
	cadLowAlertSubField := NewSubFieldFromValues("CadLowAlert", 132, 1, 0, "rpm")
	cadLowAlertSubField.AddMap(0, 18)
	dataField.AddSubField(cadLowAlertSubField)
	powerHighAlertSubField := NewSubFieldFromValues("PowerHighAlert", 132, 1, 0, "watts")
	powerHighAlertSubField.AddMap(0, 19)
	dataField.AddSubField(powerHighAlertSubField)
	powerLowAlertSubField := NewSubFieldFromValues("PowerLowAlert", 132, 1, 0, "watts")
	powerLowAlertSubField.AddMap(0, 20)
	dataField.AddSubField(powerLowAlertSubField)
	timeDurationAlertSubField := NewSubFieldFromValues("TimeDurationAlert", 134, 1000, 0, "s")
	timeDurationAlertSubField.AddMap(0, 23)
	dataField.AddSubField(timeDurationAlertSubField)
	distanceDurationAlertSubField := NewSubFieldFromValues("DistanceDurationAlert", 134, 100, 0, "m")
	distanceDurationAlertSubField.AddMap(0, 24)
	dataField.AddSubField(distanceDurationAlertSubField)
	calorieDurationAlertSubField := NewSubFieldFromValues("CalorieDurationAlert", 134, 1, 0, "calories")
	calorieDurationAlertSubField.AddMap(0, 25)
	dataField.AddSubField(calorieDurationAlertSubField)
	fitnessEquipmentStateSubField := NewSubFieldFromValues("FitnessEquipmentState", 0, 1, 0, "")
	fitnessEquipmentStateSubField.AddMap(0, 27)
	dataField.AddSubField(fitnessEquipmentStateSubField)
	sportPointSubField := NewSubFieldFromValues("SportPoint", 134, 1, 0, "")
	sportPointSubField.AddMap(0, 33)
	sportPointSubField.AddComponent(NewFieldComponentFromValues(7, false, 16, 1, 0)) // score
	sportPointSubField.AddComponent(NewFieldComponentFromValues(8, false, 16, 1, 0)) // opponent_score
	dataField.AddSubField(sportPointSubField)
	gearChangeDataSubField := NewSubFieldFromValues("GearChangeData", 134, 1, 0, "")
	gearChangeDataSubField.AddMap(0, 42)
	gearChangeDataSubField.AddMap(0, 43)
	gearChangeDataSubField.AddComponent(NewFieldComponentFromValues(11, false, 8, 1, 0)) // rear_gear_num
	gearChangeDataSubField.AddComponent(NewFieldComponentFromValues(12, false, 8, 1, 0)) // rear_gear
	gearChangeDataSubField.AddComponent(NewFieldComponentFromValues(9, false, 8, 1, 0))  // front_gear_num
	gearChangeDataSubField.AddComponent(NewFieldComponentFromValues(10, false, 8, 1, 0)) // front_gear
	dataField.AddSubField(gearChangeDataSubField)
	riderPositionSubField := NewSubFieldFromValues("RiderPosition", 0, 1, 0, "")
	riderPositionSubField.AddMap(0, 44)
	dataField.AddSubField(riderPositionSubField)
	commTimeoutSubField := NewSubFieldFromValues("CommTimeout", 132, 1, 0, "")
	commTimeoutSubField.AddMap(0, 47)
	dataField.AddSubField(commTimeoutSubField)
	radarThreatAlertSubField := NewSubFieldFromValues("RadarThreatAlert", 134, 1, 0, "")
	radarThreatAlertSubField.AddMap(0, 75)
	radarThreatAlertSubField.AddComponent(NewFieldComponentFromValues(21, false, 8, 1, 0)) // radar_threat_level_max
	radarThreatAlertSubField.AddComponent(NewFieldComponentFromValues(22, false, 8, 1, 0)) // radar_threat_count
	dataField.AddSubField(radarThreatAlertSubField)
	f[4] = dataField

	f[5] = NewField("EventGroup", 4, 2, 1, 0, "", false, types.Uint8)
	f[6] = NewField("Score", 7, 132, 1, 0, "", false, types.Uint16)
	f[7] = NewField("OpponentScore", 8, 132, 1, 0, "", false, types.Uint16)
	f[8] = NewField("FrontGearNum", 9, 10, 1, 0, "", false, types.Uint8z)
	f[9] = NewField("FrontGear", 10, 10, 1, 0, "", false, types.Uint8z)
	f[10] = NewField("RearGearNum", 11, 10, 1, 0, "", false, types.Uint8z)
	f[11] = NewField("RearGear", 12, 10, 1, 0, "", false, types.Uint8z)
	f[12] = NewField("DeviceIndex", 13, 2, 1, 0, "", false, types.DeviceIndex)
	f[13] = NewField("RadarThreatLevelMax", 21, 0, 1, 0, "", false, types.RadarThreatLevelType)
	f[14] = NewField("RadarThreatCount", 22, 2, 1, 0, "", false, types.Uint8)

	return f
}

func getDeviceInfoFields() []Field {
	f := make([]Field, 18)

	f[0] = NewField("Timestamp", 253, 134, 1, 0, "s", false, types.DateTime)
	f[1] = NewField("DeviceIndex", 0, 2, 1, 0, "", false, types.DeviceIndex)
	deviceTypeField := NewField("DeviceType", 1, 2, 1, 0, "", false, types.Uint8)
	antplusDeviceTypeSubField := NewSubFieldFromValues("AntplusDeviceType", 2, 1, 0, "")
	antplusDeviceTypeSubField.AddMap(25, 1)
	deviceTypeField.AddSubField(antplusDeviceTypeSubField)
	antDeviceTypeSubField := NewSubFieldFromValues("AntDeviceType", 2, 1, 0, "")
	antDeviceTypeSubField.AddMap(25, 0)
	f[2] = deviceTypeField
	f[3] = NewField("Manufacturer", 2, 132, 1, 0, "", false, types.Manufacturer)
	f[4] = NewField("SerialNumber", 3, 140, 1, 0, "", false, types.Uint32z)
	productField := NewField("Product", 4, 132, 1, 0, "", false, types.Uint16)
	faveroProductSubField := NewSubFieldFromValues("FaveroProduct", 132, 1, 0, "")
	faveroProductSubField.AddMap(2, 263)
	productField.AddSubField(faveroProductSubField)
	garminProductSubField := NewSubFieldFromValues("GarminProduct", 132, 1, 0, "")
	garminProductSubField.AddMap(2, 1)
	garminProductSubField.AddMap(2, 15)
	garminProductSubField.AddMap(2, 13)
	garminProductSubField.AddMap(2, 89)
	productField.AddSubField(garminProductSubField)
	f[5] = productField
	f[6] = NewField("SoftwareVersion", 5, 132, 100, 0, "", false, types.Uint16)
	f[7] = NewField("HardwareVersion", 6, 2, 1, 0, "", false, types.Uint8)
	f[8] = NewField("CumOperatingTime", 7, 134, 1, 0, "s", false, types.Uint32)
	f[9] = NewField("BatteryVoltage", 10, 132, 256, 0, "V", false, types.Uint16)
	f[10] = NewField("BatteryStatus", 11, 2, 1, 0, "", false, types.BatteryStatus)
	f[11] = NewField("SensorPosition", 18, 0, 1, 0, "", false, types.BodyLocation)
	f[12] = NewField("Descriptor", 19, 7, 1, 0, "", false, types.String)
	f[13] = NewField("AntTransmissionType", 20, 10, 1, 0, "", false, types.Uint8z)
	f[14] = NewField("AntDeviceNumber", 21, 139, 1, 0, "", false, types.Uint16z)
	f[15] = NewField("AntNetwork", 22, 0, 1, 0, "", false, types.AntNetwork)
	f[16] = NewField("SourceType", 25, 0, 1, 0, "", false, types.SourceType)
	f[17] = NewField("ProductName", 27, 7, 1, 0, "", false, types.String)

	return f
}

func getSdmProfileFields() []Field {
	f := make([]Field, 8)

	f[0] = NewField("MessageIndex", 254, 132, 1, 0, "", false, types.MessageIndex)
	f[1] = NewField("Enabled", 0, 0, 1, 0, "", false, types.Bool)
	f[2] = NewField("SdmAntId", 1, 139, 1, 0, "", false, types.Uint16z)
	f[3] = NewField("SdmCalFactor", 2, 132, 10, 0, "%", false, types.Uint16)
	f[4] = NewField("Odometer", 3, 134, 100, 0, "m", false, types.Uint32)
	f[5] = NewField("SpeedSource", 4, 0, 1, 0, "", false, types.Bool)
	f[6] = NewField("SdmAntIdTransType", 5, 10, 1, 0, "", false, types.Uint8z)
	f[7] = NewField("OdometerRollover", 7, 2, 1, 0, "", false, types.Uint8)

	return f
}

func getSportFields() []Field {
	f := make([]Field, 3)

	f[0] = NewField("Sport", 0, 0, 1, 0, "", false, types.Sport)
	f[1] = NewField("SubSport", 1, 0, 1, 0, "", false, types.SubSport)
	f[2] = NewField("Name", 3, 7, 1, 0, "", false, types.String)

	return f
}

func getRecordFields() []Field {
	f := make([]Field, 74)
	f[0] = NewField("Timestamp", 253, 134, 1, 0, "s", false, types.DateTime)
	f[1] = NewField("PositionLat", 0, 133, 1, 0, "semicircles", false, types.Sint32)
	f[2] = NewField("PositionLong", 1, 133, 1, 0, "semicircles", false, types.Sint32)
	altitudeField := NewField("Altitude", 2, 132, 5, 500, "m", false, types.Uint16)
	altitudeField.AddComponent(NewFieldComponentFromValues(78, false, 16, 5, 500)) // enhanced_altitude
	f[3] = altitudeField
	f[4] = NewField("HeartRate", 3, 2, 1, 0, "bpm", false, types.Uint8)
	f[5] = NewField("Cadence", 4, 2, 1, 0, "rpm", false, types.Uint8)
	f[6] = NewField("Distance", 5, 134, 100, 0, "m", true, types.Uint32)
	speedField := NewField("Speed", 6, 132, 1000, 0, "m/s", false, types.Uint16)
	speedField.AddComponent(NewFieldComponentFromValues(73, false, 16, 1000, 0)) // enhanced_speed
	f[7] = speedField
	f[8] = NewField("Power", 7, 132, 1, 0, "watts", false, types.Uint16)
	compressedSpeedDistanceField := NewField("CompressedSpeedDistance", 8, 13, 1, 0, "", false, types.Byte)
	compressedSpeedDistanceField.AddComponent(NewFieldComponentFromValues(6, false, 12, 100, 0)) // speed
	compressedSpeedDistanceField.AddComponent(NewFieldComponentFromValues(5, true, 12, 16, 0))   // distance
	f[9] = compressedSpeedDistanceField
	f[10] = NewField("Grade", 9, 131, 100, 0, "%", false, types.Sint16)
	f[11] = NewField("Resistance", 10, 2, 1, 0, "", false, types.Uint8)
	f[12] = NewField("TimeFromCourse", 11, 133, 1000, 0, "s", false, types.Sint32)
	f[13] = NewField("CycleLength", 12, 2, 100, 0, "m", false, types.Uint8)
	f[14] = NewField("Temperature", 13, 1, 1, 0, "C", false, types.Sint8)
	f[15] = NewField("Speed1s", 17, 2, 16, 0, "m/s", false, types.Uint8)
	cyclesField := NewField("Cycles", 18, 2, 1, 0, "cycles", false, types.Uint8)
	cyclesField.AddComponent(NewFieldComponentFromValues(19, true, 8, 1, 0)) // total_cycles
	f[16] = cyclesField
	f[17] = NewField("TotalCycles", 19, 134, 1, 0, "cycles", true, types.Uint32)
	compressedAccumulatedPowerField := NewField("CompressedAccumulatedPower", 28, 132, 1, 0, "watts", false, types.Uint16)
	compressedAccumulatedPowerField.AddComponent(NewFieldComponentFromValues(29, true, 16, 1, 0)) // accumulated_power
	f[18] = compressedAccumulatedPowerField
	f[19] = NewField("AccumulatedPower", 29, 134, 1, 0, "watts", true, types.Uint32)
	f[20] = NewField("LeftRightBalance", 30, 2, 1, 0, "", false, types.LeftRightBalance)
	f[21] = NewField("GpsAccuracy", 31, 2, 1, 0, "m", false, types.Uint8)
	f[22] = NewField("VerticalSpeed", 32, 131, 1000, 0, "m/s", false, types.Sint16)
	f[23] = NewField("Calories", 33, 132, 1, 0, "kcal", false, types.Uint16)
	f[24] = NewField("VerticalOscillation", 39, 132, 10, 0, "mm", false, types.Uint16)
	f[25] = NewField("StanceTimePercent", 40, 132, 100, 0, "percent", false, types.Uint16)
	f[26] = NewField("StanceTime", 41, 132, 10, 0, "ms", false, types.Uint16)
	f[27] = NewField("ActivityType", 42, 0, 1, 0, "", false, types.ActivityType)
	f[28] = NewField("LeftTorqueEffectiveness", 43, 2, 2, 0, "percent", false, types.Uint8)
	f[29] = NewField("RightTorqueEffectiveness", 44, 2, 2, 0, "percent", false, types.Uint8)
	f[30] = NewField("LeftPedalSmoothness", 45, 2, 2, 0, "percent", false, types.Uint8)
	f[31] = NewField("RightPedalSmoothness", 46, 2, 2, 0, "percent", false, types.Uint8)
	f[32] = NewField("CombinedPedalSmoothness", 47, 2, 2, 0, "percent", false, types.Uint8)
	f[33] = NewField("Time128", 48, 2, 128, 0, "s", false, types.Uint8)
	f[34] = NewField("StrokeType", 49, 0, 1, 0, "", false, types.StrokeType)
	f[35] = NewField("Zone", 50, 2, 1, 0, "", false, types.Uint8)
	f[36] = NewField("BallSpeed", 51, 132, 100, 0, "m/s", false, types.Uint16)
	f[37] = NewField("Cadence256", 52, 132, 256, 0, "rpm", false, types.Uint16)
	f[38] = NewField("FractionalCadence", 53, 2, 128, 0, "rpm", false, types.Uint8)
	f[39] = NewField("TotalHemoglobinConc", 54, 132, 100, 0, "g/dL", false, types.Uint16)
	f[40] = NewField("TotalHemoglobinConcMin", 55, 132, 100, 0, "g/dL", false, types.Uint16)
	f[41] = NewField("TotalHemoglobinConcMax", 56, 132, 100, 0, "g/dL", false, types.Uint16)
	f[42] = NewField("SaturatedHemoglobinPercent", 57, 132, 10, 0, "%", false, types.Uint16)
	f[43] = NewField("SaturatedHemoglobinPercentMin", 58, 132, 10, 0, "%", false, types.Uint16)
	f[44] = NewField("SaturatedHemoglobinPercentMax", 59, 132, 10, 0, "%", false, types.Uint16)
	f[45] = NewField("DeviceIndex", 62, 2, 1, 0, "", false, types.DeviceIndex)
	f[46] = NewField("LeftPco", 67, 1, 1, 0, "mm", false, types.Sint8)
	f[47] = NewField("RightPco", 68, 1, 1, 0, "mm", false, types.Sint8)
	f[48] = NewField("LeftPowerPhase", 69, 2, 0.7111111, 0, "degrees", false, types.Uint8)
	f[49] = NewField("LeftPowerPhasePeak", 70, 2, 0.7111111, 0, "degrees", false, types.Uint8)
	f[50] = NewField("RightPowerPhase", 71, 2, 0.7111111, 0, "degrees", false, types.Uint8)
	f[51] = NewField("RightPowerPhasePeak", 72, 2, 0.7111111, 0, "degrees", false, types.Uint8)
	f[52] = NewField("EnhancedSpeed", 73, 134, 1000, 0, "m/s", false, types.Uint32)
	f[53] = NewField("EnhancedAltitude", 78, 134, 5, 500, "m", false, types.Uint32)
	f[54] = NewField("BatterySoc", 81, 2, 2, 0, "percent", false, types.Uint8)
	f[55] = NewField("MotorPower", 82, 132, 1, 0, "watts", false, types.Uint16)
	f[56] = NewField("VerticalRatio", 83, 132, 100, 0, "percent", false, types.Uint16)
	f[57] = NewField("StanceTimeBalance", 84, 132, 100, 0, "percent", false, types.Uint16)
	f[58] = NewField("StepLength", 85, 132, 10, 0, "mm", false, types.Uint16)
	f[59] = NewField("AbsolutePressure", 91, 134, 1, 0, "Pa", false, types.Uint32)
	f[60] = NewField("Depth", 92, 134, 1000, 0, "m", false, types.Uint32)
	f[61] = NewField("NextStopDepth", 93, 134, 1000, 0, "m", false, types.Uint32)
	f[62] = NewField("NextStopTime", 94, 134, 1, 0, "s", false, types.Uint32)
	f[63] = NewField("TimeToSurface", 95, 134, 1, 0, "s", false, types.Uint32)
	f[64] = NewField("NdlTime", 96, 134, 1, 0, "s", false, types.Uint32)
	f[65] = NewField("CnsLoad", 97, 2, 1, 0, "percent", false, types.Uint8)
	f[66] = NewField("N2Load", 98, 132, 1, 0, "percent", false, types.Uint16)
	f[67] = NewField("Grit", 114, 136, 1, 0, "", false, types.Float32)
	f[68] = NewField("Flow", 115, 136, 1, 0, "", false, types.Float32)
	f[69] = NewField("EbikeTravelRange", 117, 132, 1, 0, "km", false, types.Uint16)
	f[70] = NewField("EbikeBatteryLevel", 118, 2, 1, 0, "percent", false, types.Uint8)
	f[71] = NewField("EbikeAssistMode", 119, 2, 1, 0, "depends on sensor", false, types.Uint8)
	f[72] = NewField("EbikeAssistLevelPercent", 120, 2, 1, 0, "percent", false, types.Uint8)
	f[73] = NewField("CoreTemperature", 139, 132, 100, 0, "C", false, types.Uint16)

	return f
}
