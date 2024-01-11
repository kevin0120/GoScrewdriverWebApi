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
	cr := tightening_device.TighteningDeviceCmd{}
	err := ctx.ReadJSON(&cr)

	if err != nil {
		// 传输结构错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	m.TighteningService.HandleHmiRequest(cr)
}
