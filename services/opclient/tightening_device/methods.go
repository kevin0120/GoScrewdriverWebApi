package tightening_device

import (
	"github.com/pkg/errors"
	"strings"
)

var EnvPsetBatchEnable = true
var EnvPsetWithEnable = true

type JobSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	WorkorderID  int64  `json:"workorder_id"`
	Total        int    `json:"total"`
	StepID       int64  `json:"step_id"`
	Job          int    `json:"job"`
	UserID       int64  `json:"user_id"`
}

func (s *JobSet) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
	}

	if s.Job <= 0 {
		return errors.New("Job Should Be Greater Than Zero")
	}

	return nil
}

type PSetBatchSet struct {
	ControllerSN string `json:"controller_sn" validate:"required"`
	ToolSN       string `json:"tool_sn" validate:"required"`
	PSet         int    `json:"pset" validate:"required"`
	Batch        int    `json:"batch" validate:"required"`
}

//todo: 接口中哪些字段是的验证条件还不确认
type PSetSet struct {
	ToolControl
	StepCode      string `json:"workstep_code" validate:"-"`
	WorkorderID   int64  `json:"workorder_id" validate:"-"`
	WorkorderCode string `json:"workorder_code" validate:"-"`
	UserID        int64  `json:"user_id" validate:"gt=0"`
	PSet          int    `json:"pset" validate:"gt=0"`
	Sequence      uint   `json:"sequence" validate:"gt=0"`
	Count         int    `json:"count" validate:"-"` //拧紧结果计数，发送请求时应该不传递
	Batch         int    `json:"batch" validate:"-"`
	Total         int    `json:"total" validate:"-"`
	IP            string `json:"ip" validate:"-"`
	PointID       string `json:"point_id" validate:"-"`
	ScannerCode   string `json:"scanner_code" validate:"-"`
}

type ToolControl struct {
	ControllerSN string `json:"controller_sn" validate:"required"`
	ToolSN       string `json:"tool_sn" validate:"required"`
	Enable       bool   `json:"enable" validate:"required"`
}

type ToolModeSelect struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Mode         string `json:"mode"`
}

func (s *ToolModeSelect) Validate() error {
	if s.ControllerSN == "" || s.ToolSN == "" {
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
	}

	return nil
}

func (s *Service) ToolControl(req *ToolControl) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}
	return nil
	//return tool.ToolControl(req.Enable)
}

func (s *Service) ToolJobSet(req *JobSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	_, err = s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}
	return nil
	//return tool.SetJob(req.Job)
}

func (s *Service) ToolPSetBatchSet(req *PSetBatchSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	if req.PSet == 0 {
		return errors.New("ToolPSetBatchSet Pset Must Be Greater Than Zero")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}
	return nil
	//return tool.SetPSetBatch(req.PSet, req.Batch)
}

func (s *Service) doTraceFromPSetReq(req *PSetSet) {

	if req.UserID == 0 {
		req.UserID = 1
	}

}

func (s *Service) ToolPSetByIP(req *PSetSet) error {
	if req == nil {
		return errors.New("ToolPSetByIP Req Is Nil")
	}

	_, err := s.findToolbyIP(req.IP)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	//controller := tool.GetParentService().(ITighteningController)
	if err = s.ToolPSetBatchSet(&PSetBatchSet{
		//ControllerSN: controller.SerialNumber(),
		//ToolSN:       tool.SerialNumber(),
		PSet:  req.PSet,
		Batch: 1,
	}); err != nil {
		//s.diag.Error("PSet Batch Set Failed", err)
	}

	return nil
}

func (s *Service) ToolPSetSet(req *PSetSet) error {

	dopset := func(tool ITighteningTool) error {
		//ctx := context.WithValue(context.Background(), "psetReq", req)
		return nil
		//return tool.SetPSet(ctx, req.PSet)
	}

	if req == nil {
		return errors.New("Req Is Nil")
	}

	s.doTraceFromPSetReq(req)

	controller, err := s.getController(req.ControllerSN)
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	if req.PSet <= 0 {
		return errors.New("ToolPSetSet.ToolPSetBatchSet Pset Must Be Greater Than 0!!!")
	}

	if EnvPsetBatchEnable {
		err := s.ToolPSetBatchSet(&PSetBatchSet{
			ControllerSN: req.ControllerSN,
			ToolSN:       req.ToolSN,
			PSet:         req.PSet,
			Batch:        req.Batch,
		})
		if err != nil {
			if !strings.Contains(err.Error(), TIGHTENING_ERR_NOT_SUPPORTED) {
				return err
			} else {
				if err := dopset(tool); err != nil {
					return err
				}
			}
		}
		if controller.Model() != ModelLexenWrench {
			if err := dopset(tool); err != nil {
				return err
			}
		}
	} else {
		if err := dopset(tool); err != nil {
			return err
		}
	}

	if req.UserID <= 0 {
		req.UserID = 1
	}

	//FIXME:手动模式下sequence可能为0
	// if req.Sequence <= 0 {
	// 	err := errors.New("Sequence Is Less Than Zero")
	// 	s.diag.Error("ToolPSetSet", err)
	// 	return err
	// }

	if req.Enable || EnvPsetWithEnable {
		//_ = tool.ToolControl(true)
	}

	//s.diag.Info(fmt.Sprintf("Pset Request Pset Number: %d Success!!!", req.PSet))

	return nil
}

func (s *Service) ToolModeSelect(req *ToolModeSelect) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	_, err = s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return err
	}

	return nil
	//return tool.ModeSelect(req.Mode)
}

type ToolInfo struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
}

type ToolPSet struct {
	ToolInfo
	PSet int `json:"pset"`
}

type ToolJob struct {
	ToolInfo
	Job int `json:"job"`
}

func (s *Service) GetToolPSetList(req *ToolInfo) ([]PSetInfo, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}
	return nil, nil
	//return tool.GetPSetList()
}

func (s *Service) GetToolPSetDetail(req *ToolPSet) (*PSetDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}
	return nil, nil
	//return tool.GetPSetDetail(req.PSet)
}

func (s *Service) GetToolJobList(req *ToolInfo) ([]int, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return nil, nil
	//return tool.GetJobList()
}

func (s *Service) GetToolJobDetail(req *ToolJob) (*JobDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	_, err := s.getTool(req.ControllerSN, req.ToolSN)
	if err != nil {
		return nil, err
	}

	return nil, nil
	//return tool.GetJobDetail(req.Job)
}

func (s *Service) findToolbyIP(ip string) (ITighteningTool, error) {
	//for _, controller := range s.runningControllers {
	//	tool, err := controller.GetToolViaIP(ip)
	//	if err == nil {
	//		return tool, nil
	//	}
	//}

	return nil, errors.New("findToolbyIP: Not Found")
}
