package main

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/user/userservice"
	"log"
)

func main() {
	userHandler := NewUserHandler()
	srv := userservice.NewServer(userHandler)
	err := srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
