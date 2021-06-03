package es

import (
	"fmt"
	"context"
	"github.com/olivere/elastic/v7"
)

type Tweet struct {
	User	string
	Message string
}

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elastic.Client
)

func InitEs(addr string) (err error) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("连接ES失败：", err)
		return nil
	}
	esClient = client
	return
}

func SendToES(topic string, data []byte) (err error) {
	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	_, err = esClient.Index().
		Index(topic).
		BodyJson(msg).
		Do(context.Background())
	if err != nil {
		panic(err)
		return
	}
	return
}
