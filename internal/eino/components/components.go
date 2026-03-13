package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/document"
	"github.com/cycling02/ai-novel-backend/internal/config"
)

// Components 包含所有 Eino 组件
type Components struct {
	ChatModel    model.ChatModel
	Embedding    embedding.Embedder
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
	var emb embedding.Embedder
	var ret retriever.Retriever

	if cfg.Vector.APIKey != "" {
		emb, err = initEmbedding(cfg.LLM)
		if err != nil {
			return nil, fmt.Errorf("初始化 Embedding 失败：%w", err)
		}

		ret, err = initRetriever(cfg.Vector, emb)
		if err != nil {
			return nil, fmt.Errorf("初始化 Retriever 失败：%w", err)
		}
	}

	chatTemplate := initChatTemplate()

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
	// 使用 eino-ext 的 OpenAI 兼容模型
	// 注意：实际项目中需要正确导入 eino-ext 包
	return NewOpenAICompatibleModel(cfg.BaseURL, cfg.APIKey, cfg.Model)
}

// initEmbedding 初始化 Embedding
func initEmbedding(cfg config.LLMConfig) (embedding.Embedder, error) {
	return NewOpenAICompatibleEmbedding(cfg.BaseURL, cfg.APIKey, "text-embedding-3-small")
}

// initRetriever 初始化 Retriever（Pinecone）
func initRetriever(cfg config.VectorConfig, emb embedding.Embedder) (retriever.Retriever, error) {
	return NewPineconeRetriever(cfg.APIKey, cfg.IndexName, cfg.Namespace, emb)
}

// initChatTemplate 初始化 ChatTemplate
func initChatTemplate() prompt.ChatTemplate {
	return NewNovelChatTemplate()
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
