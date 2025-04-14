package main

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag/tagservice"
	"log"
)

func main() {
	tagHandler := NewTagHandler()
	srv := tagservice.NewServer(tagHandler)
	err := srv.Run()
	if err != nil {
		log.Fatalf("Failed to run content service: %v", err)
	}
}
