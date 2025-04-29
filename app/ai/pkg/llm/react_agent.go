package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/common/llm"
)

func GetReactAgentLambdaNode(ctx context.Context) *compose.Lambda {

	rAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: llm.GetOpenAIChatModel(),
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{GetContentInfoTool()},
		},
		MaxStep: 10,
	})
	if err != nil {
		log.Log().Errorf("failed to create agent: %v", err)
		return nil
	}
	agentLambda := compose.InvokableLambda(func(ctx context.Context, content string) (output *ItemTags, err error) {
		fmt.Println(content)
		p, err := prompt.FromMessages(schema.FString,
			schema.SystemMessage(contentPrompt)).Format(ctx, map[string]any{
			"blog":   content,
			"format": `{"tags": ["标签1", "标签2", "标签3"]}`,
		})
		if err != nil {
			log.Log().Errorf("failed to parse message: %v", err)
			return nil, err
		}
		ret, err := rAgent.Generate(ctx, p)
		if err != nil {
			log.Log().Errorf("failed to generate message: %v", err)
			return nil, err
		}
		fmt.Println(ret.Content)
		tags := &ItemTags{}
		err = json.Unmarshal([]byte(ret.Content), tags)
		if err != nil {
			log.Log().Errorf("failed to parse message: %v", err)
			tags.Tags = make([]string, 0)
			return tags, nil
		}
		return tags, nil

	})

	return agentLambda

}
