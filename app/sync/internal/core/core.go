package core

import (
	"fmt"
	"github.com/withlin/canal-go/client"
)

func LoadCanal() *client.SimpleCanalConnector {
	// 创建Canal连接
	connector := client.NewSimpleCanalConnector(
		DefaultConfig.CanalAddr,
		DefaultConfig.CanalPort,
		"",
		"",
		DefaultConfig.CanalDest,
		60000,
		60*60*1000,
	)
	err := connector.Connect()
	if err != nil {
		panic(err)
	}

	// 订阅数据库表变更，格式为：数据库.表，可以用正则表达式
	// 例如：.*\\..*表示所有库所有表，test\\.test表示test库的test表
	err = connector.Subscribe(".*\\..*")
	if err != nil {
		fmt.Printf("connector.Subscribe failed, err:%v\n", err)
		panic(err)
	}
	return connector
}
