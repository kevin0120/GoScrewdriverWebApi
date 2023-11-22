package udp

type Sdo struct {
	sdoKey   string
	dataType int
	sdoValue interface{}
}

func SdoGenerator(sdoKey string, dataType int, sdoValue interface{}) *Sdo {
	return &Sdo{
		sdoKey:   sdoKey,
		dataType: dataType,
		sdoValue: sdoValue,
	}
}
