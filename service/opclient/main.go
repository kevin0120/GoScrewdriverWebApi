package main

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/service/opclient/tightening_device"
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
	service, err := tightening_device.NewService()
	if err != nil {
		return
	}
	fmt.Println("Op Serve Running.")
	go exit()
	select {}
}
