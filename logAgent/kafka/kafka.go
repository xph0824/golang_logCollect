package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

//连接
func InitKafka(addr string) (err error) {

	// 配置kafka生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll 			// 数据发送完毕需要leader和follow都确认
	config.Producer.Partitioner  = sarama.NewRandomPartitioner 	// 随机分配分区 partition
	// 是否开启消息发送成功后通知 successes channel
	config.Producer.Return.Successes = true						// 成功交付的消息将在success channel返回

	// 创建生产者
	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("初始化kafka生产者失败， err:", err)
		return
	}

	logs.Debug("初始化kafka生产者成功，地址：", addr)
	return
}

//发送消息
func SendToKafka(data, topic string) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		logs.Error("发送消息失败， err:%v, data:%v, topic:%v", err, data, topic)
		return
	}

	logs.Debug("发送消息成功， pid:%v, offset:%v, topic%v\n", pid, offset, topic)
	return
}