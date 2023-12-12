package main

import (
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opserver"
	"github.com/kevin0120/GoScrewdriverWebApi/services/udp/udpclient"
	"os"
	"os/signal"
	"strconv"
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
	port := 4545
	if len(os.Args) == 2 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	client := udpclient.NewClient(3000)

	go func() {
		go client.ConnectToServer("211.254.254.250", 8080, 0)
		go func() {
			err := client.ReadMultiSdoCircle([]string{"0x300803", "0x300811", "0x300814", "0x100006", "0x100007", "0x30010A", "0x300807", "0x300808", "0x300831"})
			if err != nil {
				return
			}
		}()
	}()
	//开启op服务
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	go opserver.StartOpServe(addr, client)
	fmt.Println("Op Serve Running.")
	go exit()
	select {}
}
