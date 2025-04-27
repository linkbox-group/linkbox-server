package core

import (
	"github.com/olivere/elastic/v7"
	"log"
)

func LoadEs() *elastic.Client {
	// 连接Elasticsearch
	esClient, err := elastic.NewClient(
		elastic.SetURL(DefaultConfig.ESAddrs...),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("连接Elasticsearch失败: %v", err)
	}
	return esClient
}
