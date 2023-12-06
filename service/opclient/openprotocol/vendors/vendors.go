package vendors

import (
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/openprotocol/vendors/crane"
	"github.com/masami10/rush/services/openprotocol/vendors/desoutter"
	"github.com/masami10/rush/services/openprotocol/vendors/lexen"
	"github.com/masami10/rush/services/tightening_device"
)

var OpenProtocolVendors = map[string]openprotocol.IOpenProtocolController{
	tightening_device.ModelDesoutterCvi3:        &desoutter.CVI3Controller{},
	tightening_device.ModelDesoutterCvi2:        &desoutter.CVI2Controller{},
	tightening_device.ModelDesoutterDeltaWrench: &desoutter.WrenchController{},
	tightening_device.ModelDesoutterConnector:   &desoutter.ConnectorController{},
	tightening_device.ModelCraneIQWrench:        &crane.IQWrenchController{},
	tightening_device.ModelDesoutterCvi3Twin:    &desoutter.CVI3TwinController{},
	tightening_device.ModelLexenWrench:          &lexen.WrenchController{},
}
