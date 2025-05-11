package delivery

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/ai/internal/acl"
	"github.com/linkbox-group/linkbox-server/ai/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/ai/pkg/llm"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/ai/pkg/scrape"
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ProviderSet = wire.NewSet(NewAiDelivery)

type AiDelivery struct {
	service acl.AiServiceItf
}

func (d *AiDelivery) SuggestTags(ctx context.Context, req *ai.SuggestTagsRequest) (res *ai.SuggestTagsResponse, err error) {
	itemRes, err := rpc.ItemClient.GetItem(ctx, &item.GetItemRequest{
		Id:     req.ItemId,
		UserId: req.UserId,
	})
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}
	itemData := itemRes.GetItem()
	knowledge := itemData.Note
	if itemData.Type == itemmodel.ItemType_LINK {
		firecrawlApp := scrape.NewApp()
		knowledge, err = firecrawlApp.ExtractContent(itemData.Url)
		if err != nil {
			log.Log().Error(err.Error())
			return nil, err
		}
	}
	fmt.Println(knowledge)
	tags, err := llm.GenerateSuggestion(ctx, knowledge)
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}
	res = &ai.SuggestTagsResponse{
		Tags: tags,
	}
	return res, nil
}

func (d *AiDelivery) SendMessage(ctx context.Context, req *ai.SendMessageRequest, stream ai.AIService_SendMessageServer) (err error) {
	err = d.service.SendMessage(context.Background(), req, stream)
	if err != nil {
		log.Log().Error(err.Error())
	}
	return nil
}

func (d *AiDelivery) ListMessages(ctx context.Context, req *ai.ListMessagesRequest) (res *ai.ListMessagesResponse, err error) {
	// 调用usecase层获取消息列表
	messages, total, err := d.service.ListMessages(ctx, req)
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}
	pageInfo := req.GetPagination()

	// 将domain.Message转换为proto.Message
	var protoMessages []*ai.Message
	for _, msg := range messages {
		protoMessages = append(protoMessages, &ai.Message{
			Id:         msg.ID,
			UserId:     msg.UserID,
			Content:    msg.Content,
			SenderType: msg.SenderType,
			SendTime:   timestamppb.New(msg.SendTime),
		})
	}

	return &ai.ListMessagesResponse{
		Messages: protoMessages,
		Pagination: &pagination.PaginationMeta{
			TotalItems: int32(total),
			Page:       pageInfo.GetPage(),
			PageSize:   pageInfo.GetPageSize(),
			TotalPages: (int32(total) / pageInfo.PageSize) + 1,
		},
	}, nil
}

func (d *AiDelivery) DeleteMessage(ctx context.Context, req *ai.DeleteMessageRequest) (res *ai.DeleteMessageResponse, err error) {
	success, err := d.service.DeleteMessage(ctx, req.UserId, req.Ids)
	return &ai.DeleteMessageResponse{Success: success}, err
}

func NewAiDelivery(service acl.AiServiceItf) *AiDelivery {
	return &AiDelivery{
		service: service,
	}
}
