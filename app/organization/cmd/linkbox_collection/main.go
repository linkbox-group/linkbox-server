package main

import (
	"github.com/linkbox-group/linkbox-server/organization/internal/core"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization/organizationservice"
	"log"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal("load config failed", err)
	}
	organizationHandler := NewOrganizationHandler()
	srv := organizationservice.NewServer(organizationHandler)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Failed to run content service: %v", err)
	}
}
