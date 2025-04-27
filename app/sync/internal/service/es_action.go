package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/olivere/elastic/v7"
	protocol "github.com/withlin/canal-go/protocol/entry"
	"log"
	"strings"
)

// 处理插入事件
func HandleInsert(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
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
func HandleUpdate(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
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
func HandleDelete(esClient *elastic.Client, esIndex string, rowData *protocol.RowData) {
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
