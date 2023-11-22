package udp

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"math"
)

const (
	Group         = 0
	Bool          = 1
	I8            = 2
	U8            = 3
	I16           = 4
	U16           = 5
	I32           = 6
	U32           = 7
	I64           = 8
	U64           = 9
	F32           = 10
	F64           = 11
	String        = 12
	StringUnicode = 13
	ByteArray     = 14
	Remap         = 15
)

const (
	MSdo        = 0x0001
	QSdo        = 0x0002
	QSdoDefines = 0x0003
	SUB_ITEM    = 0x0004
	SUB_TOPIC   = 0x0005
	ReadJson    = 0x0030
	//UnPack
	UnpackSdoRet      = 0x8001
	UnpackQSdoData    = 0x8002
	UnpackSdoDict     = 0x8003
	UnpackSub         = 0x8004
	UnpackTopicsQuery = 0x8005
	UnpackReadJson    = 0x8030
)

func GetRawData(val interface{}, dataType int) []byte {
	var data []byte
	switch dataType {
	case Bool:
		b, _ := val.(bool)
		if b {
			data = append(data, uint8(1))
		} else {
			data = append(data, uint8(0))
		}
		break
	case I8:
		b, _ := val.(int8)
		data = append(data, byte(b))
		break
	case U8:
		b, _ := val.(uint8)
		data = append(data, b)
		break
	case I16:
		b, _ := val.(int16)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		break
	case U16:
		b, _ := val.(uint16)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		break
	case I32:
		b, _ := val.(int32)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		data = append(data, byte(b>>16))
		data = append(data, byte(b>>24))
		break
	case U32:
		b, _ := val.(uint32)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		data = append(data, byte(b>>16))
		data = append(data, byte(b>>24))
		break
	case I64:
		b, _ := val.(int64)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		data = append(data, byte(b>>16))
		data = append(data, byte(b>>24))
		data = append(data, byte(b>>32))
		data = append(data, byte(b>>40))
		data = append(data, byte(b>>48))
		data = append(data, byte(b>>56))
		break
	case U64:
		b, _ := val.(uint64)
		data = append(data, byte(b))
		data = append(data, byte(b>>8))
		data = append(data, byte(b>>16))
		data = append(data, byte(b>>24))
		data = append(data, byte(b>>32))
		data = append(data, byte(b>>40))
		data = append(data, byte(b>>48))
		data = append(data, byte(b>>56))
		break

	case F32:
		b, _ := val.(float32)
		bits := math.Float32bits(b)
		data = append(data, byte(bits))
		data = append(data, byte(bits>>8))
		data = append(data, byte(bits>>16))
		data = append(data, byte(bits>>24))
		break
	case F64:
		b, _ := val.(float64)
		bits := math.Float64bits(b)
		data = append(data, byte(bits))
		data = append(data, byte(bits>>8))
		data = append(data, byte(bits>>16))
		data = append(data, byte(bits>>24))
		data = append(data, byte(bits>>32))
		data = append(data, byte(bits>>40))
		data = append(data, byte(bits>>48))
		data = append(data, byte(bits>>56))
		break
	case String:
		b, _ := val.(string)
		c := b + "\x00"
		data = append(data, []byte(c)...)
		break

	case StringUnicode:
		b, _ := val.(string)
		c := b + "\x00"
		encoder := simplifiedchinese.GBK.NewEncoder()
		gbkBytes, _ := encoder.Bytes([]byte(c))
		data = append(data, gbkBytes...)
		break

	case ByteArray:
		b, _ := val.([]byte)
		data = append(data, b...)
		bytesLen := uint32(len(b))
		data = append(data, byte(bytesLen))
		data = append(data, byte(bytesLen>>8))
		data = append(data, byte(bytesLen>>16))
		data = append(data, byte(bytesLen>>24))
		break

	case Remap:
		b, _ := val.([]byte)
		data = append(data, b...)
		break

	}
	return data
}

func sdoPack() {

}
