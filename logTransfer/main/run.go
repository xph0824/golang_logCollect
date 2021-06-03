package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"golang_logCollect/logTransfer/es"
	"time"
)

func run() (err error) {

	//kafka消费数据
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)

	if err != nil {
		logs.Error("获取kafka分区列表失败， err:", err)
		return
	}


	for partition := range partitionList {
		//  遍历所有的分区，并且针对每一个分区建立对应的消费者
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("分区启动消费者失败，partition: %d: %s\n", partition, err)
			return
		}

		defer pc.AsyncClose()

		go func(pc sarama.PartitionConsumer) {

			for msg := range pc.Messages() {
				logs.Debug("分区 Partition:%d, 偏移量 Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				err = es.SendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("发送数据到ES失败, err:%v", err)
				}
			}

		}(pc)
	}

	time.Sleep(time.Hour)
	return
}