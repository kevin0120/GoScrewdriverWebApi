package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"github.com/kevin0120/GoScrewdriverWebApi/utils"
	"strconv"
	"strings"
)

var CONFLICT_VIA_LNR = false
var validate *validator.Validate

func init() {
	ENV_CONFLICT_VIA_LNR := utils.GetEnv("ENV_CONFLICT_VIA_LNR", "false") //是否车序控制conflict
	CONFLICT_VIA_LNR, _ = strconv.ParseBool(ENV_CONFLICT_VIA_LNR)
	validate = validator.New()
}

func (s *Service) putNotify(ctx iris.Context) {
	var req WSNotify
	resp := HMICommonResponse{StatusCode: iris.StatusBadRequest}

	if err := ctx.ReadJSON(&req); err != nil {
		resp.Message = "putNotify Read JSON Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func (s *Service) listWorkorders(ctx iris.Context) {
	var p listWorkordersURLParams
	resp := HMICommonResponse{StatusCode: iris.StatusBadRequest}

	if err := ctx.ReadQuery(&p); err != nil {
		resp.Message = "listWorkorders Validate Param Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	if err := validate.Struct(p); err != nil {
		resp.Message = "listWorkorders Validate Param Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

}

func (s *Service) getNextWorkorder(ctx iris.Context) {

	hmi_sn := ctx.URLParam("hmi_sn")
	workcenterCode := ctx.URLParam("workcenter_code")

	if hmi_sn == "" && workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn or workcenter_code is required")
		return
	}

}

// deprecated
func (s *Service) getLocalResults(ctx iris.Context) {

	var rt []LocalResults

	body, _ := json.Marshal(rt)
	ctx.Header("content-type", "application/json")
	ctx.Write(body)
}

func (s *Service) filterValue(filters string, key string, value interface{}) interface{} {
	if filters == "" || strings.Contains(filters, key) {
		return value
	}
	return nil
}

// 根据hmi序列号以及vin或knr取得工单
func (s *Service) getWorkorderDetail(ctx iris.Context) {
	//m.service.diag.Debug("getWorkorderDetail start")

	var err error
	hmi_sn := ctx.URLParam("hmi_sn")
	workcenterCode := ctx.URLParam("workcenter_code")
	code := ctx.URLParam("code")

	if hmi_sn == "" && workcenterCode == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("hmi_sn or workcenter_code is required")
		return
	}

	if code == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("code is required")
		return
	}

	err = nil
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("workOrder not found")
		return
	}

	//m.service.diag.Debug(fmt.Sprintf("getWorkorderDetail finish with body:%s", string(body)))
}

func (s *Service) putToolControl(ctx iris.Context) {
	var controllerSn string
	req := ToolEnable{}
	err := ctx.ReadJSON(&req)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	if req.Controller_SN == "" {
		_, err := s.TighteningService.GetControllerByToolSN(req.GunSN)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
		controllerSn = ""
	} else {
		controllerSn = req.Controller_SN
	}

	err = s.TighteningService.ToolControl(&tightening_device.ToolControl{
		ControllerSN: controllerSn,
		ToolSN:       req.GunSN,
		Enable:       req.Enable,
	})

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func (m *Service) postAK2(ctx iris.Context) {
	cr := tightening_device.TighteningResult{}
	err := ctx.ReadJSON(&cr)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	if cr.ControllerSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("controller_sn is required")
		return
	}

	if cr.ToolSN == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("gun_sn is required")
		return
	}

	//if cr.PSet == 0 {
	//	ctx.StatusCode(iris.StatusBadRequest)
	//	ctx.WriteString("pset is required")
	//	return
	//}

	if cr.Count == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("count is required")
		return
	}

	if cr.Seq == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("seq is required")
		return
	}
	if cr.WorkorderID == 0 {
		//没有工单直接点击了放行按钮，猜测是为了联动
		//fixme: 现在使用控制器名称作为站点的信息
		return
	}
}

func (s *Service) putPSets(ctx iris.Context) {

	var err error
	var pset PSet
	var controllerSn string
	err = ctx.ReadJSON(&pset)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	s.diag.Debug(fmt.Sprintf("new pset: %#v", pset))

	if err := validate.Struct(pset); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("Validate Error: %s", err.Error())
		return
	}

	err = s.TighteningService.ToolPSetSet(&tightening_device.PSetSet{
		ToolControl: tightening_device.ToolControl{
			ControllerSN: controllerSn,
			ToolSN:       pset.GunSN,
		},
		WorkorderID:   pset.WorkorderID,
		WorkorderCode: "",
		UserID:        pset.UserID,
		PSet:          pset.PSet,
		Sequence:      uint(pset.GroupSeq),
		Count:         pset.Count,
		Total:         1,
		Batch:         pset.Count,
	})

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
}
