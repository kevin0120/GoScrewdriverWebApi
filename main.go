package main

import (
	"github.com/kevin0120/GoScrewdriverWebApi/service/opserver"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	//开启op服务
	go opserver.StartOpServe()

	time.Sleep(time.Second)
	go exit()
	select {}
}
