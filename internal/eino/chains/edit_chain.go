package chains

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// ContentEditChain 内容编辑链
type ContentEditChain struct {
	chain *compose.Chain[map[string]any, map[string]any]
}

// NewContentEditChain 创建内容编辑链
func NewContentEditChain(components *components.Components) (*ContentEditChain, error) {
	chain := compose.NewChain[map[string]any, map[string]any]()

	// Node 1: 格式化提示词
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
			template := components.ChatTemplate
			args := map[string]any{
				"Content": input["content"],
			}
			return template.Format(ctx, "edit_content", args)
		}),
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("ContentEditing"))

	// Node 3: 解析编辑结果
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (map[string]any, error) {
			return map[string]any{
				"original":  "",
				"edited":    msg.Content,
				"changes":   []string{},
				"improved":  true,
			}, nil
		}),
		compose.WithNodeName("ResultParse"),
	)

	return &ContentEditChain{chain: chain}, nil
}

// Edit 编辑内容
func (c *ContentEditChain) Edit(ctx context.Context, input map[string]any) (map[string]any, error) {
	return c.chain.Invoke(ctx, input)
}
