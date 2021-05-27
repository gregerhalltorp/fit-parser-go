package utils

import (
	"encoding/binary"
	"math"
)

type ByteConverter struct {
	data        []byte
	position    int
	isBigEndian bool
	reader      binary.ByteOrder
}

func NewByteConverter(data []byte, isBigEndian bool) ByteConverter {
	var reader binary.ByteOrder
	if isBigEndian {
		reader = binary.BigEndian
	} else {
		reader = binary.LittleEndian
	}

	return ByteConverter{
		data:        data,
		position:    0,
		isBigEndian: isBigEndian,
		reader:      reader,
	}
}

func (b *ByteConverter) GetString(size int) string {
	// fmt.Println("GetString at position", b.position, "reading", size, "bytes")
	bytes := b.data[b.position : b.position+size]
	b.position += size
	nonZero := false
	actualBytes := make([]byte, 0)
	for _, byt := range bytes {
		if byt != 0 {
			nonZero = true
			break
		}
	}
	if nonZero {
		for _, byt := range bytes {
			actualBytes = append(actualBytes, byt)
			if byt == 0 {
				return string(actualBytes)
			}
		}
	}
	return ""
}

func (b *ByteConverter) GetByte() byte {
	// fmt.Println("GetByte at position", b.position, "reading 1 byte")
	value := b.data[b.position]
	b.position++
	return value
}

func (b *ByteConverter) GetSByte() int8 {
	// fmt.Println("GetSByte at position", b.position, "reading 1 byte")
	value := int8(b.data[b.position])
	b.position++
	return value
}

func (b *ByteConverter) GetSInt16() int16 {
	// fmt.Println("GetSInt16 at position", b.position, "reading 2 bytes")
	uValue := b.reader.Uint16(b.data[b.position : b.position+2])
	b.position += 2
	return int16(uValue)
}

func (b *ByteConverter) GetUInt16() uint16 {
	// fmt.Println("GetUInt16 at position", b.position, "reading 2 bytes")
	value := b.reader.Uint16(b.data[b.position : b.position+2])
	b.position += 2
	return value
}

func (b *ByteConverter) GetSInt32() int32 {
	// fmt.Println("GetSInt32 at position", b.position, "reading 4 bytes")
	value := b.reader.Uint32(b.data[b.position : b.position+4])
	b.position += 4
	return int32(value)
}

func (b *ByteConverter) GetUInt32() uint32 {
	// fmt.Println("GetUInt32 at position", b.position, "reading 4 bytes")
	value := b.reader.Uint32(b.data[b.position : b.position+4])
	b.position += 4
	return value
}

func (b *ByteConverter) GetSInt64() int64 {
	// fmt.Println("GetSInt64 at position", b.position, "reading 8 bytes")
	value := b.reader.Uint64(b.data[b.position : b.position+8])
	b.position += 8
	return int64(value)
}

func (b *ByteConverter) GetUInt64() uint64 {
	// fmt.Println("GetUInt64 at position", b.position, "reading 8 bytes")
	value := b.reader.Uint64(b.data[b.position : b.position+8])
	b.position += 8
	return value
}

func (b *ByteConverter) GetFloat32() float32 {
	// fmt.Println("GetFloat32 at position", b.position, "reading 4 bytes")
	bits := b.reader.Uint32(b.data[b.position : b.position+4])
	b.position += 4
	return math.Float32frombits(bits)
}

func (b *ByteConverter) GetFloat64() float64 {
	// fmt.Println("GetFloat64 at position", b.position, "reading 8 bytes")
	bits := b.reader.Uint64(b.data[b.position : b.position+8])
	b.position += 8
	return math.Float64frombits(bits)
}

func (b *ByteConverter) GetBytes(size byte) []byte {
	// fmt.Println("GetBytes at position", b.position, "reading", size, "bytes")
	value := b.data[b.position : b.position+int(size)]
	b.position += int(size)
	return value
}
