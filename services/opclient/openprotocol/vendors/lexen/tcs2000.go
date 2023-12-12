package lexen

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol"
	"time"
)

type TCS2000Controller struct {
	Controller
}

func (c *TCS2000Controller) New() openprotocol.IOpenProtocolController {
	controller := TCS2000Controller{}
	controller.SetInstance(&controller)
	return &controller
}

func (c *TCS2000Controller) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							*每个MID对应的REV版本
		openprotocol.MID_0001_START:                 "001",
		openprotocol.MID_0018_PSET:                  "001",
		openprotocol.MID_0014_PSET_SUBSCRIBE:        "001",
		openprotocol.MID_0060_LAST_RESULT_SUBSCRIBE: "001",
		openprotocol.MID_0062_LAST_RESULT_ACK:       "001",
		openprotocol.MID_0064_OLD_SUBSCRIBE:         "006",
		openprotocol.MID_0012_PSET_DETAIL_REQUEST:   "002",
		openprotocol.MID_0010_PSET_LIST_REQUEST:     "009",
		openprotocol.MID_0042_TOOL_DISABLE:          "001",
		openprotocol.MID_0043_TOOL_ENABLE:           "001",
		openprotocol.MID_0019_PSET_BATCH_SET:        "001",
		//openprotocol.MID_0070_ALARM_SUBSCRIBE:       "001",
		openprotocol.MID_0040_TOOL_INFO_REQUEST: "005",

		//openprotocol.MID_7408_LAST_CURVE_SUBSCRIBE: "001",
		openprotocol.MID_7411_LAST_CURVE_DATA_ACK: "001",
	}

	return vendorModels
}

func (c *TCS2000Controller) OpenProtocolParams() *openprotocol.OpenProtocolParams {
	return &openprotocol.OpenProtocolParams{
		MaxKeepAliveCheck: 3,
		MaxReplyTime:      5 * time.Second,
		KeepAlivePeriod:   time.Duration(c.DeviceConf.KeepAlive),
	}
}
