package rpc

import (
	"log"

	"github.com/spf13/viper"

	"github.com/linkbox-group/linkbox-server/common/clientsuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth/authservice"

	"sync"

	"github.com/cloudwego/kitex/client"
)

var (
	AuthClient   authservice.Client
	once         sync.Once
	err          error
	registryAddr string
	commonSuite  client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = viper.GetString("consul.address")
		serviceName := viper.GetString("service.name")
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		})
		initAuthClient()

	})
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient("user", commonSuite)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
