package rpc

import (
	"github.com/linkbox-group/linkbox-server/common/clientsuite"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth/authservice"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item/itemservice"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization/organizationservice"
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag/tagservice"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user/userservice"
	"github.com/spf13/viper"
	"log"

	"sync"

	"github.com/cloudwego/kitex/client"
)

var (
	UserClient         userservice.Client
	AuthClient         authservice.Client
	TagClient          tagservice.Client
	ItemClient         itemservice.Client
	OrganizationClient organizationservice.Client
	once               sync.Once
	err                error
	registryAddr       string
	commonSuite        client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = viper.GetString("consul.address")
		serviceName := viper.GetString("service.name")
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		})
		initUserClient()
		initAuthClient()
		initTagClient()
		initOrganizationClient()
		initItemClient()
	})
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// 初始化Auth客户端
func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth", commonSuite)
	if err != nil {
		log.Fatalf("初始化Auth客户端失败: %s", err.Error())
	}
}
func initItemClient() {
	ItemClient, err = itemservice.NewClient("item", commonSuite)
	if err != nil {
		log.Fatalf(err.Error())
	}

}

// 初始化Tag客户端
func initTagClient() {
	TagClient, err = tagservice.NewClient("tag", commonSuite)
	if err != nil {
		log.Fatalf("初始化Tag客户端失败: %s", err.Error())
	}
}

func initOrganizationClient() {
	OrganizationClient, err = organizationservice.NewClient("organization", commonSuite)
	if err != nil {
		log.Fatalf("初始化Tag客户端失败: %s", err.Error())
	}

}
