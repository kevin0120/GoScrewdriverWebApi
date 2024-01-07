package udpclient

import (
	"github.com/kataras/iris/v12"
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
)

func (c *UdpClient) HandleHttpRequest() {
	var r httpd.Route
	//增加显示屏用户切换接口
	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/hmi-user",
		HandlerFunc: c.getDoc,
	}
	_ = httpd.Httpd.Handler[0].AddRoute(r)

	return
}

func (c *UdpClient) getDoc(ctx iris.Context) {

	ctx.Header("content-type", "application/json")
	_, _ = ctx.Write([]byte("f"))
}
