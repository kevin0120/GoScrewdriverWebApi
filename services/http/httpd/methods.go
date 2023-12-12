package httpd

import (
	"github.com/kataras/iris/v12"
	"io/ioutil"
)

type Methods struct {
	service *Service
}

type HMICommonResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"msg"`
	Extra      string `json:"extra"`
}

func newHttpMethods(s *Service) Methods {
	return Methods{
		service: s,
	}
}

func (m *Methods) getDoc(ctx iris.Context) {
	f, _ := ioutil.ReadFile(m.service.ApiDoc)

	ctx.Header("content-type", "application/json")
	ctx.Write(f) //nolint
}

//
func NewCommonResponseBody(resp *HMICommonResponse, ctx iris.Context) error {
	ctx.StatusCode(resp.StatusCode)
	if err := ctx.JSON(resp); err != nil {
		return err
	}
	return nil
}
