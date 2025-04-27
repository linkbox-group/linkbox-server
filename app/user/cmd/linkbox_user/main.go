package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/linkbox-group/linkbox-server/common/serversuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user/userservice"
	"github.com/linkbox-group/linkbox-server/user/internal/core"
	"github.com/linkbox-group/linkbox-server/user/internal/infra/rpc"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {
	core.LoadLog()
	err := core.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	rpc.InitClient()
	userHandler := NewUserHandler()
	srv := userservice.NewServer(userHandler, kitexInit()...)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", viper.GetString("service.address"))
	if err != nil {
		panic(err)
	}

	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithSuite(
		serversuite.CommonServerSuite{
			CurrentServiceName: viper.GetString("service.name"),
			RegistryAddr:       viper.GetString("consul.address")}),
	)
	return
}
