package opserver

import "fmt"

func GeneratePackage(mid, rev string, noack bool, station, spindle, data string) string {
	noAckStr := ""
	if noack {
		noAckStr = "1"
	}
	h := OpenProtocolHeader{
		MID:      mid,
		LEN:      LenHeader + len(data),
		Revision: rev,
		NoAck:    noAckStr,
		Station:  station,
		Spindle:  spindle,
		Spare:    "",
	}

	return h.Serialize() + data + string(rune(OpTerminal))
}

func (c *Connect) Err004Response(mid string, code string) error {
	data := fmt.Sprintf("%-4s%-2s", mid, code)
	res := GeneratePackage("0004", DefaultRev, true, "", "", data)
	_, err := c.tcpConnect.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) Success005Response(mid string) error {
	res := GeneratePackage("0005", DefaultRev, true, "", "", mid)
	_, err := c.tcpConnect.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) MidResponse(mid string, data string) error {
	res := GeneratePackage(mid, DefaultRev, true, "", "", data)
	_, err := c.tcpConnect.Write([]byte(res))
	if err != nil {
		return err
	}
	return nil
}
