package vendors

import (
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol/vendors/desoutter"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol/vendors/lexen"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
)

var OpenProtocolVendors = map[string]openprotocol.IOpenProtocolController{
	tightening_device.ModelDesoutterCvi3: &desoutter.CVI3Controller{},
	tightening_device.ModelLexenWrench:   &lexen.WrenchController{},
	tightening_device.ModelLeetxTCS2000:  &lexen.TCS2000Controller{},
}
