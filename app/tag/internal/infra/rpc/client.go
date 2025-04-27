package rpc

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/item/itemservice"
	"log"

	"github.com/spf13/viper"

	"github.com/linkbox-group/linkbox-server/common/clientsuite"
	"sync"

	"github.com/cloudwego/kitex/client"
)

var (
	ItemClient   itemservice.Client
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
		initItemClient()

	})
}

func initItemClient() {
	ItemClient, err = itemservice.NewClient("item", commonSuite)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
