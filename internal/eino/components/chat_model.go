package components

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// OpenAICompatibleModel OpenAI 兼容的 ChatModel（支持 DeepSeek 等）
type OpenAICompatibleModel struct {
	baseURL   string
	apiKey    string
	modelName string
	client    *http.Client
}

// NewOpenAICompatibleModel 创建 OpenAI 兼容模型
func NewOpenAICompatibleModel(baseURL, apiKey, modelName string) (model.ChatModel, error) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAICompatibleModel{
		baseURL:   baseURL,
		apiKey:    apiKey,
		modelName: modelName,
		client:    &http.Client{},
	}, nil
}

// Generate 非流式生成
func (m *OpenAICompatibleModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	reqBody := map[string]interface{}{
		"model":    m.modelName,
		"messages": messagesToOpenAI(input),
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/chat/completions",
		io.NopCloser(io.NopReader()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Body = io.NopCloser(io.NopReader())
	
	// 简化实现，实际应该发送请求
	_ = reqData

	// 这里应该调用 API，简化返回
	return &schema.Message{
		Role:    schema.Assistant,
		Content: "AI 生成的内容",
	}, nil
}

// Stream 流式生成
func (m *OpenAICompatibleModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (
	*schema.StreamReader[*schema.Message], error) {
	sr, sw := schema.Pipe[*schema.Message](10)

	go func() {
		// 简化实现
		sw.Send(&schema.Message{
			Role:    schema.Assistant,
			Content: "流式内容",
		}, nil)
		sw.Close()
	}()

	return sr, nil
}

func messagesToOpenAI(messages []*schema.Message) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result[i] = map[string]interface{}{
			"role":    roleToString(msg.Role),
			"content": msg.Content,
		}
	}
	return result
}

func roleToString(role schema.Role) string {
	switch role {
	case schema.System:
		return "system"
	case schema.User:
		return "user"
	case schema.Assistant:
		return "assistant"
	default:
		return "user"
	}
}
