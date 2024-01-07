package udps

type Sdo struct {
	SdoKey   string
	DataType int
	SdoValue []byte
}

func (s *Sdo) Bytes() []byte {
	data := append(append(append(GetRawData(s.SdoKey, String), GetRawData(s.DataType, U8)...), GetRawData(len(s.SdoValue), U32)...), s.SdoValue...)
	return data
}

func SdoGenerator(sdoKey string, dataType int, sdoValue []byte) *Sdo {
	return &Sdo{
		SdoKey:   sdoKey,
		DataType: dataType,
		SdoValue: sdoValue,
	}
}

func bytesPacker(data []byte, opType uint16, rid uint32) []byte {
	var r []byte
	r = append(append(append(append(append(r, GetRawData(len(data)+6, U32)...)), GetRawData(rid, U32)...), GetRawData(opType, U16)...), data...)
	return r
}

func SdoPack(data []byte, opType uint16, rid uint32) []byte {
	var r []byte
	switch opType {
	case MSdo:
		r = bytesPacker(data, MSdo, rid)
		break
	case QSdo:
		r = bytesPacker(data, QSdo, rid)
		break
	case QSdoDefines:
		r = bytesPacker([]byte{0x00}, QSdoDefines, rid)
		break
	case SUB_ITEM:
		r = bytesPacker(data, SUB_ITEM, rid)
		break
	case SUB_TOPIC:
		r = bytesPacker([]byte("\x00"), QSdoDefines, rid)
		break
	case ReadJson:
		buff := GetRawData(len(data), U32)
		r = bytesPacker(append(buff, data...), QSdoDefines, rid)
		break
	}
	return r
}

//func SdoUnPack(data []byte) (error uint16 interface{}) {
//	var r interface{}
//	switch opType {
//	case MSdo:
//		r = false
//		if len(data) == 1 {
//			r = data[0] == 0x01
//		}
//		break
//	case QSdo:
//		r = 0
//		if len(data) == 1 {
//			r = int8(data[0])
//		}
//		break
//	case QSdoDefines:
//		r = 0
//		if len(data) == 1 {
//			r = data[0]
//		}
//		break
//	case SUB_ITEM:
//		r = 0
//		if len(data) == 2 {
//			r = int16(binary.LittleEndian.Uint16(data))
//		}
//		break
//	case SUB_TOPIC:
//		r = 0
//		if len(data) == 2 {
//			r = binary.LittleEndian.Uint16(data)
//		}
//		break
//	case ReadJson:
//		r = 0
//		if len(data) == 4 {
//			r = int32(binary.LittleEndian.Uint32(data))
//		}
//		break
//	}
//	return r
//}
