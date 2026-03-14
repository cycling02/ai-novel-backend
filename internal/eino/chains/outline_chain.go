package chains

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// OutlineExpandChain 大纲扩写链
type OutlineExpandChain struct {
	chain   *compose.Chain[map[string]any, string]
	runnable compose.Runnable[map[string]any, string]
}

// NewOutlineExpandChain 创建大纲扩写链
func NewOutlineExpandChain(components *components.Components) (*OutlineExpandChain, error) {
	chain := compose.NewChain[map[string]any, string]()

	// Node 1: ChatTemplate
	chain = chain.AppendChatTemplate(
		components.ChatTemplate,
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("OutlineExpansion"))

	// Node 3: 解析输出
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input any) (string, error) {
			content := ""
			switch msg := input.(type) {
			case *schema.Message:
				content = msg.Content
			case schema.Message:
				content = msg.Content
			default:
				content = fmt.Sprintf("%v", input)
			}
			// 解析结构化大纲
			return parseOutline(content), nil
		}),
		compose.WithNodeName("OutputParse"),
	)

	// 编译 Chain
	rCtx := context.Background()
	runnable, err := chain.Compile(rCtx)
	if err != nil {
		return nil, fmt.Errorf("编译 Chain 失败：%w", err)
	}

	return &OutlineExpandChain{
		chain:   chain,
		runnable: runnable,
	}, nil
}

// Expand 扩写大纲
func (c *OutlineExpandChain) Expand(ctx context.Context, input map[string]any) (string, error) {
	args := map[string]any{
		"template_name": "expand_outline",
		"NovelTitle":   input["novel_title"],
		"Genre":        input["genre"],
		"BriefOutline": input["brief_outline"],
	}
	return c.runnable.Invoke(ctx, args)
}

func parseOutline(content string) string {
	// 简化实现，实际应该解析 AI 返回的结构化大纲
	return content
}

// PlotSuggestChain 情节建议链
type PlotSuggestChain struct {
	chain   *compose.Chain[map[string]any, []string]
	runnable compose.Runnable[map[string]any, []string]
}

// NewPlotSuggestChain 创建情节建议链
func NewPlotSuggestChain(components *components.Components) (*PlotSuggestChain, error) {
	chain := compose.NewChain[map[string]any, []string]()

	// Node 1: ChatTemplate
	chain = chain.AppendChatTemplate(
		components.ChatTemplate,
		compose.WithNodeName("PromptFormat"),
	)

	// Node 2: ChatModel
	chain = chain.AppendChatModel(components.ChatModel, compose.WithNodeName("PlotSuggestion"))

	// Node 3: 解析建议列表
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input any) ([]string, error) {
			content := ""
			switch msg := input.(type) {
			case *schema.Message:
				content = msg.Content
			case schema.Message:
				content = msg.Content
			default:
				content = fmt.Sprintf("%v", input)
			}
			return parseSuggestions(content), nil
		}),
		compose.WithNodeName("SuggestionParse"),
	)

	// 编译 Chain
	rCtx := context.Background()
	runnable, err := chain.Compile(rCtx)
	if err != nil {
		return nil, fmt.Errorf("编译 Chain 失败：%w", err)
	}

	return &PlotSuggestChain{
		chain:   chain,
		runnable: runnable,
	}, nil
}

// Suggest 获取情节建议
func (c *PlotSuggestChain) Suggest(ctx context.Context, input map[string]any) ([]string, error) {
	args := map[string]any{
		"template_name":  "suggest_plot",
		"NovelTitle":    input["novel_title"],
		"Genre":         input["genre"],
		"Summary":       input["summary"],
		"WorldSettings": input["world_settings"],
	}
	return c.runnable.Invoke(ctx, args)
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