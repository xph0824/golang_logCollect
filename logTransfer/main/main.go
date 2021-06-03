package main

import (
	"github.com/astaxie/beego/logs"
	"golang_logCollect/logTransfer/es"
)

func main() {

	//初始化配置项
	err := InitConfig("ini", "/Users/syl/golang_logCollect/logTransfer/config/logTransfer.conf")
	if err != nil {
		panic(err)
		return
	}
	logs.Debug("初始化配置成功") //日志模块还没有初始化，所以这行是写不写进去的

	//初始化日志模块
	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}
	logs.Debug("初始化日志成功")

	//初始化kafka
	err = InitKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("初始化kafka失败， err:", err)
		return
	}
	logs.Debug("初始化kafka成功")

	//初始化Es
	err = es.InitEs(logConfig.EsAddr)
	if err != nil {
		logs.Error("初始化ES失败， err:", err)
		return
	}
	logs.Debug("初始化ES成功")

	//run
	err = run()
	if err != nil {
		logs.Error("运行错误：", err)
		return
	}

	logs.Warn("logTransfer 退出")
}