package llm

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/common/llm"
)

func getChatPromptTemplate() prompt.ChatTemplate {
	return prompt.FromMessages(schema.FString,
		schema.SystemMessage(chatPrompt),
		schema.MessagesPlaceholder("chat_histories", false),
		schema.UserMessage("{question}"),
	)

}

func Chat(ctx context.Context, question string, item string) (*schema.StreamReader[*schema.Message], error) {

	workflow := compose.NewChain[map[string]any, *schema.Message]()
	workflow.
		AppendChatTemplate(getChatPromptTemplate()).
		AppendChatModel(llm.GetOpenAIChatModel())

	r, err := workflow.Compile(ctx)
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}

	in := map[string]any{
		"question": question,
		//"item":     item,
		"chat_histories": []*schema.Message{
			{},
		},
	}
	ret, err := r.Stream(ctx, in)
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}

	return ret, nil
}
