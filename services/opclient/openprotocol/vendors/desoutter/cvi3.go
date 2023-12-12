package desoutter

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol"
)

type CVI3Controller struct {
	Controller
}

func (c *CVI3Controller) New() openprotocol.IOpenProtocolController {
	controller := CVI3Controller{}
	controller.SetInstance(&controller)
	return &controller
}

func (c *CVI3Controller) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							  *每个MID对应的REV版本
		openprotocol.MID_0001_START:                   "004",
		openprotocol.MID_0018_PSET:                    "001",
		openprotocol.MID_0014_PSET_SUBSCRIBE:          "001",
		openprotocol.MID_0034_JOB_INFO_SUBSCRIBE:      "004",
		openprotocol.MID_0250_SELECTOR_SUBSCRIBE:      "001",
		openprotocol.MID_0060_LAST_RESULT_SUBSCRIBE:   "998",
		openprotocol.MID_0062_LAST_RESULT_ACK:         "998",
		openprotocol.MID_0150_IDENTIFIER_SET:          "001",
		openprotocol.MID_0038_JOB_SELECT:              "002",
		openprotocol.MID_0064_OLD_SUBSCRIBE:           "006",
		openprotocol.MID_0130_JOB_OFF:                 "001",
		openprotocol.MID_0012_PSET_DETAIL_REQUEST:     "002",
		openprotocol.MID_0010_PSET_LIST_REQUEST:       "001",
		openprotocol.MID_0032_JOB_DETAIL_REQUEST:      "003",
		openprotocol.MID_0030_JOB_LIST_REQUEST:        "002",
		openprotocol.MID_0042_TOOL_DISABLE:            "001",
		openprotocol.MID_0043_TOOL_ENABLE:             "001",
		openprotocol.MID_0200_CONTROLLER_RELAYS:       "001",
		openprotocol.MID_0019_PSET_BATCH_SET:          "001",
		openprotocol.MID_0210_INPUT_SUBSCRIBE:         "001",
		openprotocol.MID_0127_JOB_ABORT:               "001",
		openprotocol.MID_0100_MULTI_SPINDLE_SUBSCRIBE: "001",
		openprotocol.MID_0051_VIN_SUBSCRIBE:           "002",
		openprotocol.MID_0070_ALARM_SUBSCRIBE:         "001",
		openprotocol.MID_0040_TOOL_INFO_REQUEST:       "002",

		openprotocol.MID_7408_LAST_CURVE_SUBSCRIBE: "001",
		openprotocol.MID_7411_LAST_CURVE_DATA_ACK:  "001",
	}

	return vendorModels
}

// 可重写所有TighteningController中的方法
