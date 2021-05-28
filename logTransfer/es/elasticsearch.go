package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Tweet struct {
	User	string
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
