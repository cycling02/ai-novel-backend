package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino-ext/components/model/openai"
	openaiemb "github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cycling02/ai-novel-backend/internal/config"
	"github.com/cycling02/ai-novel-backend/internal/eino/tools"
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

	var emb embedding.Embedder
	var ret retriever.Retriever

	if cfg.Vector.APIKey != "" {
		emb, err = initEmbedding(cfg.LLM)
		if err != nil {
			fmt.Printf("警告：Embedding 初始化失败：%v\n", err)
		} else {
			ret, err = NewPineconeRetriever(cfg.Vector.APIKey, cfg.Vector.IndexName, cfg.Vector.Namespace, emb)
			if err != nil {
				fmt.Printf("警告：Retriever 初始化失败：%v\n", err)
			}
		}
	}

	return &Components{
		ChatModel:    chatModel,
		Embedding:    emb,
		Retriever:    ret,
		ChatTemplate: NewNovelChatTemplate(),
		Tools:        initNovelTools(),
	}, nil
}

// initChatModel 初始化 ChatModel（使用 OpenAI 兼容接口支持 DeepSeek）
func initChatModel(cfg config.LLMConfig) (model.ChatModel, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	// 使用 eino-ext 的 OpenAI 兼容模型
	return openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: baseURL,
	})
}

// initEmbedding 初始化 Embedding
func initEmbedding(cfg config.LLMConfig) (embedding.Embedder, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	return openaiemb.NewEmbedder(context.Background(), &openaiemb.EmbeddingConfig{
		Model:   "text-embedding-3-small",
		APIKey:  cfg.APIKey,
		BaseURL: baseURL,
	})
}

// initNovelTools 初始化小说创作工具
func initNovelTools() []tool.BaseTool {
	return []tool.BaseTool{
		tools.NewCharacterQueryTool(),
		tools.NewWorldSettingQueryTool(),
		tools.NewPlotOutlineTool(),
		tools.NewStyleTransferTool(),
	}
}
