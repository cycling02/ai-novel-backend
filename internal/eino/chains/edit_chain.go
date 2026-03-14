package chains

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// ContentEditChain 内容编辑链
type ContentEditChain struct {
	chain   *compose.Chain[map[string]any, map[string]any]
	runnable compose.Runnable[map[string]any, map[string]any]
}

// NewContentEditChain 创建内容编辑链
func NewContentEditChain(components *components.Components) (*ContentEditChain, error) {
	chain := compose.NewChain[map[string]any, map[string]any]()

	// Node 1: ChatTemplate
	chain = chain.AppendChatTemplate(
		components.ChatTemplate,
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("ContentEditing"))

	// Node 3: 解析编辑结果
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input any) (map[string]any, error) {
			content := ""
			switch msg := input.(type) {
			case *schema.Message:
				content = msg.Content
			case schema.Message:
				content = msg.Content
			default:
				content = fmt.Sprintf("%v", input)
			}
			return map[string]any{
				"original":  "",
				"edited":    content,
				"changes":   []string{},
				"improved":  true,
			}, nil
		}),
		compose.WithNodeName("ResultParse"),
	)

	// 编译 Chain
	rCtx := context.Background()
	runnable, err := chain.Compile(rCtx)
	if err != nil {
		return nil, fmt.Errorf("编译 Chain 失败：%w", err)
	}

	return &ContentEditChain{
		chain:   chain,
		runnable: runnable,
	}, nil
}

// Edit 编辑内容
func (c *ContentEditChain) Edit(ctx context.Context, input map[string]any) (map[string]any, error) {
	// 准备模板变量
	args := map[string]any{
		"template_name": "edit_content",
		"Content":       input["content"],
	}

	return c.runnable.Invoke(ctx, args)
}