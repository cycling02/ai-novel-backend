package components

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/config"
)

// Components 包含所有 Eino 组件
type Components struct {
	ChatModel    model.ChatModel
	EmbedFn      EmbeddingFunc
	Retriever    retriever.Retriever
	ChatTemplate prompt.ChatTemplate
	Tools        []tool.BaseTool
}

// EmbeddingFunc 向量化函数类型
type EmbeddingFunc func(ctx context.Context, texts []string) ([][]float32, error)

// InitComponents 初始化所有 Eino 组件
func InitComponents(cfg *config.Config) (*Components, error) {
	chatModel, err := initChatModel(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("初始化 ChatModel 失败：%w", err)
	}

	// embedding 和 retriever 可选
	var embedFn EmbeddingFunc
	var ret retriever.Retriever

	if cfg.Vector.APIKey != "" && cfg.LLM.APIKey != "" {
		embedFn = createEmbeddingFunc(cfg.LLM.BaseURL, cfg.LLM.APIKey)

		ret, err = NewPineconeRetriever(cfg.Vector.APIKey, cfg.Vector.IndexName, cfg.Vector.Namespace, embedFn)
		if err != nil {
			// Retriever 初始化失败不影响整体
			fmt.Printf("警告：Retriever 初始化失败：%v\n", err)
		}
	}

	chatTemplate := NewNovelChatTemplate()
	tools := initNovelTools()

	return &Components{
		ChatModel:    chatModel,
		EmbedFn:      embedFn,
		Retriever:    ret,
		ChatTemplate: chatTemplate,
		Tools:        tools,
	}, nil
}

// initChatModel 初始化 ChatModel（使用 OpenAI 兼容接口支持 DeepSeek）
func initChatModel(cfg config.LLMConfig) (model.ChatModel, error) {
	return NewOpenAICompatibleModel(cfg.BaseURL, cfg.APIKey, cfg.Model)
}

// createEmbeddingFunc 创建向量化函数
func createEmbeddingFunc(baseURL, apiKey string) EmbeddingFunc {
	client := &http.Client{}
	
	return func(ctx context.Context, texts []string) ([][]float32, error) {
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}

		reqBody := map[string]interface{}{
			"model": "text-embedding-3-small",
			"input": texts,
		}

		reqData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/embeddings",
			bytes.NewReader(reqData))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// 解析响应
		var apiResp struct {
			Data []struct {
				Embedding []float32 `json:"embedding"`
			} `json:"data"`
		}

		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return nil, err
		}

		embeddings := make([][]float32, len(apiResp.Data))
		for i, d := range apiResp.Data {
			embeddings[i] = d.Embedding
		}

		return embeddings, nil
	}
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
