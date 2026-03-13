package components

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/config"
)

// Embedding 简单的 Embedding 实现
type Embedding struct {
	baseURL   string
	apiKey    string
	modelName string
}

// NewEmbedding 创建 Embedding
func NewEmbedding(cfg config.LLMConfig) embedding.Embedder {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	return &Embedding{
		baseURL:   baseURL,
		apiKey:    cfg.APIKey,
		modelName: "text-embedding-3-small",
	}
}

// EmbedStrings 向量化文本
func (e *Embedding) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float32, error) {
	// TODO: 实现真实的 API 调用
	// 返回占位向量
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = make([]float32, 1536) // OpenAI embedding 维度
	}
	return result, nil
}

// EmbedDocuments 向量化文档
func (e *Embedding) EmbedDocuments(ctx context.Context, docs []*schema.Document, opts ...embedding.Option) ([][]float32, error) {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.Content
	}
	return e.EmbedStrings(ctx, texts, opts...)
}
