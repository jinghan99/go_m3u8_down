package main

import (
	"fmt"
	"go_m3u8_down/conf"
	"go_m3u8_down/routers"
	"io"
	"log"
	"os"
)

func init() {
	initLog()
}

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	//	启动gin 服务器
	routers.HttpServer(conf.Port)
	defer conf.Ndb.Close()
}

//初始化 log
func initLog() {
	// 追加日志 可添加 O_APPEND
	writer1, logErr := os.OpenFile("go_m3u8_down.log", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if logErr != nil {
		fmt.Println("创建日志失败")
		os.Exit(1)
	}
	//os.Stdout代表标准输出流
	writer2 := os.Stdout
	// //io.MultiWriter实现多目的地输出 组合一下即可，
	multiWriter := io.MultiWriter(writer1, writer2)
	//	设置 logger 输出日志文件
	log.SetOutput(multiWriter)
	log.SetPrefix("[go_m3u8_down]")
	//	设置 logger 输出时间格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
