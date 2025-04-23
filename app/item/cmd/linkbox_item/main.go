package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/linkbox-group/linkbox-server/common/serversuite"
	"github.com/linkbox-group/linkbox-server/item/internal/core"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item/itemservice"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {
	core.LoadLog()
	err := core.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	itemHandler := NewItemHandler()
	srv := itemservice.NewServer(itemHandler, kitexInit()...)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Failed to run item service: %v", err)
	}
}
func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", viper.GetString("service.address"))
	if err != nil {
		log.Fatal(err)
	}

	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithSuite(
		serversuite.CommonServerSuite{
			CurrentServiceName: viper.GetString("service.name"),
			RegistryAddr:       viper.GetString("consul.address")}),
	)
	return
}
