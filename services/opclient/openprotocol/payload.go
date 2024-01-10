package openprotocol

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	OpTerminal = 0x00

	OpenprotocolModeJob  = "1"
	OpenprotocolModePset = "0"
)

const (
	JobInfoNotCompleted = 0
	MaxIdsNum           = 4
)

//const (
//	EvtControllerNoErr          = "E000"
//	EvtControllerToolDisconnect = "I003"
//	EvtControllerToolConnect    = "I002"
//)

type OpenProtocolParams struct {
	MaxKeepAliveCheck int
	MaxReplyTime      time.Duration
	KeepAlivePeriod   time.Duration
}

type IOStatus struct {
	No     int    `json:"no"`
	Status string `json:"status"`
}

const (
	LenHeader        = 20
	DefaultRev       = "000"
	LenSingleSpindle = 18
	IoModel          = "IO_MODEL"

	MID_0001_START                        = "0001"
	MID_0002_START_ACK                    = "0002"
	MID_0003_STOP                         = "0003"
	MID_0004_CMD_ERR                      = "0004"
	MID_0005_CMD_OK                       = "0005"
	MID_0014_PSET_SUBSCRIBE               = "0014"
	MID_0018_PSET                         = "0018"
	MID_0012_PSET_DETAIL_REQUEST          = "0012"
	MID_0013_PSET_DETAIL_REPLY            = "0013"
	MID_0010_PSET_LIST_REQUEST            = "0010"
	MID_0011_PSET_LIST_REPLY              = "0011"
	MID_0034_JOB_INFO_SUBSCRIBE           = "0034"
	MID_0040_TOOL_INFO_REQUEST            = "0040"
	MID_0041_TOOL_INFO_REPLY              = "0041"
	MID_0060_LAST_RESULT_SUBSCRIBE        = "0060"
	MID_0062_LAST_RESULT_ACK              = "0062"
	MID_7402_Cycle_Phase_Result_Subscribe = "7402"
	MID_7408_LAST_CURVE_SUBSCRIBE         = "7408"
	MID_7411_LAST_CURVE_DATA_ACK          = "7411"
	MID_0151_IDENTIFIER_SUBSCRIBE         = "0151"
	MID_0150_IDENTIFIER_SET               = "0150"
	MID_0038_JOB_SELECT                   = "0038"
	MID_0064_OLD_SUBSCRIBE                = "0064"
	MID_0065_OLD_DATA                     = "0065"
	MID_0070_ALARM_SUBSCRIBE              = "0070"
	MID_0071_ALARM                        = "0071"
	MID_0076_ALARM_STATUS                 = "0076"

	MID_0080_TIME_SUBSCRIBE = "0080"
	MID_0081_TIME           = "0081"

	MID_0130_JOB_OFF                 = "0130"
	MID_0250_SELECTOR_SUBSCRIBE      = "0250"
	MID_0042_TOOL_DISABLE            = "0042"
	MID_0043_TOOL_ENABLE             = "0043"
	MID_0030_JOB_LIST_REQUEST        = "0030"
	MID_0031_JOB_LIST_REPLY          = "0031"
	MID_0032_JOB_DETAIL_REQUEST      = "0032"
	MID_0033_JOB_DETAIL_REPLY        = "0033"
	MID_0200_CONTROLLER_RELAYS       = "0200"
	MID_0019_PSET_BATCH_SET          = "0019"
	MID_0020_PSET_BATCH_RESET        = "0020"
	MID_0035_JOB_INFO                = "0035"
	MID_0210_INPUT_SUBSCRIBE         = "0210"
	MID_0211_INPUT_MONITOR           = "0211"
	MID_0127_JOB_ABORT               = "0127"
	MID_0100_MULTI_SPINDLE_SUBSCRIBE = "0100"
	MID_0101_MULTI_SPINDLE_RESULT    = "0101"
	MID_0051_VIN_SUBSCRIBE           = "0051"
	MID_0052_VIN                     = "0052"

	MID_0008_DATA_SUB = "0008"

	MID_0061_LAST_RESULT  = "0061"
	MID_7410_LAST_CURVE   = "7410"
	MID_7404_PHASE_RESULT = "7404"
	MID_7406_CYCLE_RESULT = "7406"
	MID_9999_ALIVE        = "9999"
)

const (
	ROTATION_CW  = "CW"
	ROTATION_CCW = "CCW"
)

const (
	MID_0038_REV_1 = "001"
	MID_0038_REV_2 = "002"
)

var requestErrors = map[string]string{
	"00": "OK",
	"01": "Invalid data",
	"02": "Parameter set ID not present",
	"03": "Parameter set can not be set",
	"04": "Parameter set not running",
	"06": "VIN upload subscription already exists",
	"07": "VIN upload subscription does not exists",
	"08": "VIN input source not granted",
	"09": "Last tightening result subscription already exists",
	"10": "Last tightening result subscription does not exist",
	"11": "Alarm subscription already exists",
	"12": "Alarm subscription does not exist",
	"13": "Parameter set selection subscription already exists",
	"14": "Parameter set selection subscription does not exist",
	"15": "Tightening ID requested not found",
	"16": "Connection rejected protocol busy",
	"17": "Job ID not present",
	"18": "Job info subscription already exists",
	"19": "Job info subscription does not exist",
	"20": "Job can not be set",
	"21": "Job not running",
	"22": "Not possible to execute dynamic Job request",
	"23": "Job batch decrement failed",
	"24": "Not possible to create Pset",
	"25": "Programming control not granted",
	"30": "TighteningController is not a sync Master/station controller",
	"31": "Multi-spindle status subscription already exists",
	"32": "Multi-spindle status subscription does not exist",
	"33": "Multi-spindle result subscription already exists",
	"34": "Multi-spindle result subscription does not exist",
	"40": "Job line control info subscription already exists",
	"41": "Job line control info subscription does not exist",
	"42": "Identifier input source not granted",
	"43": "Multiple identifiers work order subscription already exists",
	"44": "Multiple identifiers work order subscription does not exist",
	"50": "Status external monitored inputs subscription already exists",
	"51": "Status external monitored inputs subscription does not exist",
	"52": "io tightening_device not connected",
	"53": "Faulty io tightening_device ID",
	"54": "Tool Tag ID unknown",
	"55": "Tool Tag ID subscription already exists",
	"56": "Tool Tag ID subscription does not exist",
	"58": "No alarm present",
	"59": "Tool currently in use",
	"60": "No histogram available",
	"70": "Calibration failed",
	"79": "Command failed",
	"80": "Audi emergency status subscription exists",
	"81": "Audi emergency status subscription does not exist",
	"82": "Automatic/Manual mode subscribe already exist",
	"83": "Automatic/Manual mode subscribe does not exist",
	"84": "The relay function subscription already exists",
	"85": "The relay function subscription does not exist",
	"86": "The selector socket info subscription already exist",
	"87": "The selector socket info subscription does not exist",
	"88": "The digin info subscription already exist",
	"89": "The digin info subscription does not exist",
	"90": "Lock at bach done subscription already exist",
	"91": "Lock at bach done subscription does not exist",
	"92": "Open protocol commands disabled",
	"93": "Open protocol commands disabled subscription already exists",
	"94": "Open protocol commands disabled subscription does not exist",
	"95": "Reject request, PowerMACS is in manual mode",
	"96": "Client already connected",
	"97": "MID revision unsupported",
	"98": "TighteningController internal request timeout",
	"99": "Unknown MID",
}

type OpenProtocolHeader struct {
	LEN      int
	MID      string
	Revision string
	NoAck    string
	Station  string
	Spindle  string
	Spare    string
}

func (h *OpenProtocolHeader) Serialize() string {
	return fmt.Sprintf("%04d%04s%03s%-1s%-2s%-2s%-4s", h.LEN, h.MID, h.Revision, h.NoAck, h.Station, h.Spindle, h.Spare)
}

func (h *OpenProtocolHeader) Deserialize(str string) {
	if len(str) != LenHeader {
		return
	}

	n, _ := strconv.ParseInt(str[0:4], 10, 32)
	h.LEN = int(n) - LenHeader
	h.MID = str[4:8]
	h.Revision = str[8:11]
	h.NoAck = str[11:12]
	h.Station = str[12:14]
	h.Spindle = str[14:16]
	h.Spare = str[16:20]
}

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

type IOMonitor struct {
	ControllerSN string `json:"controller_sn"`
	Inputs       string `json:"inputs"`
}

func (iom *IOMonitor) Deserialize(str string) error {

	iom.Inputs = str

	return nil
}

//type baseResultData struct {
//	IResultData
//	CellID         int    `start:"3"  end:"6"`
//	ChannelID      int    `start:"9"  end:"10"`
//	ControllerName string `start:"13"  end:"37"`
//	VIN            string `start:"40"  end:"64"`
//	JobID          int    `start:"67"  end:"70"`
//	PSetID         int    `start:"73"  end:"75"`
//}

type CycleResultData struct {
	//rev1
	ChannelID   int `start:"3"  end:"4"`
	CycleNumber int `start:"7"  end:"9"`

	VIN       string `start:"12"  end:"36"`
	TimeStamp string `start:"39"  end:"57"`

	ControllerName   string `start:"60"  end:"70"`
	ToolSerialNumber string `start:"73"  end:"83"`
	TighteningID     string `start:"86"  end:"95"`
	ID1              string `start:"98"  end:"107"`
	BatchSize        int    `start:"110"  end:"112"`
	BatchCount       int    `start:"115"  end:"117"`
}

type PhaseResult struct {
	PhaseTorque []float64
	PhaseAngle  []float64
	Phase       []PhaseResultData
}
type PhaseResultData struct {
	ChannelID         int     `start:"3"  end:"4"`
	CycleNumber       int     `start:"7"  end:"9"`
	PhaseNumber       int     `start:"12"  end:"13"`
	PhaseName         string  `start:"16"  end:"17"`
	MethodName        string  `start:"20"  end:"21"`
	TorqueMin         float64 `start:"24"  end:"29"`
	TorqueMax         float64 `start:"32"  end:"37"`
	TorqueFinalTarget float64 `start:"56"  end:"61"`
	Torque            float64 `start:"64"  end:"69"`
	AngleMin          float64 `start:"74"  end:"77"`
	AngleMax          float64 `start:"80"  end:"85"`
	FinalAngleTarget  float64 `start:"88"  end:"93"`
	Angle             float64 `start:"96"  end:"101"`
}

// 0061 Result payload
type ResultDataRev1 struct {
	CellID         int    `start:"3"  end:"6"`
	ChannelID      int    `start:"9"  end:"10"`
	ControllerName string `start:"13"  end:"37"`
	VIN            string `start:"40"  end:"64"`
	JobID          int    `start:"67"  end:"68"`
	PSetID         int    `start:"71"  end:"73"`

	BatchSize         int     `start:"76"  end:"79"`
	BatchCount        int     `start:"82"  end:"85"`
	TighteningStatus  string  `start:"88"  end:"88"`
	TorqueStatus      string  `start:"91"  end:"91"`
	AngleStatus       string  `start:"94"  end:"94"`
	TorqueMin         float64 `start:"97"  end:"102"`
	TorqueMax         float64 `start:"105"  end:"110"`
	TorqueFinalTarget float64 `start:"113"  end:"118"`
	Torque            float64 `start:"121"  end:"126"`
	AngleMin          float64 `start:"129"  end:"133"`
	AngleMax          float64 `start:"136"  end:"140"`
	FinalAngleTarget  float64 `start:"143"  end:"147"`
	Angle             float64 `start:"150"  end:"154"`

	TimeStamp    string `start:"157"  end:"175"`
	BatchStatus  string `start:"199"  end:"199"`
	TighteningID string `start:"202"  end:"211"`
}

func (rd *ResultDataRev1) ToTighteningResult() *tightening_device.TighteningResult {
	measureResult := tightening_device.RESULT_OK
	if rd.TighteningStatus == "0" {
		measureResult = tightening_device.RESULT_NOK
	}

	strategy := ""

	dt := time.Now()
	//if rd.TimeStamp != "" {
	//	dts := []byte(rd.TimeStamp)
	//	dts[10] = 'T'
	//	dts = append(dts, 'Z')
	//	dt, _ = time.Parse(time.RFC3339, string(dts))
	//}

	errs, _ := parseOpenProtocolErrorCode("00")
	r := tightening_device.TighteningResult{
		BaseResult: tightening_device.BaseResult{
			MeasureResult: measureResult,
			MeasureTorque: rd.Torque / 100,
			MeasureAngle:  rd.Angle,
		},

		JobInfo: tightening_device.JobInfo{
			Job: rd.JobID,
		},

		//错误信息解析
		ErrorCode: errs,

		// 工具通道号
		ChannelID: rd.ChannelID,

		// 收到时间
		UpdateTime: dt,

		// pset号
		PSet: rd.PSetID,

		// 拧紧ID
		TighteningID: rd.TighteningID,

		// 拧紧策略
		Strategy: strategy,

		// 最大扭矩
		TorqueMax: rd.TorqueMax / 100,

		// 最小扭矩
		TorqueMin: rd.TorqueMin / 100,

		// 扭矩阈值
		TorqueThreshold: rd.TorqueMax / 100,

		// 目标扭矩
		TorqueTarget: rd.TorqueFinalTarget / 100,

		// 最大角度
		AngleMax: rd.AngleMax,

		// 最小角度
		AngleMin: rd.AngleMin,

		// 目标角度
		AngleTarget: rd.FinalAngleTarget,

		ControllerName: strings.TrimSpace(rd.ControllerName),

		BatchCount: rd.BatchCount,
	}

	r.StepResults = []tightening_device.StepData{
		tighteningResult2StepData(&r),
	}

	return &r
}

func (rd *ResultDataRev1) GetInstance() interface{} {
	return rd
}

type ResultData struct {
	//rev2
	CellID         int    `start:"3"  end:"6"`
	ChannelID      int    `start:"9"  end:"10"`
	ControllerName string `start:"13"  end:"37"`
	VIN            string `start:"40"  end:"64"`
	JobID          int    `start:"67"  end:"70"`
	PSetID         int    `start:"73"  end:"75"`
	Strategy       string `start:"78"  end:"79"`
	//StrategyOption                []byte	`start:"3"  end:"6"`
	BatchSize                     int     `start:"89"  end:"92"`
	BatchCount                    int     `start:"95"  end:"98"`
	TighteningStatus              string  `start:"101"  end:"101"`
	BatchStatus                   string  `start:"104"  end:"104"`
	TorqueStatus                  string  `start:"107"  end:"107"`
	AngleStatus                   string  `start:"110"  end:"110"`
	RundownAngleStatus            string  `start:"113"  end:"113"`
	CurrentMonitoringStatus       string  `start:"116"  end:"116"`
	SelftapStatus                 string  `start:"119"  end:"119"`
	PrevailTorqueMonitoringStatus string  `start:"122"  end:"122"`
	PrevailTorqueCompensateStatus string  `start:"125"  end:"125"`
	TighteningErrorStatus         string  `start:"128"  end:"137"`
	TorqueMin                     float64 `start:"140"  end:"145"`
	TorqueMax                     float64 `start:"148"  end:"153"`
	TorqueFinalTarget             float64 `start:"156"  end:"161"`
	Torque                        float64 `start:"164"  end:"169"`
	AngleMin                      float64 `start:"172"  end:"176"`
	AngleMax                      float64 `start:"179"  end:"183"`
	FinalAngleTarget              float64 `start:"186"  end:"190"`
	Angle                         float64 `start:"193"  end:"197"`
	//RundownAngleMin               float64	`start:"3"  end:"6"`
	//RundownAngleMax               float64	`start:"3"  end:"6"`
	//RundownAngle                  float64	`start:"130"  end:"134"`
	//CurrentMonitoringMin          float64	`start:"3"  end:"6"`
	//CurrentMonitoringMax          float64	`start:"3"  end:"6"`
	//CurrentMonitoring             float64	`start:"137"  end:"139"`
	//SelftapMin                    float64	`start:"3"  end:"6"`
	//SelftapMax                    float64	`start:"3"  end:"6"`
	//SelftapTorque                 float64	`start:"142"  end:"147"`
	//PrevailTorqueMonitoringMin    float64	`start:"3"  end:"6"`
	//PrevailTorqueMonitoringMax    float64	`start:"3"  end:"6"`
	//PrevailTorque                 float64	`start:"150"  end:"155"`
	TighteningID string `start:"284"  end:"293"`
	//JobSequenceNumber             int32		`start:"3"  end:"6"`
	//SyncTighteningID              int32		`start:"3"  end:"6"`
	ToolSerialNumber string `start:"310"  end:"323"`
	TimeStamp        string `start:"326"  end:"344"`
	//TimeStampPSetLastChange       string	`start:"3"  end:"6"`

	//rev3
	//PSetName   string	`start:"3"  end:"6"`
	TorqueUnit string `start:"395"  end:"395"`
	ResultType string `start:"398"  end:"399"`

	//rev4
	ID2 string `start:"402"  end:"426"`
	ID3 string `start:"429"  end:"453"`
	ID4 string `start:"456"  end:"480"`

	//rev5
	//CustomerErrorCode string	`start:"3"  end:"6"`

	//rev6

	//rev998
	NumberOfStages int `start:"3"  end:"6"`
	//NumberOfStageResults int	`start:"3"  end:"6"`
	StageResult string `start:"3"  end:"..."`
}

func NewResultData(rev string) IResultData {
	switch rev {
	case "001":
		return &ResultDataRev1{}
	case "998":
		return &ResultData{}
	default:
		return &ResultData{}
	}
}

func parseOpenProtocolErrorCode(errors string) (string, error) {
	valErrors, err := strconv.ParseInt(errors, 10, 64)
	if err != nil {
		return "", err
	}

	strErrors := strconv.FormatInt(valErrors, 2)
	return fmt.Sprintf("%032s", strErrors), nil
}

func (rd *CycleResultData) ToTighteningResult(p *PhaseResult) tightening_device.TighteningResult {
	measureResult := tightening_device.RESULT_OK
	if rd.BatchSize != rd.BatchCount {
		measureResult = tightening_device.RESULT_NOK
	}

	strategy := ""

	dt := time.Now()
	if rd.TimeStamp != "" {
		dts := []byte(rd.TimeStamp)
		dts[10] = 'T'
		dts = append(dts, 'Z')
		dt, _ = time.Parse(time.RFC3339, string(dts))
	}

	errs, _ := parseOpenProtocolErrorCode("00")
	r := tightening_device.TighteningResult{
		BaseResult: tightening_device.BaseResult{
			ToolSN:        strings.TrimSpace(rd.ToolSerialNumber),
			MeasureResult: measureResult,
			MeasureTorque: 0,
			MeasureAngle:  0,
		},

		JobInfo: tightening_device.JobInfo{
			Job: 0,
		},

		//错误信息解析
		ErrorCode: errs,

		// 工具通道号
		ChannelID: rd.ChannelID,

		// 收到时间
		UpdateTime: dt,

		// pset号
		PSet: 0,

		// 拧紧ID
		TighteningID: rd.TighteningID,

		// 拧紧策略
		Strategy: strategy,

		// 最大扭矩
		TorqueMax: 0,

		// 最小扭矩
		TorqueMin: 0,

		// 扭矩阈值
		TorqueThreshold: 0,

		// 目标扭矩
		TorqueTarget: 0,

		// 最大角度
		AngleMax: 0,

		// 最小角度
		AngleMin: 0,

		// 目标角度
		AngleTarget: 0,

		ControllerName: strings.TrimSpace(rd.ControllerName),

		BatchCount: rd.BatchCount,
	}

	r.StepResults = []tightening_device.StepData{}
	for _, phase := range p.Phase {
		cvi2R := tighteningResult2StepDataCvi2(&phase)
		r.StepResults = append(r.StepResults, cvi2R)
	}
	return r
}

func (rd *ResultData) ToTighteningResult() *tightening_device.TighteningResult {
	measureResult := tightening_device.RESULT_OK
	if rd.TighteningStatus == "0" {
		measureResult = tightening_device.RESULT_NOK
	}

	strategy := ""

	switch rd.Strategy {
	case "01":
		strategy = tightening_device.STRATEGY_AW

	case "02":
		strategy = tightening_device.STRATEGY_AW

	case "03":
		strategy = tightening_device.STRATEGY_ADW

	case "04":
		strategy = tightening_device.STRATEGY_AD
	}

	if rd.ResultType == "02" {
		measureResult = tightening_device.RESULT_LSN
		strategy = tightening_device.STRATEGY_LN
	}

	dt := time.Now()
	if rd.TimeStamp != "" {
		dts := []byte(rd.TimeStamp)
		dts[10] = 'T'
		dts = append(dts, 'Z')
		dt, _ = time.Parse(time.RFC3339, string(dts))
	}

	errs, _ := parseOpenProtocolErrorCode(rd.TighteningErrorStatus)
	r := tightening_device.TighteningResult{
		BaseResult: tightening_device.BaseResult{
			ToolSN:        strings.TrimSpace(rd.ToolSerialNumber),
			MeasureResult: measureResult,
			MeasureTorque: rd.Torque / 100,
			MeasureAngle:  rd.Angle,
		},

		JobInfo: tightening_device.JobInfo{
			Job: rd.JobID,
		},

		//错误信息解析
		ErrorCode: errs,

		// 工具通道号
		ChannelID: rd.ChannelID,

		// 收到时间
		UpdateTime: dt,

		// pset号
		PSet: rd.PSetID,

		// 拧紧ID
		TighteningID: rd.TighteningID,

		// 拧紧策略
		Strategy: strategy,

		// 最大扭矩
		TorqueMax: rd.TorqueMax / 100,

		// 最小扭矩
		TorqueMin: rd.TorqueMin / 100,

		// 扭矩阈值
		TorqueThreshold: rd.TorqueMax / 100,

		// 目标扭矩
		TorqueTarget: rd.TorqueFinalTarget / 100,

		// 最大角度
		AngleMax: rd.AngleMax,

		// 最小角度
		AngleMin: rd.AngleMin,

		// 目标角度
		AngleTarget: rd.FinalAngleTarget,

		ControllerName: strings.TrimSpace(rd.ControllerName),

		BatchCount: rd.BatchCount,
	}

	r.StepResults = []tightening_device.StepData{
		tighteningResult2StepData(&r),
	}

	return &r
}

func (rd *ResultData) GetInstance() interface{} {
	return rd
}

func tighteningResult2StepDataCvi2(result *PhaseResultData) tightening_device.StepData {
	var r tightening_device.StepData
	r.PSetDefine = tightening_device.PSetDefine{
		Mp: result.TorqueMax,
		Mm: result.TorqueMin,
		Ms: 0.0, //没有阈值扭矩
		Ma: result.Torque,
		Wp: result.AngleMax,
		Wm: result.AngleMin,
		Wa: result.Angle,
	}
	return r
}

func tighteningResult2StepData(result *tightening_device.TighteningResult) tightening_device.StepData {
	var r tightening_device.StepData
	r.TighteningResult = *result
	r.Strategy = result.Strategy
	r.PSetDefine = tightening_device.PSetDefine{
		Mp: result.TorqueMax,
		Mm: result.TorqueMin,
		Ms: 0.0, //没有阈值扭矩
		Ma: result.TorqueTarget,
		Wp: result.AngleMax,
		Wm: result.AngleMin,
		Wa: result.AngleTarget,
	}
	return r
}

type ResultDataOld struct {

	//rev2
	VIN                           string  `start:"15"  end:"39"`
	JobID                         int     `start:"42"  end:"45"`
	PSetID                        int     `start:"48"  end:"50"`
	Strategy                      string  `start:"53"  end:"54"`
	BatchSize                     int     `start:"64"  end:"67"`
	BatchCount                    int     `start:"70"  end:"73"`
	TighteningStatus              string  `start:"76"  end:"76"`
	BatchStatus                   string  `start:"79"  end:"79"`
	TorqueStatus                  string  `start:"82"  end:"82"`
	AngleStatus                   string  `start:"85"  end:"85"`
	RundownAngleStatus            string  `start:"88"  end:"88"`
	CurrentMonitoringStatus       string  `start:"91"  end:"91"`
	SelftapStatus                 string  `start:"94"  end:"94"`
	PrevailTorqueMonitoringStatus string  `start:"97"  end:"97"`
	PrevailTorqueCompensateStatus string  `start:"100"  end:"100"`
	Torque                        float64 `start:"115"  end:"120"`
	Angle                         float64 `start:"123"  end:"127"`
	RundownAngle                  float64 `start:"130"  end:"134"`
	CurrentMonitoring             float64 `start:"137"  end:"139"`
	SelftapTorque                 float64 `start:"142"  end:"147"`
	PrevailTorque                 float64 `start:"150"  end:"155"`
	TightingID                    string  `start:"3"  end:"12"`
	ToolSerialNumber              string  `start:"172"  end:"285"`
	TimeStamp                     string  `start:"188"  end:"206"`

	//rev3
	TorqueUnit string `start:"209"  end:"209"`
	ResultType string `start:"212"  end:"213"`

	//rev4
	ID2 string `start:"216"  end:"240"`
	ID3 string `start:"243"  end:"267"`
	ID4 string `start:"270"  end:"294"`

	//rev5

	//rev6

	//rev998
	StageResult string `start:"3"  end:"..."`
}

func (ord *ResultDataOld) ToTighteningResult() tightening_device.TighteningResult {
	measureResult := tightening_device.RESULT_OK
	if ord.TighteningStatus == "0" {
		measureResult = tightening_device.RESULT_NOK
	}

	strategy := ""

	switch ord.Strategy {
	case "01":
		strategy = tightening_device.STRATEGY_AW

	case "02":
		strategy = tightening_device.STRATEGY_AW

	case "03":
		strategy = tightening_device.STRATEGY_ADW

	case "04":
		strategy = tightening_device.STRATEGY_AD
	}

	if ord.ResultType == "02" {
		measureResult = tightening_device.RESULT_LSN
		strategy = tightening_device.STRATEGY_LN
	}

	return tightening_device.TighteningResult{
		BaseResult: tightening_device.BaseResult{
			ToolSN:        strings.TrimSpace(ord.ToolSerialNumber),
			MeasureResult: measureResult,
			MeasureTorque: ord.Torque / 100,
			MeasureAngle:  ord.Angle,
		},

		JobInfo: tightening_device.JobInfo{
			Job: ord.JobID,
		},

		// 收到时间
		UpdateTime: time.Now(),

		// pset号
		PSet: ord.PSetID,

		// 拧紧ID
		TighteningID: ord.TightingID,

		// 拧紧策略
		Strategy: strategy,
	}
}

func DeserializePSetDetail(str string) (*tightening_device.PSetDetail, error) {
	var p tightening_device.PSetDetail
	var err error

	p.PSetID, err = strconv.Atoi(str[2:5])
	if err != nil {
		return nil, err
	}

	p.PSetName = strings.TrimSpace(str[7:32])
	p.RotationDirection = str[34:35]

	switch p.RotationDirection {
	case "1":
		p.RotationDirection = ROTATION_CW

	case "2":
		p.RotationDirection = ROTATION_CCW
	}

	p.BatchSize, err = strconv.Atoi(str[37:39])
	if err != nil {
		return nil, err
	}

	p.TorqueMin, err = strconv.ParseFloat(str[41:47], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueMin = p.TorqueMin / 100

	p.TorqueMax, err = strconv.ParseFloat(str[49:55], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueMax = p.TorqueMax / 100

	p.TorqueTarget, err = strconv.ParseFloat(str[57:63], 64)
	if err != nil {
		return nil, err
	}

	p.TorqueTarget = p.TorqueTarget / 100

	p.AngleMin, err = strconv.ParseFloat(str[65:70], 64)
	if err != nil {
		return nil, err
	}

	p.AngleMax, err = strconv.ParseFloat(str[72:77], 64)
	if err != nil {
		return nil, err
	}

	p.AngleTarget, err = strconv.ParseFloat(str[79:84], 64)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

type PSetListWithInfo struct {
	num   int
	psets []tightening_device.PSetInfo
}

func (p *PSetListWithInfo) Deserialize(str string) error {
	infoLen := 18 // 3char pset ID + 15 char pset name
	nums := len(str) / infoLen
	for i := 0; i < nums; i++ {
		if (i*infoLen)+infoLen > len(str) {
			continue
		}
		item := str[i*infoLen : (i*infoLen)+infoLen]
		pId, err := strconv.Atoi(item[:3])
		if err != nil {
			continue
		}
		p.psets = append(p.psets, tightening_device.PSetInfo{
			Name: item[3:],
			ID:   pId,
		})
	}

	return nil
}

type PSetList struct {
	num   int
	psets []tightening_device.PSetInfo
}

func (p *PSetList) Deserialize(str string) error {
	var err error
	p.num, err = strconv.Atoi(str[0:3])
	if err != nil {
		return err
	}

	for i := 0; i < p.num; i++ {
		pset, _ := strconv.Atoi(str[(i+1)*3 : (i+1)*3+3])
		p.psets = append(p.psets, tightening_device.PSetInfo{
			ID: pset,
		})
	}

	return nil
}

type JobList struct {
	num  int
	jobs []int
}

func (p *JobList) Deserialize(str string) error {
	var err error
	p.num, err = strconv.Atoi(str[0:4])
	if err != nil {
		return err
	}

	for i := 0; i < p.num; i++ {
		job, _ := strconv.Atoi(str[(i+1)*4 : (i+1)*4+4])
		p.jobs = append(p.jobs, job)
	}

	return nil
}

func DeserializeJobDetail(str string) (*tightening_device.JobDetail, error) {
	var err error
	var p tightening_device.JobDetail

	p.JobID, err = strconv.Atoi(str[2:6])
	if err != nil {
		return nil, err
	}

	p.JobName = strings.TrimSpace(str[8:33])

	order := str[35:36]
	switch order {
	case "0":
		p.OrderStrategy = "free"

	case "1":
		p.OrderStrategy = "forced"

	case "2":
		p.OrderStrategy = "free and forced"
	}

	countType := str[51:52]
	switch countType {
	case "0":
		p.CountType = "only the OK tightenings are counted"

	case "1":
		p.CountType = "both the OK and NOK tightenings are counted"
	}

	if str[54:55] == "0" {
		p.LockAtJobDone = false
	} else {
		p.LockAtJobDone = true
	}

	if str[57:58] == "0" {
		p.UseLineControl = false
	} else {
		p.UseLineControl = true
	}

	if str[60:61] == "0" {
		p.RepeatJob = false
	} else {
		p.RepeatJob = true
	}

	loosening, err := strconv.Atoi(str[63:65])
	if err != nil {
		return nil, err
	}

	switch loosening {
	case 0:
		p.LooseningStrategy = "enable"

	case 1:
		p.LooseningStrategy = "disable"

	case 2:
		p.LooseningStrategy = "enable only on NOK tightening"
	}

	stepStr := str[75 : len(str)-1]
	steps := strings.Split(stepStr, ";")
	jobStep := tightening_device.JobStep{}
	for _, v := range steps {
		values := strings.Split(v, ":")

		jobStep.ChannelID, _ = strconv.Atoi(values[0])
		jobStep.PSetID, _ = strconv.Atoi(values[1])
		jobStep.BatchSize, _ = strconv.Atoi(values[3])
		jobStep.Socket, _ = strconv.Atoi(values[4])
		jobStep.StepName = strings.TrimSpace(values[5])

		p.Steps = append(p.Steps, jobStep)
	}

	return &p, nil
}

type AlarmInfo struct {
	ErrorCode      string `start:"3"  end:"6"`
	IsControllerOK bool   `start:"9"  end:"9"`
	IsToolOK       bool   `start:"12"  end:"12"`
}

type ToolInfo struct {
	ToolSN               string `json:"serial_no"`
	ControllerSN         string `json:"controller_sn"`
	TotalTighteningCount int    `json:"times"`
	CountSinLastService  int    `json:"sin_last_service"`
}

func (ti *ToolInfo) Deserialize(msg string) (err error) {

	if len(msg) < 20 {
		return errors.New("tool info msg len is not enough")
	}

	ti.ToolSN = strings.TrimSpace(msg[2:16]) //去除空格获取序列号

	ti.TotalTighteningCount, err = strconv.Atoi(msg[18:28])
	if err != nil {
		return errors.Wrap(err, "tool info Convert TotalTighteningCount Error")
	}

	ti.ControllerSN = strings.TrimSpace(msg[51:61])

	ti.CountSinLastService, err = strconv.Atoi(msg[92:102])

	if err != nil {
		return errors.Wrap(err, "tool info Convert CountSinLastService Error")
	}

	return
}

func (ti *ToolInfo) ToMaintenanceInfo() *tightening_device.ToolMaintenanceInfo {
	return &tightening_device.ToolMaintenanceInfo{
		ToolSN:               ti.ToolSN,
		ControllerSN:         ti.ControllerSN,
		TotalTighteningCount: ti.TotalTighteningCount,
		CountSinLastService:  ti.CountSinLastService,
	}
}

type JobInfo struct {
	JobID               int    `start:"3"  end:"6"`
	JobStatus           int    `start:"9"  end:"9"`
	JobBatchMode        int    `start:"12"  end:"12"`
	JobBatchSize        int    `start:"15"  end:"17"`
	JobBatchCounter     int    `start:"21"  end:"24"`
	Timestamp           string `start:"27"  end:"45"`
	JobCurrentStep      int    `start:"48"  end:"50"`
	JobTotalStep        int    `start:"53"  end:"55"`
	JobStepType         int    `start:"58"  end:"59"`
	JobTighteningStatus int    `start:"62"  end:"63"`
}

func DeserializeIDS(str string) []string {
	var rt []string

	vin := strings.TrimSpace(str[2:27])
	rt = append(rt, vin)

	id2 := strings.TrimSpace(str[29:54])
	rt = append(rt, id2)

	id3 := strings.TrimSpace(str[56:71])
	rt = append(rt, id3)

	id4 := strings.TrimSpace(str[73:98])
	rt = append(rt, id4)

	return rt
}

type SingleSpindleResult struct {
	SpindleNo int
	ChannelID int
	Result    string
	Torque    float64
	Angle     float64
}

type MultiSpindleResult struct {
	TotalSpindleNumber int
	Vin                string
	JobID              int
	PSetID             int
	BatchSize          int
	BatchCount         int
	BatchStatus        int
	TorqueMin          float64
	TorqueMax          float64
	TorqueFinalTarget  float64
	AngleMin           float64
	AngleMax           float64
	FinalAngleTarget   float64

	Spindles []SingleSpindleResult
}

func (msr *MultiSpindleResult) Deserialize(str string) {

	sps := str[154:]

	spNum := len(sps) / LenSingleSpindle
	sp := SingleSpindleResult{}
	for i := 0; i < spNum; i++ {
		targetSp := sps[i*LenSingleSpindle : i*LenSingleSpindle+LenSingleSpindle]
		sp.SpindleNo, _ = strconv.Atoi(targetSp[0:2])
		if targetSp[4:5] == "0" {
			sp.Result = "NOK"
		} else {
			sp.Result = "OK"
		}

		sp.Torque, _ = strconv.ParseFloat(targetSp[6:12], 64)
		sp.Torque = sp.Torque / 100

		sp.Angle, _ = strconv.ParseFloat(targetSp[13:LenSingleSpindle], 64)

		msr.Spindles = append(msr.Spindles, sp)
	}
}

type AlarmStatus struct {
	Status         string `start:"3"  end:"3"`
	ErrorCode      string `start:"6"  end:"9"`
	IsControllerOK bool   `start:"12"  end:"12"`
	IsToolOK       bool   `start:"15"  end:"15"`
}

type CurveBody struct {
	ToolChannelNumber int    `start:"3"  end:"4"`
	TimeString        string `start:"12"  end:"25"`
	TorqueString      string `start:"28"  end:"41"`
	AngleString       string `start:"44"  end:"57"`
	MeasurePoints     int    `start:"60"  end:"63"`
	Num               string `start:"66"  end:"67"` //曲线总共分几段
	Id                string `start:"70"  end:"71"` //当前为第几段
	Data              string `start:"72"  end:"..."`
}

type ControllerSn struct {
	SN string `start:"134"  end:"143"` //SN
}
