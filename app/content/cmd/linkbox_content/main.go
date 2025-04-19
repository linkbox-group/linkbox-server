package main

import (
	"github.com/linkbox-group/linkbox-server/content/internal/core"
	"github.com/linkbox-group/linkbox-server/rpc-gen/content/contentservice"
	"log"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	contentHandler := NewContentHandler()
	srv := contentservice.NewServer(contentHandler)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Failed to run content service: %v", err)
	}
}
