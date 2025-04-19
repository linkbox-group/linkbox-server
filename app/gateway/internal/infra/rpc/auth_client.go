package rpc

import (
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/linkbox-group/linkbox-server/common/clientsuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth/authservice"
	"github.com/spf13/viper"
)

var (
	AuthClient authservice.Client
)

// 初始化Auth客户端
func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth", commonSuite)
	if err != nil {
		log.Fatalf("初始化Auth客户端失败: %s", err.Error())
	}
}

// 在InitClient函数中添加对initAuthClient的调用
func init() {
	// 确保InitClient函数中调用initAuthClient
	once.Do(func() {
		registryAddr = viper.GetString("consul.address")
		serviceName := viper.GetString("service.name")
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		})
		initUserClient()
		initAuthClient()
	})
}
