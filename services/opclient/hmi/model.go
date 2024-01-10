package hmi

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"github.com/pkg/errors"
	"time"
)

const (
	HMINotifyInfo    = "Info"
	HMINotifyWarning = "Warn"
	HMINotifyError   = "Error"
)

type HMICommonResponse = httpd.HMICommonResponse

type listWorkordersURLParams struct {
	Status          string `param:"status" validate:"required"`
	HMISerialNumber string `param:"hmi_sn" validate:"-"`
	WorkCenterCode  string `param:"workcenter_code" validate:"-"`
}

type WorkSiteLeaving struct {
	WorkcenterCode string    `json:"workcenter_code"`
	WorkorderCode  string    `json:"workorder_code"`
	ExecuteTime    time.Time `json:"execute_time"`
}

type LocalResults struct {
	Id           interface{} `json:"id,omitempty"`
	HmiSN        interface{} `json:"hmi_sn,omitempty"`
	Vin          interface{} `json:"vin,omitempty"`
	VehicleType  interface{} `json:"vehicle_type,omitempty"`
	JobID        interface{} `json:"job_id,omitempty"`
	PSetID       interface{} `json:"pset_id,omitempty"`
	ControllerSN interface{} `json:"controller_sn,omitempty"`
	ToolSN       interface{} `json:"tool_sn,omitempty"`
	Result       interface{} `json:"result,omitempty"`
	Torque       interface{} `json:"torque,omitempty"`
	Angle        interface{} `json:"angle,omitempty"`
	Spent        interface{} `json:"spent,omitempty"`
	TimeStamp    interface{} `json:"timestamp,omitempty"`
	Batch        interface{} `json:"batch,omitempty"`
	HasUpload    interface{} `json:"has_upload"`
	Type         string      `json:"type"`
}

type WSOrderReq struct {
	ID            int64  `json:"id"`
	WorkorderCode string `json:"workorder_code"`
	Status        string `json:"status"`
}

type WSOrderReqData struct {
	ID            int64  `json:"id"`
	WorkorderCode string `json:"workorder_code"`
	Data          string `json:"data"`
}

type WSOrderConflictReq struct {
	ID             int64  `json:"id" validate:"required"`
	WorkCenterCode string `json:"workcenter_code"`
	HmiSn          string `json:"hmi_sn"`
}

type WSNextOrderReq = WSOrderConflictReq

type WSOrderConflictResp struct {
	WSNextOrderResp
	Conflict bool `json:"conflict"`
}

type WSNextOrderResp struct {
	Code            string `json:"code"` //工单号
	TrackCode       string `json:"track_code"`
	FinishedProduct string `json:"finished_product"` //产成品类型
	Origin          string `json:"origin"`
}

type WSStepReq struct {
	ID       int64  `json:"id"`
	StepCode string `json:"workstep_code"`
	Status   string `json:"status"`
}

type WSStepReqData struct {
	ID       int64  `json:"id"`
	StepCode string `json:"workstep_code"`
	Data     string `json:"data"`
}

type WSOrderReqCode struct {
	Code       string `json:"code"`
	Workcenter string `json:"workcenter"`
}

type WSWorkcenter struct {
	WorkCenter string `json:"workcenter"`
}

type WSLocalResults struct {
	HmiSN   string   `json:"hmi_sn" validate:"-"`
	Filters []string `json:"filters" validate:"required"`
	Limit   int      `json:"limit" validate:"-"`
}

type WSLocalResultsUpload struct {
	Id int `json:"id" validate:"gt=0"`
}

type WsScannerConflictResp struct {
	IsConflict bool `json:"isConflict"`
}

type WSOrderStartReq = WSOrderResumeReq

type WSOrderFinishReq = WSOrderPendingReq

type WSOrderPendingReq struct {
	Type    string    `json:"except_type" validate:"required"`
	Code    string    `json:"except_code" validate:"required"`
	Desc    string    `json:"except_desc" validate:"required"`
	EndTime time.Time `json:"end_time" validate:"required"`
	WSOrderCommonReq
}

type WSOrderCommonReq struct {
	OrderName      string `json:"order_name" validate:"required"`
	WorkCenterCode string `json:"workcenter_code" validate:"required"`
}

type WSOrderStartSimulateReq struct {
	WSOrderCommonReq
	Employee []string `json:"employee" validate:"required"`
}

type WSOrderResumeReq struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	WSOrderCommonReq
}

type WSNotify struct {
	Type    string      `json:"type" validate:"required,oneof=Info Warn Error"`
	Content string      `json:"content"`
	Config  interface{} `json:"config"`
}

func (s *WSNotify) Validate() error {
	if s.Type != HMINotifyInfo && s.Type != HMINotifyWarning && s.Type != HMINotifyError {
		return errors.Errorf("Type Error: %s", s.Type)
	}

	return nil
}

type NextWorkorder struct {
	Vin     string `json:"vin"`
	Model   string `json:"model"`
	LongPin string `json:"long_pin"`
	Knr     string `json:"knr"`
	Lnr     string `json:"lnr"`
}

type Result struct {
	ID            int64   `json:"id"`
	Controller_SN string  `json:"controller_sn"`
	GunSN         string  `json:"gun_sn"`
	PSet          int     `json:"pset"`
	MaxRedoTimes  int     `json:"max_redo_times"`
	X             float64 `json:"offset_x"`
	Y             float64 `json:"offset_y"`
	Seq           int     `json:"sequence"`
	GroupSeq      int     `json:"group_sequence"`
	NutNo         string  `json:"nut_no"`
}

type HMIWorkorderResp struct {
	Workorder_id   int64  `json:"workorder_id"`
	HMI_sn         string `json:"hmi_sn"`
	Vin            string `json:"vin"`
	Knr            string `json:"knr"`
	LongPin        string `json:"long_pin"`
	Status         string `json:"status"`
	WorkSheet      string `json:"work_sheet"`
	VehicleTypeImg string `json:"vehicleTypeImg"`

	MaxOpTime int      `json:"max_op_time"`
	Job       int      `json:"job_id"`
	Lnr       string   `json:"lnr"`
	Model     string   `json:"model"`
	Points    []Result `json:"points"`

	Reasons []string `json:"reasons"`
}

type ToolEnable struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	Enable        bool   `json:"enable"`
	Reason        string `json:"reason"`
}

type PSet struct {
	ControllerSN string `json:"controller_sn"`
	GunSN        string `json:"gun_sn" validate:"required"`
	PSet         int    `json:"pset" validate:"required"`
	Result_id    int64  `json:"result_id" validate:"-"`
	Count        int    `json:"count" validate:"gte=0"` //必须大于等于0
	UserID       int64  `json:"user_id" validate:"-"`
	GroupSeq     int    `json:"group_sequence" validate:"gt=0"`
	WorkorderID  int64  `json:"workorder_id" validate:"gt=0"` //必须大于0
}
