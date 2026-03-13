package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	pinecone "github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// PineconeRetriever Pinecone 向量检索器
type PineconeRetriever struct {
	client    *pinecone.Client
	indexName string
	namespace string
	embedding embedding.Embedder
}

// NewPineconeRetriever 创建 Pinecone Retriever
func NewPineconeRetriever(apiKey, indexName, namespace string, emb embedding.Embedder) (retriever.Retriever, error) {
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
		embedding: emb,
	}, nil
}

// Retrieve 检索相关文档
func (r *PineconeRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	// 1. 将查询文本向量化
	embeddings, err := r.embedding.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("向量化失败：%w", err)
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
	docs := make([]*schema.Document, len(resp.Matches))
	for i, match := range resp.Matches {
		docs[i] = &schema.Document{
			ID:      match.ID,
			Content: match.Metadata["content"].(string),
			Meta:    match.Metadata,
		}
	}

	return docs, nil
}
