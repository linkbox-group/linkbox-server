package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/linkbox-group/linkbox-server/auth/internal/core"
	"github.com/linkbox-group/linkbox-server/common/serversuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth/authservice"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(viper.GetString("REDIS_PASSWORD"))
	authHandler := NewAuthHandler()
	srv := authservice.NewServer(authHandler, kitexInit()...)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
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
