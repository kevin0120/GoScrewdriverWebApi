package tightening_device

import (
	"fmt"
	"github.com/masami10/rush/services/wsnotify"
	"gopkg.in/resty.v1"
	"testing"
	"time"
)

func getHttpClient() *resty.Client {
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(3 * time.Second)
	client.SetContentLength(true)
	// Headers for all request
	client.
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	return client
}

func wsTest(msg *wsnotify.WSMsg) {
	http := getHttpClient()
	resp, err := http.R().SetBody(msg).Put("http://127.0.0.1:8082/rush/v1/ws-test")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp.StatusCode())
}

func enableTool(controllerSN string, toolSN string, enable bool) {
	wsTest(&wsnotify.WSMsg{
		SeqNumber: 0,
		Type:      wsnotify.WS_TOOL_ENABLE,
		Data: ToolControl{
			ControllerSN: controllerSN,
			ToolSN:       toolSN,
			Enable:       enable,
		},
	})
}

func setPSet(controllerSN string, toolSN string, pset int) {
	wsTest(&wsnotify.WSMsg{
		SeqNumber: 0,
		Type:      wsnotify.WS_TOOL_PSET,
		Data: PSetSet{
			ToolControl: ToolControl{
				ControllerSN: controllerSN,
				ToolSN:       toolSN,
			},
			PSet: pset,
		},
	})
}

func Test_SetPSet(t *testing.T) {
	controllerSN := "w1"
	toolSN := "xx4443"
	pset := 4

	enableTool(controllerSN, toolSN, true)
	setPSet(controllerSN, toolSN, pset)
}
