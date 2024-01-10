package hmi

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"github.com/kevin0120/GoScrewdriverWebApi/utils"
	"strconv"
)

var CONFLICT_VIA_LNR = false
var validate *validator.Validate

func init() {
	ENV_CONFLICT_VIA_LNR := utils.GetEnv("ENV_CONFLICT_VIA_LNR", "false") //是否车序控制conflict
	CONFLICT_VIA_LNR, _ = strconv.ParseBool(ENV_CONFLICT_VIA_LNR)
	validate = validator.New()
}

func (m *Service) tighteningControl(ctx iris.Context) {
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
