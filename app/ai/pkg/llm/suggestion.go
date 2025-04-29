package llm

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
)

const (
	ReactLambda = "react_lambda"
)

type ItemTags struct {
	Tags []string `json:"tags"`
}

func GenerateSuggestion(ctx context.Context, content string) (tags []string, err error) {

	workflow := compose.NewGraph[string, *ItemTags]()
	_ = workflow.AddLambdaNode(ReactLambda, GetReactAgentLambdaNode(ctx))
	//edge
	_ = workflow.AddEdge(compose.START, ReactLambda)
	_ = workflow.AddEdge(ReactLambda, compose.END)

	r, err := workflow.Compile(ctx, compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}

	ret, err := r.Invoke(ctx, content)
	if err != nil {
		log.Log().Error(err.Error())
		return nil, err
	}

	return ret.Tags, nil
}
