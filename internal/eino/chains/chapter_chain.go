package chains

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// ChapterGenerateChain 章节生成链
type ChapterGenerateChain struct {
	chain *compose.Chain[map[string]any, string]
}

// NewChapterGenerateChain 创建章节生成链
func NewChapterGenerateChain(components *components.Components) (*ChapterGenerateChain, error) {
	// 创建 Chain：输入 map -> 输出 string
	chain := compose.NewChain[map[string]any, string]()

	// Node 1: Lambda - 检索相关知识
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			// 从输入中提取查询信息
			novelID, _ := input["novel_id"].(string)
			_ = novelID

			// 如果有 Retriever，检索相关知识
			// 这里简化处理，直接传递输入
			input["retrieved_context"] = ""
			return input, nil
		}),
		compose.WithNodeName("KnowledgeRetrieval"),
	)

	// Node 2: ChatTemplate - 格式化提示词
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
			template := components.ChatTemplate

			args := map[string]any{
				"NovelTitle":    input["novel_title"],
				"Genre":         input["genre"],
				"ChapterTitle":  input["chapter_title"],
				"Outline":       input["outline"],
				"PrevContent":   input["prev_content"],
				"WorldSettings": input["world_settings"],
				"Characters":    input["characters"],
			}

			return template.Format(ctx, "generate_chapter", args)
		}),
		compose.WithNodeName("PromptFormat"),
	)

	// Node 3: ChatModel - 调用 LLM 生成内容
	chain = chain.AppendChatModel(
		components.ChatModel,
		compose.WithNodeName("ChapterGeneration"),
	)

	// Node 4: Lambda - 后处理
	chain = chain.AppendLambda(
		compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (string, error) {
			return msg.Content, nil
		}),
		compose.WithNodeName("OutputProcess"),
	)

	return &ChapterGenerateChain{chain: chain}, nil
}

// Generate 生成章节
func (c *ChapterGenerateChain) Generate(ctx context.Context, input map[string]any) (string, error) {
	result, err := c.chain.Invoke(ctx, input)
	if err != nil {
		return "", fmt.Errorf("生成章节失败：%w", err)
	}
	return result, nil
}

// GenerateStream 流式生成章节
func (c *ChapterGenerateChain) GenerateStream(ctx context.Context, input map[string]any) (*schema.StreamReader[string], error) {
	// 获取流式输出
	streamReader, err := c.chain.Stream(ctx, input)
	if err != nil {
		return nil, err
	}

	// 转换类型
	outSr, outSw := schema.Pipe[string](10)
	go func() {
		defer outSw.Close()
		for {
			msg, err := streamReader.Recv()
			if err != nil {
				break
			}
			outSw.Send(msg, nil)
		}
	}()

	return outSr, nil
}
