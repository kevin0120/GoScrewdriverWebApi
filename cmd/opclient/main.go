package main

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/config"
	"github.com/kevin0120/GoScrewdriverWebApi/services/diagnostic"
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/hmi"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol/vendors"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"os"
	"os/signal"
	"syscall"
)

// 监听中断信号，以便在按下Ctrl+C时保存配置并退出程序
func exit() {
	// 监听中断信号，以便在按下Ctrl+C时保存配置并退出程序
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		os.Exit(0)
	}()
}

func main() {
	// 获取命令行输入的参数
	// 检查是否至少有一个参数传入
	diagService := diagnostic.NewService(config.GetConfig().Logging, os.Stdout, os.Stderr)

	if err := diagService.Open(); err != nil {
		return
	}
	op := diagService.NewOpenProtocolHandler()
	tighteningSerivce, err := tightening_device.NewService(config.GetConfig().TighteningDevice, []tightening_device.ITighteningProtocol{
		openprotocol.NewService(config.GetConfig().OpenProtocol, op, vendors.OpenProtocolVendors),
	})
	if err != nil {
		return
	}
	err = tighteningSerivce.Open()
	if err != nil {
		return
	}
	httpDiag := diagService.NewHTTPDHandler()
	httpdService, err := httpd.NewService(config.GetConfig().DocPath, config.GetConfig().HTTP, config.GetConfig().Hostname, httpDiag, diagService)
	if err != nil {
		panic("!!!Panic: Can Not Open Http Service!!!")
	}
	err = httpdService.Open()
	if err != nil {
		return
	}
	hmiService := hmi.NewService(httpDiag, httpdService, "My", tighteningSerivce)
	err = hmiService.Open()
	if err != nil {
		return
	}
	fmt.Println("Op Client Running.")
	go exit()
	select {}
}
