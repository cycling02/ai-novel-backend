package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	pinecone "github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// EmbeddingFunc 向量化函数类型
type EmbeddingFunc func(ctx context.Context, texts []string) ([][]float32, error)

// PineconeRetriever Pinecone 向量检索器
type PineconeRetriever struct {
	client    *pinecone.Client
	indexName string
	namespace string
	embedFn   EmbeddingFunc
}

// NewPineconeRetriever 创建 Pinecone Retriever
func NewPineconeRetriever(apiKey, indexName, namespace string, embedFn EmbeddingFunc) (retriever.Retriever, error) {
	client, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &PineconeRetriever{
		client:    client,
		indexName: indexName,
		namespace: namespace,
		embedFn:   embedFn,
	}, nil
}

// Retrieve 检索相关文档
func (r *PineconeRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	if r.embedFn == nil {
		return nil, fmt.Errorf("向量化函数未配置")
	}

	// 1. 将查询文本向量化
	embeddings, err := r.embedFn(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("向量化失败：%w", err)
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		// 向量化失败，返回空结果
		return []*schema.Document{}, nil
	}

	// 2. 查询 Pinecone
	index, err := r.client.Index(r.indexName)
	if err != nil {
		return nil, fmt.Errorf("获取索引失败：%w", err)
	}
	defer index.Close()

	resp, err := index.QueryVector(ctx, &pinecone.QueryVectorRequest{
		Vector:          embeddings[0],
		TopK:            5,
		Namespace:       r.namespace,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	// 3. 转换为 Eino Document
	docs := make([]*schema.Document, 0, len(resp.Matches))
	for _, match := range resp.Matches {
		content := ""
		if c, ok := match.Metadata["content"].(string); ok {
			content = c
		}
		docs = append(docs, &schema.Document{
			ID:      match.ID,
			Content: content,
			Meta:    match.Metadata,
		})
	}

	return docs, nil
}
