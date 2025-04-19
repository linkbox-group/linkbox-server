package main

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/user/userservice"
	"github.com/linkbox-group/linkbox-server/user/internal/core"
	"github.com/linkbox-group/linkbox-server/user/internal/infra/rpc"
	"log"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	rpc.InitClient()
	userHandler := NewUserHandler()
	srv := userservice.NewServer(userHandler)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
