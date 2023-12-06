package desoutter

import (
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/tightening_device"
	"time"
)

type CVI2Controller struct {
	Controller
}

func (c *CVI2Controller) New() openprotocol.IOpenProtocolController {
	controller := CVI2Controller{}
	controller.SetInstance(&controller)
	return &controller
}

func (c *CVI2Controller) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							*每个MID对应的REV版本
		openprotocol.MID_0001_START:                        "001",
		openprotocol.MID_0018_PSET:                         "001",
		openprotocol.MID_0014_PSET_SUBSCRIBE:               "001",
		openprotocol.MID_0060_LAST_RESULT_SUBSCRIBE:        "001",
		openprotocol.MID_0062_LAST_RESULT_ACK:              "001",
		openprotocol.MID_0012_PSET_DETAIL_REQUEST:          "001",
		openprotocol.MID_0010_PSET_LIST_REQUEST:            "001",
		openprotocol.MID_0042_TOOL_DISABLE:                 "001",
		openprotocol.MID_0043_TOOL_ENABLE:                  "001",
		openprotocol.MID_0019_PSET_BATCH_SET:               "001",
		openprotocol.MID_0070_ALARM_SUBSCRIBE:              "001",
		openprotocol.MID_0040_TOOL_INFO_REQUEST:            "001",
		openprotocol.MID_0210_INPUT_SUBSCRIBE:              "001",
		openprotocol.MID_0051_VIN_SUBSCRIBE:                "001",
		openprotocol.MID_7402_Cycle_Phase_Result_Subscribe: "001",
		openprotocol.MID_7408_LAST_CURVE_SUBSCRIBE:         "001",
		openprotocol.MID_7411_LAST_CURVE_DATA_ACK:          "001",
		openprotocol.IoModel: io.IoConfig{
			InputNum:  0,
			OutputNum: 0,
		},
	}

	return vendorModels
}

// 可重写所有TighteningController中的方法
func (c *CVI2Controller) CreateIO() tightening_device.ITighteningIO {
	return nil
}

func (c *CVI2Controller) OpenProtocolParams() *openprotocol.OpenProtocolParams {
	return &openprotocol.OpenProtocolParams{
		MaxKeepAliveCheck: 3,
		MaxReplyTime:      30 * time.Second,
		KeepAlivePeriod:   3 * time.Second,
	}
}

func (c *CVI2Controller) InitSubscribeInfos() {
	c.ControllerSubscribes = []openprotocol.ControllerSubscribe{
		c.ResultSubscribe,
		//c.SelectorSubscribe,
		//c.JobInfoSubscribe,
		//c.IOInputSubscribe,
		//c.VinSubscribe,
		//c.AlarmSubscribe,
		c.CurveSubscribe,
		c.CycleAndPhaseResultSubscribe,
	}
}
