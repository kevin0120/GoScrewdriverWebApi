package opserver

type OpenProtocolPayload struct {
	OpenProtocolHeader
	data string
}

func (c *Connect) Deserialize(request string) *OpenProtocolPayload {
	if len(request) < 20 {
		c.Err004Response("0000", "01")
		return nil
	}
	p := OpenProtocolPayload{}
	p.Deserialize(request[:20])

	if request[len(request)-1] != OpTerminal {
		c.Err004Response(p.MID, "01")
		return nil
	}
	if len(request) > 21 {
		p.data = request[20 : len(request)-1]
	}
	return &p
}
func (p *OpenProtocolPayload) Handle(c *Connect) {
	if p.MID == MID_0001_START {
		c.HandleMid0001Request(p)
		return
	}
	if !c.start {
		c.Err004Response(p.MID, "79")
		return
	}
	switch p.MID {
	case MID_0010_PSET_LIST_REQUEST:
		c.HandleMid0010Request(p)
		break
	case MID_0014_PSET_SUBSCRIBE:
		c.HandleMid0014Request(p)
		break
	case MID_0018_PSET:
		c.HandleMid0018Request(p)
		break
	case MID_0042_TOOL_DISABLE:
		c.HandleMid0042Request(p)
		break
	case MID_0043_TOOL_ENABLE:
		c.HandleMid0043Request(p)
		break
	case MID_0060_LAST_RESULT_SUBSCRIBE:
		c.HandleMid0060Request(p)
		break
	case MID_0150_IDENTIFIER_SET:
		c.HandleMid0150Request(p)
		break
	default:
		c.Err004Response(p.MID, "99")

	}
}

func (c *Connect) HandleMid0001Request(request *OpenProtocolPayload) error {
	c.start = true
	c.Success005Response(MID_0001_START)
	return nil
}
func (c *Connect) HandleMid0010Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0010_PSET_LIST_REQUEST)
	return nil
}
func (c *Connect) HandleMid0014Request(request *OpenProtocolPayload) error {

	c.Success005Response(MID_0014_PSET_SUBSCRIBE)
	return nil
}
func (c *Connect) HandleMid0018Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0018_PSET)
	return nil
}
func (c *Connect) HandleMid0042Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0042_TOOL_DISABLE)
	return nil
}
func (c *Connect) HandleMid0043Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0043_TOOL_ENABLE)
	return nil
}
func (c *Connect) HandleMid0060Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0060_LAST_RESULT_SUBSCRIBE)
	return nil
}
func (c *Connect) HandleMid0150Request(request *OpenProtocolPayload) error {
	c.Success005Response(MID_0150_IDENTIFIER_SET)
	return nil
}
