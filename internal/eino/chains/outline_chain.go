package chains

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// OutlineExpandChain 大纲扩写链
type OutlineExpandChain struct {
	chain *compose.Chain[map[string]any, string]
}

// NewOutlineExpandChain 创建大纲扩写链
func NewOutlineExpandChain(components *components.Components) (*OutlineExpandChain, error) {
	chain := compose.NewChain[map[string]any, string]()

	// Node 1: 格式化提示词
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
			template := components.ChatTemplate
			args := map[string]any{
				"NovelTitle":   input["novel_title"],
				"Genre":        input["genre"],
				"BriefOutline": input["brief_outline"],
			}
			return template.Format(ctx, "expand_outline", args)
		}),
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("OutlineExpansion"))

	// Node 3: 解析输出
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (string, error) {
			content := msg.Content
			// 解析结构化大纲
			return parseOutline(content), nil
		}),
		compose.WithNodeName("OutputParse"),
	)

	return &OutlineExpandChain{chain: chain}, nil
}

// Expand 扩写大纲
func (c *OutlineExpandChain) Expand(ctx context.Context, input map[string]any) (string, error) {
	return c.chain.Invoke(ctx, input)
}

func parseOutline(content string) string {
	// 简化实现，实际应该解析 AI 返回的结构化大纲
	return content
}

// PlotSuggestChain 情节建议链
type PlotSuggestChain struct {
	chain *compose.Chain[map[string]any, []string]
}

// NewPlotSuggestChain 创建情节建议链
func NewPlotSuggestChain(components *components.Components) (*PlotSuggestChain, error) {
	chain := compose.NewChain[map[string]any, []string]()

	// Node 1: 格式化提示词
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
			template := components.ChatTemplate
			args := map[string]any{
				"NovelTitle":    input["novel_title"],
				"Genre":         input["genre"],
				"Summary":       input["summary"],
				"WorldSettings": input["world_settings"],
			}
			return template.Format(ctx, "suggest_plot", args)
		}),
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("PlotSuggestion"))

	// Node 3: 解析建议列表
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) ([]string, error) {
			return parseSuggestions(msg.Content), nil
		}),
		compose.WithNodeName("SuggestionParse"),
	)

	return &PlotSuggestChain{chain: chain}, nil
}

// Suggest 获取情节建议
func (c *PlotSuggestChain) Suggest(ctx context.Context, input map[string]any) ([]string, error) {
	return c.chain.Invoke(ctx, input)
}

func parseSuggestions(content string) []string {
	// 简化实现，按行分割
	lines := strings.Split(content, "\n")
	var suggestions []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			suggestions = append(suggestions, strings.TrimSpace(line))
		}
	}
	return suggestions
}
