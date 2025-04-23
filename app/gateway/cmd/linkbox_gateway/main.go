package main

import (
	"github.com/linkbox-group/linkbox-server/gateway/internal/core"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/gateway/internal/router"
	"github.com/spf13/viper"
	"log"
)

func main() {
	core.LoadLog()
	err := core.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	r := router.InitRouter()
	rpc.InitClient()
	err = r.Run(viper.GetString("service.address"))
	if err != nil {
		log.Fatalf(err.Error())
	}
}
