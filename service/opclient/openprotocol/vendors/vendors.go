package vendors

import (
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/openprotocol"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/openprotocol/vendors/desoutter"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/openprotocol/vendors/lexen"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/tightening_device"
)

var OpenProtocolVendors = map[string]openprotocol.IOpenProtocolController{
	tightening_device.ModelDesoutterCvi3: &desoutter.CVI3Controller{},
	tightening_device.ModelLexenWrench:   &lexen.WrenchController{},
}
