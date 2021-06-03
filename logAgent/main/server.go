package main

import (
	"github.com/astaxie/beego/logs"
	"golang_logCollect/logAgent/kafka"
	"golang_logCollect/logAgent/tailf"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if nil != err {
			logs.Error("数据发送kafka失败, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	//fmt.Printf("读取 msg:%s, topic:%s\n", msg.Msg, msg.Topic)
	_ = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
