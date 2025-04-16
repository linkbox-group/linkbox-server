package main

import (
	"github.com/linkbox-group/linkbox-server/auth/internal/core"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth/authservice"
	"github.com/spf13/viper"
	"log"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(viper.GetString("REDIS_PASSWORD"))
	authHandler := NewAuthHandler()
	srv := authservice.NewServer(authHandler)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
