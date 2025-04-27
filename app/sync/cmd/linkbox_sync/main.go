package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/linkbox-group/linkbox-server/sync/internal/core"
	"github.com/linkbox-group/linkbox-server/sync/internal/service"
	protocol "github.com/withlin/canal-go/protocol/entry"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var SyncTables = map[string]string{
	"item": "item",
}

// 配置信息

func main() {
	core.LoadConfig()
	esClient := core.LoadEs()
	connector := core.LoadCanal()

	// 定义关闭函数
	closeFunc := func() {
		connector.DisConnection()
		log.Println("关闭Canal连接...")
	}

	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-signalChan
		closeFunc()
		os.Exit(0)
	}()

	defer closeFunc()

	log.Println("启动Canal MySQL到ES同步服务...")
	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Printf("获取Canal消息失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(time.Second)
			continue
		}
		//处理消息
		for i := range message.Entries {

			// 跳过事务开始和结束的消息
			if message.Entries[i].GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || message.Entries[i].GetEntryType() == protocol.EntryType_TRANSACTIONEND {
				continue
			}

			// 只处理行级变更消息
			if message.Entries[i].GetEntryType() == protocol.EntryType_ROWDATA {
				// 获取表信息
				tableName := message.Entries[i].GetHeader().GetTableName()
				_ = message.Entries[i].GetHeader().GetSchemaName()

				// 检查是否为需要同步的表
				esIndex, ok := SyncTables[tableName]
				if !ok {
					continue // 不是需要同步的表，跳过
				}
				log.Printf("开始处理")

				// 解析行变更事件
				rowChange := new(protocol.RowChange)
				err := proto.Unmarshal(message.Entries[i].GetStoreValue(), rowChange)
				if err != nil {
					log.Printf("解析Canal消息失败: %v", err)
					continue
				}

				eventType := rowChange.GetEventType()

				// 处理不同类型的事件
				for _, rowData := range rowChange.GetRowDatas() {
					switch eventType {
					case protocol.EventType_INSERT:
						service.HandleInsert(esClient, esIndex, rowData)
					case protocol.EventType_UPDATE:
						service.HandleUpdate(esClient, esIndex, rowData)
					case protocol.EventType_DELETE:
						service.HandleDelete(esClient, esIndex, rowData)
					default:
						log.Printf("未处理的事件类型: %v", eventType)
					}
				}
			}
		}

		// 提交消息确认
		//connector.Ack(batchId)
	}
}
