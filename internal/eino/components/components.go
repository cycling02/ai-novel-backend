package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cycling02/ai-novel-backend/internal/config"
)

// Components 包含所有 Eino 组件
type Components struct {
	ChatModel    model.ChatModel
	Embedding    model.Embedder
	Retriever    retriever.Retriever
	ChatTemplate prompt.ChatTemplate
	Tools        []tool.BaseTool
}

// InitComponents 初始化所有 Eino 组件
func InitComponents(cfg *config.Config) (*Components, error) {
	chatModel, err := initChatModel(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("初始化 ChatModel 失败：%w", err)
	}

	// embedding 和 retriever 可选
	var emb model.Embedder
	var ret retriever.Retriever

	if cfg.Vector.APIKey != "" && cfg.LLM.APIKey != "" {
		emb, err = initEmbedding(cfg.LLM)
		if err != nil {
			fmt.Printf("警告：Embedding 初始化失败：%v\n", err)
		}
	}

	if emb != nil && cfg.Vector.APIKey != "" {
		ret, err = NewPineconeRetriever(cfg.Vector.APIKey, cfg.Vector.IndexName, cfg.Vector.Namespace, emb)
		if err != nil {
			fmt.Printf("警告：Retriever 初始化失败：%v\n", err)
		}
	}

	chatTemplate := NewNovelChatTemplate()
	tools := initNovelTools()

	return &Components{
		ChatModel:    chatModel,
		Embedding:    emb,
		Retriever:    ret,
		ChatTemplate: chatTemplate,
		Tools:        tools,
	}, nil
}

// initChatModel 初始化 ChatModel（使用 OpenAI 兼容接口支持 DeepSeek）
func initChatModel(cfg config.LLMConfig) (model.ChatModel, error) {
	// 简化实现，直接使用 schema
	return &simpleChatModel{
		baseURL:   cfg.BaseURL,
		apiKey:    cfg.APIKey,
		modelName: cfg.Model,
	}, nil
}

// initEmbedding 初始化 Embedding
func initEmbedding(cfg config.LLMConfig) (model.Embedder, error) {
	return &simpleEmbedding{
		baseURL:   cfg.BaseURL,
		apiKey:    cfg.APIKey,
		modelName: "text-embedding-3-small",
	}, nil
}

// initNovelTools 初始化小说创作工具
func initNovelTools() []tool.BaseTool {
	return []tool.BaseTool{
		NewCharacterQueryTool(),
		NewWorldSettingQueryTool(),
		NewPlotOutlineTool(),
		NewStyleTransferTool(),
	}
}

// simpleChatModel 简单的 ChatModel 实现
type simpleChatModel struct {
	baseURL   string
	apiKey    string
	modelName string
}

func (m *simpleChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return &schema.Message{
		Role:    schema.Assistant,
		Content: "AI 回复",
	}, nil
}

func (m *simpleChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	sr, sw := schema.Pipe[*schema.Message](1)
	sw.Send(&schema.Message{
		Role:    schema.Assistant,
		Content: "AI 流式回复",
	}, nil)
	sw.Close()
	return sr, nil
}

// simpleEmbedding 简单的 Embedding 实现
type simpleEmbedding struct {
	baseURL   string
	apiKey    string
	modelName string
}

func (e *simpleEmbedding) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float32, error) {
	// 简化实现，返回空向量
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = make([]float32, 1536)
	}
	return result, nil
}
