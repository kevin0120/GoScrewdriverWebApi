package openprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/utils"
	"github.com/reactivex/rxgo/v2"

	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils/ascii"
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
	client := c.getClient(pkg.SN)
	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: requestErrors["00"],
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	go c.processSubscribeControllerInfo(pkg.SN)
	if needOldTighteningResult {
		c.diag.Info("Do Request Old Tightening Result!!!")
		//c.solveOldResult(pkg.SN)
	}

	return nil
}

// 处理曲线
func HandleMid7410LastCurve(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	//讲收到的曲线先做反序列化处理
	var curve CurveBody
	var tool tightening_device.ITighteningTool
	err := ascii.Unmarshal(pkg.Body, &curve)
	if err != nil {
		c.diag.Error("ascii.Unmarshal", err)
	}
	if curve.ToolChannelNumber == 0 {
		e := errors.New("收到的结果曲线数据不合法，未指定工具号")
		c.diag.Error("handleMID_7410_LAST_CURVE", e)
		return e
	}

	torqueCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TorqueString), 64)
	angleCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.AngleString), 64)
	timeCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TimeString), 64)
	Torque, Angle := c.CurveDataDecoding([]byte(curve.Data), torqueCoefficient, angleCoefficient, c.diag)

	client.tempResultCurve.CUR_M = append(client.tempResultCurve.CUR_M, Torque...)
	client.tempResultCurve.CUR_W = append(client.tempResultCurve.CUR_W, Angle...)
	//当本次数据为本次拧紧曲线的最后一次数据时
	if curve.Num == curve.Id {
		//若取到的点的数量大于协议解析出来该曲线的点数，多出的部分删掉，否则有多少发多少.
		if curve.MeasurePoints < len(client.tempResultCurve.CUR_M) {
			client.tempResultCurve.CUR_M = client.tempResultCurve.CUR_M[0:curve.MeasurePoints]
			client.tempResultCurve.CUR_W = client.tempResultCurve.CUR_W[0:curve.MeasurePoints]
		}

		client.tempResultCurve.GenerateTimeCurveByCoef(float32(timeCoefficient * 1000))

		//本次曲线全部解析完毕后,降临时存储的数据清空
		tool, err = c.getInstance().GetToolViaChannel(curve.ToolChannelNumber)
		if err != nil {
			return err
		}
		//todo: 找到唯一的工具
		sn := tool.SerialNumber()

		//defer delete(client.tempResultCurve, curve.ToolChannelNumber)
		client.tempResultCurve.TighteningUnit = sn
		client.tempResultCurve.UpdateTime = time.Now()
		c.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherCurve), client.tempResultCurve)
		// dispatch完后创建新的缓存拧紧曲线
		client.tempResultCurve = tightening_device.NewTighteningCurve()
	}
	err = client.SendOpenProtocolAckMsg(MID_7411_LAST_CURVE_DATA_ACK, "", "", "")
	if err != nil {
		return err
	}
	return nil
}

func HandleMid7404PhaseResult(c *TighteningController, pkg *handlerPkg) (err error) {
	var phaseResultData PhaseResultData
	client := c.getClient(pkg.SN)
	err = ascii.Unmarshal(pkg.Body, &phaseResultData)
	if err != nil {
		return
	}

	phaseResultData.Torque = phaseResultData.Torque / 100
	phaseResultData.TorqueMin = phaseResultData.TorqueMin / 100
	phaseResultData.TorqueMax = phaseResultData.TorqueMax / 100
	phaseResultData.TorqueFinalTarget = phaseResultData.TorqueFinalTarget / 100
	phaseResultData.Angle = phaseResultData.Angle / 10
	phaseResultData.AngleMin = phaseResultData.AngleMin / 10
	phaseResultData.AngleMax = phaseResultData.AngleMax / 10
	phaseResultData.FinalAngleTarget = phaseResultData.FinalAngleTarget / 10

	client.tempPhaseResult.Phase = append(client.tempPhaseResult.Phase, phaseResultData)
	client.tempPhaseResult.PhaseTorque = append(client.tempPhaseResult.PhaseTorque, phaseResultData.Torque)
	client.tempPhaseResult.PhaseAngle = append(client.tempPhaseResult.PhaseAngle, phaseResultData.Angle)

	return nil
}

func HandleMid7406CycleResult(c *TighteningController, pkg *handlerPkg) (err error) {
	client := c.getClient(pkg.SN)
	defer func() {
		client.tempPhaseResult = &PhaseResult{}
	}()
	var cycleResultData CycleResultData
	err = ascii.Unmarshal(pkg.Body, &cycleResultData)
	if err != nil {
		return
	}
	tighteningResult := cycleResultData.ToTighteningResult(client.tempPhaseResult)
	err = c.handleResult(&tighteningResult)
	return err

}

// 处理结果
func HandleMid0061LastResult(c *TighteningController, pkg *handlerPkg) (err error) {
	var tool tightening_device.ITighteningTool
	resultData := NewResultData(pkg.Header.Revision)

	client := c.getClient(pkg.SN)
	err = ascii.Unmarshal(pkg.Body, resultData.GetInstance())
	if err != nil && !strings.Contains(err.Error(), "message is not enough") {
		return
	}
	if err != nil {
		return
	}
	tighteningResult := resultData.ToTighteningResult()
	// rev1 没有工具序列号所以从配置文件读取
	if tighteningResult.ToolSN == "" {
		tool, err = c.GetToolViaChannel(tighteningResult.ChannelID)
		if err == nil {
			tighteningResult.ToolSN = tool.SerialNumber()
		}
	}
	if err = c.handleResult(tighteningResult); err == nil {
		// 发送ack报文
		err = client.SendOpenProtocolAckMsg(MID_0062_LAST_RESULT_ACK, "", "", "")
		if err != nil {
			return
		}
	} else {

	}
	return
}

// 处理历史结果
func HandleMid0065OldData(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	resultData := ResultDataOld{}
	err := ascii.Unmarshal(pkg.Body, &resultData)
	if err != nil {
		return err
	}

	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: resultData.ToTighteningResult(),
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// pset详情
func HandleMid0013PsetDetailReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	psetDetail, err := DeserializePSetDetail(pkg.Body)
	if err != nil {
		return err
	}

	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: psetDetail,
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// pset列表
func HandleMid0011PsetListReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	psetList := PSetList{}
	err := psetList.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: psetList,
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// job列表
func HandleMid0031JobListReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	jobList := JobList{}
	err := jobList.Deserialize(pkg.Body)
	if err != nil {
		return err
	}
	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: jobList,
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// job详情
func HandleMid0033JobDetailReply(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	jobDetaill, err := DeserializeJobDetail(pkg.Body)
	if err != nil {
		return err
	}
	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: jobDetaill,
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// 请求错误
func HandleMid0004CmdErr(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	errCode := pkg.Body[4:6]
	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: fmt.Sprintf("Error Code: %s Is Not Defined!", errCode),
	}
	if _, ok := requestErrors[errCode]; ok {
		resp.Body = requestErrors[errCode]
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// 请求成功
func HandleMid0005CmdOk(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: requestErrors["00"],
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}

// job推送信息
func HandleMid0035JobInfo(c *TighteningController, pkg *handlerPkg) error {
	jobInfo := JobInfo{}
	err := ascii.Unmarshal(pkg.Body, &jobInfo)
	if err != nil {
		return err
	}

	// 加入判断，防止重复推送
	if jobInfo.JobStatus == JobInfoNotCompleted &&
		jobInfo.JobBatchCounter == 0 &&
		jobInfo.JobCurrentStep > 0 &&
		jobInfo.JobTotalStep > 0 &&
		jobInfo.JobTighteningStatus == 0 {
		// 推送job选择信息

		jobSelect := tightening_device.JobInfo{
			Job: jobInfo.JobID,
		}

		c.doDispatch(dispatcherbus.DispatcherJob, jobSelect)
	}

	return nil

}

// 控制器输入变化
func HandleMid0211InputMonitor(c *TighteningController, pkg *handlerPkg) error {
	inputs := IOMonitor{}
	err := inputs.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	inputs.ControllerSN = c.SerialNumber()

	c.inputs = inputs.Inputs

	c.NotifyIOContact(io.IoTypeInput, c.inputs)

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
	ids := DeserializeIDS(pkg.Body)
	conf := c.ProtocolService.config()
	bc := ""
	for _, v := range conf.VinIndex {
		if v < 0 || v > (MaxIdsNum-1) {
			continue
		}

		bc += ids[v]
	}
	v := conf.DataIndex
	b, error := strconv.ParseInt(ids[v], 10, 64) // 获取数据段第一个数据，用来判定是否为dummy标识
	if error != nil {
		b = 0
	}
	if c.subscribeVINStatus == BarcodeDone && (b == DummyBarCode) {
		c.diag.Debug("Switch VIN Subscribe Status Waiting For VIN!!!")
		c.subscribeVINStatus = WaitingForBarcode
		return
	}
	if c.subscribeVINStatus == WaitingForBarcode {
		c.getDefaultTransportClient().vinSubscribeBuf <- rxgo.Of(bc)
		c.subscribeVINStatus = BarcodeDone
		return
	}

	return
}

// 报警信息
func HandleMid0071Alarm(c *TighteningController, pkg *handlerPkg) error {
	var ai AlarmInfo
	err := ascii.Unmarshal(pkg.Body, &ai)
	if err != nil {
		return err
	}

	// 参见 项目管理,长安项目中文件:http://116.62.21.97/web#id=325&view_type=form&model=ir.attachment&active_id=3&menu_id=90
	// 第11页,错误代码:Tool calibration required:E305
	//if ai.ErrorCode == "E305" {
	//	// do nothing,当前未确认是否为这个错误代码
	//}

	//switch ai.ErrorCode {
	//case EvtControllerToolConnect:
	//	c.getInstance().UpdateToolStatus(pkg.SeqNumber, device.BaseDeviceStatusOnline)
	//
	//case EvtControllerToolDisconnect:
	//	c.getInstance().UpdateToolStatus(pkg.SeqNumber, device.BaseDeviceStatusOffline)
	//}

	return nil
}

// 报警状态
func HandleMid0076AlarmStatus(c *TighteningController, pkg *handlerPkg) error {
	var as = AlarmStatus{}
	err := ascii.Unmarshal(pkg.Body, &as)
	if err != nil {
		return err
	}

	//switch as.ErrorCode {
	//case EvtControllerNoErr:
	//	c.getInstance().UpdateToolStatus(pkg.SeqNumber, device.BaseDeviceStatusOnline)
	//
	//case EvtControllerToolDisconnect:
	//	c.getInstance().UpdateToolStatus(pkg.SeqNumber, device.BaseDeviceStatusOffline)
	//}

	return nil
}

// 工具状态(维护)
func HandleMid0041ToolInfoReply(c *TighteningController, pkg *handlerPkg) error {
	var ti ToolInfo
	err := ti.Deserialize(pkg.Body)

	if err != nil {
		return err
	}

	if ti.ToolSN == "" {
		return errors.New("Tool Serial Number Is Empty String ")
	}

	if ti.TotalTighteningCount == 0 || ti.CountSinLastService == 0 {
		//不需要尝试创建维修/标定单据
		return nil
	}

	c.doDispatch(dispatcherbus.DispatcherToolMaintenance, ti.ToMaintenanceInfo())
	return nil
}

/*
lexen 扩展
*/

// 处理曲线
func HandleMid7410LastCurveNoAck(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	//讲收到的曲线先做反序列化处理
	var curve = CurveBody{}
	err := ascii.Unmarshal(pkg.Body, &curve)
	if err != nil {
		c.diag.Error("ascii.Unmarshal", err)
	}
	if curve.ToolChannelNumber == 0 {
		e := errors.New("收到的结果曲线数据不合法，未指定工具号")
		c.diag.Error("handleMID_7410_LAST_CURVE", e)
		return e
	}

	torqueCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TorqueString), 64)
	angleCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.AngleString), 64)
	timeCoefficient, _ := strconv.ParseFloat(strings.TrimSpace(curve.TimeString), 64)
	Torque, Angle := c.CurveDataDecoding([]byte(curve.Data), torqueCoefficient, angleCoefficient, c.diag)

	client.tempResultCurve.CUR_M = append(client.tempResultCurve.CUR_M, Torque...)
	client.tempResultCurve.CUR_W = append(client.tempResultCurve.CUR_W, Angle...)
	//当本次数据为本次拧紧曲线的最后一次数据时
	if curve.Num == curve.Id {
		//若取到的点的数量大于协议解析出来该曲线的点数，多出的部分删掉，否则有多少发多少.
		if curve.MeasurePoints < len(client.tempResultCurve.CUR_M) {
			client.tempResultCurve.CUR_M = client.tempResultCurve.CUR_M[0:curve.MeasurePoints]
			client.tempResultCurve.CUR_W = client.tempResultCurve.CUR_W[0:curve.MeasurePoints]
		}

		client.tempResultCurve.GenerateTimeCurveByCoef(float32(timeCoefficient * 1000))

		//本次曲线全部解析完毕后,降临时存储的数据清空
		tool, err := c.getInstance().GetToolViaChannel(curve.ToolChannelNumber)
		if err != nil {
			return err
		}
		//todo: 找到唯一的工具
		sn := tool.SerialNumber()

		//defer delete(client.tempResultCurve, curve.ToolChannelNumber)
		client.tempResultCurve.ToolSN = sn
		client.tempResultCurve.UpdateTime = time.Now()
		c.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherCurve), client.tempResultCurve)
		// dispatch完后创建新的缓存拧紧曲线
		client.tempResultCurve = tightening_device.NewTighteningCurve()
	}
	return nil
}

func HandleMid0011PsetListInfo(c *TighteningController, pkg *handlerPkg) error {
	client := c.getClient(pkg.SN)
	psetList := PSetListWithInfo{}
	err := psetList.Deserialize(pkg.Body)
	if err != nil {
		return err
	}

	resp := &respPkg{
		Seq:  pkg.Seq,
		Body: psetList,
	}
	if client.IsNeedResponse() {
		client.responseChannel <- resp
	}

	return nil
}
