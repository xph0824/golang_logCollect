package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"strings"
)

type KafkaClient struct {
	client  sarama.Consumer
	addr 	string
	topic 	string
}

var (
	kafkaClient *KafkaClient
)

//初始化kafka消费者
func InitKafka(addr string, topic string) (err error) {
	kafkaClient = &KafkaClient{}
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("启动kafka消费失败：", err)
		return nil
	}
	kafkaClient.client = consumer
	kafkaClient.addr   = addr
	kafkaClient.topic  = topic
	return
}
