package crane

import (
	"errors"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/tightening_device"
	"time"
)

type IQWrenchController struct {
	Controller
}

func (c *IQWrenchController) New() openprotocol.IOpenProtocolController {
	controller := IQWrenchController{}
	controller.SetInstance(&controller)
	return &controller
}

func (c *IQWrenchController) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							*每个MID对应的REV版本
		openprotocol.MID_0001_START:                 "004",
		openprotocol.MID_0018_PSET:                  "001",
		openprotocol.MID_0014_PSET_SUBSCRIBE:        "001",
		openprotocol.MID_0060_LAST_RESULT_SUBSCRIBE: "006",
		openprotocol.MID_0062_LAST_RESULT_ACK:       "006",
		openprotocol.MID_0064_OLD_SUBSCRIBE:         "006",
		openprotocol.MID_0012_PSET_DETAIL_REQUEST:   "002",
		openprotocol.MID_0010_PSET_LIST_REQUEST:     "001",
		openprotocol.MID_0042_TOOL_DISABLE:          "001",
		openprotocol.MID_0043_TOOL_ENABLE:           "001",
		openprotocol.MID_0019_PSET_BATCH_SET:        "001",
		openprotocol.MID_0070_ALARM_SUBSCRIBE:       "001",
		openprotocol.MID_0040_TOOL_INFO_REQUEST:     "005",

		openprotocol.IoModel: io.IoConfig{
			InputNum:  0,
			OutputNum: 0,
		},
	}

	return vendorModels
}

// 可重写所有TighteningController中的方法
func (c *IQWrenchController) GetToolViaChannel(channel int) (tightening_device.ITighteningTool, error) {
	for _, v := range c.Children() {
		return v.(tightening_device.ITighteningTool), nil
	}

	return nil, errors.New("IQWrenchController.GetToolViaChannel: Tool Not Found")
}

func (c *IQWrenchController) HandleStatus(sn string, status string) {
	c.TighteningController.HandleStatus(sn, status)

	for _, tool := range c.Children() {
		tool.(*openprotocol.TighteningTool).UpdateStatus(status)
	}
}

func (c *IQWrenchController) CreateIO() tightening_device.ITighteningIO {
	return nil
}

func (c *IQWrenchController) OpenProtocolParams() *openprotocol.OpenProtocolParams {
	return &openprotocol.OpenProtocolParams{
		MaxKeepAliveCheck: 3,
		MaxReplyTime:      3 * time.Second,
		KeepAlivePeriod:   5 * time.Second,
	}
}
