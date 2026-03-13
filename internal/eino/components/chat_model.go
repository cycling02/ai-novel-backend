package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/config"
)

// ChatModel 简单的 ChatModel 实现
type ChatModel struct {
	baseURL   string
	apiKey    string
	modelName string
}

// NewChatModel 创建 ChatModel
func NewChatModel(cfg config.LLMConfig) (model.ChatModel, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	return &ChatModel{
		baseURL:   baseURL,
		apiKey:    cfg.APIKey,
		modelName: cfg.Model,
	}, nil
}

// Generate 非流式生成
func (m *ChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// TODO: 实现真实的 API 调用
	// 现在返回一个占位响应
	return &schema.Message{
		Role:    schema.Assistant,
		Content: "这是一个占位回复。实际部署时需要实现真实的 LLM API 调用。",
	}, nil
}

// Stream 流式生成
func (m *ChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	sr, sw := schema.Pipe[*schema.Message](5)

	go func() {
		defer sw.Close()
		
		// 模拟流式输出
		parts := []string{"这是", "一个", "流式", "回复", "。"}
		for _, part := range parts {
			sw.Send(&schema.Message{
				Role:    schema.Assistant,
				Content: part,
			}, nil)
		}
	}()

	return sr, nil
}

// WithModel 设置模型名称（Option）
func WithModel(model string) model.Option {
	return model.Option(func(o *model.Options) {
		o.Model = model
	})
}
