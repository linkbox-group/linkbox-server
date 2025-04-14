package main

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag/tagservice"
	"github.com/linkbox-group/linkbox-server/tag/internal/core"
	"log"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal("load config failed", err)
	}
	tagHandler := NewTagHandler()
	srv := tagservice.NewServer(tagHandler)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Failed to run content service: %v", err)
	}
}
