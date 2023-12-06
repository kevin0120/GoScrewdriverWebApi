package desoutter

import (
	"fmt"
	"github.com/masami10/rush/services/openprotocol"
)

type Controller struct {
	openprotocol.TighteningController
}

var midHandlers = map[string]openprotocol.MidHandler{
	openprotocol.MID_9999_ALIVE:                openprotocol.HandleMid9999Alive,
	openprotocol.MID_0002_START_ACK:            openprotocol.HandleMid0002StartAck,
	openprotocol.MID_0005_CMD_OK:               openprotocol.HandleMid0005CmdOk,
	openprotocol.MID_0004_CMD_ERR:              openprotocol.HandleMid0004CmdErr,
	openprotocol.MID_7404_PHASE_RESULT:         openprotocol.HandleMid7404PhaseResult,
	openprotocol.MID_7406_CYCLE_RESULT:         openprotocol.HandleMid7406CycleResult,
	openprotocol.MID_7410_LAST_CURVE:           openprotocol.HandleMid7410LastCurve,
	openprotocol.MID_0061_LAST_RESULT:          openprotocol.HandleMid0061LastResult,
	openprotocol.MID_0065_OLD_DATA:             openprotocol.HandleMid0065OldData,
	openprotocol.MID_0013_PSET_DETAIL_REPLY:    openprotocol.HandleMid0013PsetDetailReply,
	openprotocol.MID_0011_PSET_LIST_REPLY:      openprotocol.HandleMid0011PsetListReply,
	openprotocol.MID_0031_JOB_LIST_REPLY:       openprotocol.HandleMid0031JobListReply,
	openprotocol.MID_0033_JOB_DETAIL_REPLY:     openprotocol.HandleMid0033JobDetailReply,
	openprotocol.MID_0035_JOB_INFO:             openprotocol.HandleMid0035JobInfo,
	openprotocol.MID_0211_INPUT_MONITOR:        openprotocol.HandleMid0211InputMonitor,
	openprotocol.MID_0101_MULTI_SPINDLE_RESULT: openprotocol.HandleMid0101MultiSpindleResult,
	openprotocol.MID_0052_VIN:                  openprotocol.HandleMid0052Vin,
	openprotocol.MID_0071_ALARM:                openprotocol.HandleMid0071Alarm,
	openprotocol.MID_0076_ALARM_STATUS:         openprotocol.HandleMid0076AlarmStatus,
	openprotocol.MID_0041_TOOL_INFO_REPLY:      openprotocol.HandleMid0041ToolInfoReply,
}

func (c *Controller) GetMidHandler(mid string) (openprotocol.MidHandler, error) {
	h, exist := midHandlers[mid]
	if !exist {
		return nil, fmt.Errorf("Handler Not Found, Mid: %s", mid)
	}

	return h, nil
}
