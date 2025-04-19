package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/linkbox-group/linkbox-server/common/serversuite"
	"github.com/linkbox-group/linkbox-server/organization/internal/core"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization/organizationservice"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {
	err := core.LoadConfig()
	if err != nil {
		log.Fatal("load config failed", err)
	}
	organizationHandler := NewOrganizationHandler()
	srv := organizationservice.NewServer(organizationHandler, kitexInit()...)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Failed to run content service: %v", err)
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
