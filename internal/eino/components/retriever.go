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
	if r.embedding == nil {
		return nil, fmt.Errorf("embedding 未配置")
	}

	// 1. 向量化查询
	embeddings, err := r.embedding.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("向量化失败：%w", err)
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return []*schema.Document{}, nil
	}

	// 转换 float64 到 float32
	vectorFloat32 := make([]float32, len(embeddings[0]))
	for i, v := range embeddings[0] {
		vectorFloat32[i] = float32(v)
	}

	// 2. 描述索引获取 Host
	indexDesc, err := r.client.DescribeIndex(ctx, r.indexName)
	if err != nil {
		return nil, fmt.Errorf("获取索引信息失败：%w", err)
	}

	// 3. 创建索引连接
	idxConnParams := pinecone.NewIndexConnParams{
		Host:      indexDesc.Host,
		Namespace: r.namespace,
	}
	index, err := r.client.Index(idxConnParams)
	if err != nil {
		return nil, fmt.Errorf("连接索引失败：%w", err)
	}
	defer index.Close()

	// 4. 查询 Pinecone
	resp, err := index.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          vectorFloat32,
		TopK:            5,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	// 5. 转换为 Document
	docs := make([]*schema.Document, 0, len(resp.Matches))
	for _, match := range resp.Matches {
		if match.Vector == nil {
			continue
		}

		content := ""
		if match.Vector.Metadata != nil {
			metaMap := match.Vector.Metadata.AsMap()
			if c, ok := metaMap["content"].(string); ok {
				content = c
			}
		}

		// 转换 metadata
		meta := make(map[string]any)
		if match.Vector.Metadata != nil {
			meta = match.Vector.Metadata.AsMap()
		}

		docs = append(docs, &schema.Document{
			ID:       match.Vector.Id,
			Content:  content,
			MetaData: meta,
		})
	}

	return docs, nil
}