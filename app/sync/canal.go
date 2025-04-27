package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/olivere/elastic/v7"
	"github.com/withlin/canal-go/client"
	protocol "github.com/withlin/canal-go/protocol/entry"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// 配置信息
type Config struct {
	CanalAddr      string
	CanalPort      int
	CanalUser      string
	CanalPassword  string
	CanalDest      string
	CanalBatchSize int
	ESAddrs        []string
	SyncTables     map[string]string // key为MySQL表名，value为ES索引名
	MySQLDatabase  string
}

func main() {
	// 配置信息
	cfg := Config{
		CanalAddr:      "canal-server",
		CanalPort:      11111,
		CanalUser:      "canal",
		CanalPassword:  "canal",
		CanalDest:      "example",
		CanalBatchSize: 100,
		ESAddrs:        []string{"http://es.xyq777.com"},
		MySQLDatabase:  "tag",
		SyncTables: map[string]string{
			"item": "item", // 将MySQL的users表同步到ES的users索引
		},
	}

	// 连接Elasticsearch
	esClient, err := elastic.NewClient(
		elastic.SetURL(cfg.ESAddrs...),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("连接Elasticsearch失败: %v", err)
	}

	// 创建Canal连接
	connector := client.NewSimpleCanalConnector(
		"canal-server", 11111, "", "", "example", 60000, 60*60*1000)
	err = connector.Connect()
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
				esIndex, ok := cfg.SyncTables[tableName]
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
						handleInsert(esClient, esIndex, rowData)
					case protocol.EventType_UPDATE:
						handleUpdate(esClient, esIndex, rowData)
					case protocol.EventType_DELETE:
						handleDelete(esClient, esIndex, rowData)
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

// 处理插入事件
func handleInsert(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
	doc, id := convertToMap(rowData.GetAfterColumns())
	if id == "" {
		log.Printf("插入操作缺少ID字段，跳过: %v", doc)
		return
	}

	// 执行ES插入
	_, err := esClient.Index().
		Index(esIndex).
		Id(id).
		BodyJson(doc).
		Refresh("false").
		Do(context.Background())

	if err != nil {
		log.Printf("ES插入数据失败: %v, 文档: %v", err, doc)
	} else {
		log.Printf("插入数据到ES, 索引: %s, ID: %s", esIndex, id)
	}
}

// 处理更新事件
func handleUpdate(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
	doc, id := convertToMap(rowData.GetAfterColumns())
	if id == "" {
		log.Printf("更新操作缺少ID字段，跳过: %v", doc)
		return
	}

	// 执行ES更新
	_, err := esClient.Update().
		Index(esIndex).
		Id(id).
		Doc(doc).
		DocAsUpsert(true).
		Refresh("false").
		Do(context.Background())

	if err != nil {
		log.Printf("ES更新数据失败: %v, 文档: %v", err, doc)
	} else {
		log.Printf("更新ES中的数据, 索引: %s, ID: %s", esIndex, id)
	}
}

// 处理删除事件
func handleDelete(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
	doc, id := convertToMap(rowData.GetBeforeColumns())
	if id == "" {
		log.Printf("删除操作缺少ID字段，跳过: %v", doc)
		return
	}

	// 执行ES删除
	_, err := esClient.Delete().
		Index(esIndex).
		Id(id).
		Refresh("false").
		Do(context.Background())

	if err != nil {
		// 如果是文档不存在的错误，可以忽略
		if strings.Contains(err.Error(), "404") {
			log.Printf("要删除的文档不存在，索引: %s, ID: %s", esIndex, id)
		} else {
			log.Printf("ES删除数据失败: %v, 文档ID: %s", err, id)
		}
	} else {
		log.Printf("从ES删除数据, 索引: %s, ID: %s", esIndex, id)
	}
}

// 将Column列表转换为map并返回ID
func convertToMap(columns []*protocol.Column) (map[string]interface{}, string) {
	doc := make(map[string]interface{})
	var id string

	for _, col := range columns {
		value := col.GetValue()

		// 根据数据类型进行转换
		switch col.GetMysqlType() {
		case "int", "tinyint", "smallint", "mediumint", "bigint":
			// 保留原始字符串，ES会自动转换
			doc[col.GetName()] = value
		case "float", "double", "decimal":
			doc[col.GetName()] = value
		case "date", "datetime", "timestamp", "time", "year":
			doc[col.GetName()] = value
		default: // string, text, blob等
			doc[col.GetName()] = value
		}

		// 记录ID字段
		if strings.ToLower(col.GetName()) == "id" {
			id = value
		}
	}

	return doc, id
}

func printEntry(entries []protocol.Entry) {
	for i := range entries {
		// 忽略事务开启和事务关闭类型
		if entries[i].GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN ||
			entries[i].GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		// RowChange对象，包含了一行数据变化的所有特征
		rowChange := new(protocol.RowChange)
		// protobuf解析
		err := proto.Unmarshal(entries[i].GetStoreValue(), rowChange)
		if err != nil {
			fmt.Printf("proto.Unmarshal failed, err:%v\n", err)
		}
		if rowChange == nil {
			continue
		}
		// 获取并打印Header信息
		header := entries[i].GetHeader()
		fmt.Printf("binlog[%s : %d], name[%s,%s], eventType: %s\n",
			header.GetLogfileName(),
			header.GetLogfileOffset(),
			header.GetSchemaName(),
			header.GetTableName(),
			header.GetEventType(),
		)
		//判断是否为DDL语句
		if rowChange.GetIsDdl() {
			fmt.Printf("isDdl:true, sql:%v\n", rowChange.GetSql())
		}

		// 获取操作类型：insert/update/delete等
		eventType := rowChange.GetEventType()
		for _, rowData := range rowChange.GetRowDatas() {
			if eventType == protocol.EventType_DELETE {
				printColumn(rowData.GetBeforeColumns())
			} else if eventType == protocol.EventType_INSERT || eventType == protocol.EventType_UPDATE {
				printColumn(rowData.GetAfterColumns())
			} else {
				fmt.Println("---before---")
				printColumn(rowData.GetBeforeColumns())
				fmt.Println("---after---")
				printColumn(rowData.GetAfterColumns())
			}
		}
	}
}

func printColumn(columns []*protocol.Column) {
	for _, col := range columns {
		fmt.Printf("%s:%s  updated=%v\n", col.GetName(), col.GetValue(), col.GetUpdated())
	}
}
