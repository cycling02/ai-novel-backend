package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
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
	chatModel, err := NewChatModel(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("初始化 ChatModel 失败：%w", err)
	}

	var emb embedding.Embedder
	var ret retriever.Retriever

	if cfg.Vector.APIKey != "" {
		emb = NewEmbedding(cfg.LLM)
		ret, err = NewPineconeRetriever(cfg.Vector.APIKey, cfg.Vector.IndexName, cfg.Vector.Namespace, emb)
		if err != nil {
			fmt.Printf("警告：Retriever 初始化失败：%v\n", err)
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

// initNovelTools 初始化小说创作工具
func initNovelTools() []tool.BaseTool {
	return []tool.BaseTool{
		NewCharacterQueryTool(),
		NewWorldSettingQueryTool(),
		NewPlotOutlineTool(),
		NewStyleTransferTool(),
	}
}
