package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/linkbox-group/linkbox-server/ai/internal/core"
	"github.com/linkbox-group/linkbox-server/ai/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/common/serversuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai/aiservice"
	"github.com/spf13/viper"
	"net"
)

func main() {

	core.LoadLog()
	rpc.InitClient()

	err := core.LoadConfig()

	if err != nil {
		log.Log().Fatalf("load config failed: %v", err)
	}
	AiHandler := NewAiHandler()
	srv := aiservice.NewServer(AiHandler, kitexInit()...)
	err = srv.Run()
	if err != nil {
		log.Log().Fatalf("Failed to run content service: %v", err)
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
