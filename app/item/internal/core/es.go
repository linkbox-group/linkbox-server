package core

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
)

func NewEs() *elasticsearch.TypedClient {
	cfg := elasticsearch.Config{
		//Username: "elastic",
		Addresses: []string{
			viper.GetString("es.address"),
		},
	}
	esClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return esClient
}
