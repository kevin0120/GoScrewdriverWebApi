package openprotocol

import (
	"github.com/kevin0120/GoScrewdriverWebApi/utils"
)

var needOldTighteningResult = utils.GetEnvBool("ENV_NEED_OLD_RESULT", false) //是否需要处理历史拧紧结果

type MidHandler func(controller *TighteningController, pkg *handlerPkg) error

/*
标准op协议 & desoutter扩展
*/

func HandleMid9999Alive(c *TighteningController, pkg *handlerPkg) error {
	return nil
}

func HandleMid0002StartAck(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 处理曲线
func HandleMid7410LastCurve(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

func HandleMid7404PhaseResult(c *TighteningController, pkg *handlerPkg) (err error) {

	return nil
}

func HandleMid7406CycleResult(c *TighteningController, pkg *handlerPkg) (err error) {

	return nil

}

// 处理结果
func HandleMid0061LastResult(c *TighteningController, pkg *handlerPkg) (err error) {

	return nil
}

// 处理历史结果
func HandleMid0065OldData(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// pset详情
func HandleMid0013PsetDetailReply(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// pset列表
func HandleMid0011PsetListReply(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// job列表
func HandleMid0031JobListReply(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// job详情
func HandleMid0033JobDetailReply(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 请求错误
func HandleMid0004CmdErr(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 请求成功
func HandleMid0005CmdOk(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// job推送信息
func HandleMid0035JobInfo(c *TighteningController, pkg *handlerPkg) error {

	return nil

}

// 控制器输入变化
func HandleMid0211InputMonitor(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 多轴结果
func HandleMid0101MultiSpindleResult(c *TighteningController, pkg *handlerPkg) error {
	ms := MultiSpindleResult{}
	ms.Deserialize(pkg.Body)

	return nil
}

const DummyBarCode = 0

// 收到条码推送
func HandleMid0052Vin(c *TighteningController, pkg *handlerPkg) (err error) {

	return
}

// 报警信息
func HandleMid0071Alarm(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 报警状态
func HandleMid0076AlarmStatus(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

// 工具状态(维护)
func HandleMid0041ToolInfoReply(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

/*
lexen 扩展
*/

// 处理曲线
func HandleMid7410LastCurveNoAck(c *TighteningController, pkg *handlerPkg) error {

	return nil
}

func HandleMid0011PsetListInfo(c *TighteningController, pkg *handlerPkg) error {

	return nil
}
