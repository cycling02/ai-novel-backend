package graphs

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// NovelCreationGraph 小说创作图（多阶段流程）
type NovelCreationGraph struct {
	graph *compose.Graph[map[string]any, map[string]any]
}

// NovelCreationState 创作状态
type NovelCreationState struct {
	Stage       string
	NovelTitle  string
	Genre       string
	Outline     string
	ChapterTitle string
	Content     string
	Feedback    string
}

// NewNovelCreationGraph 创建小说创作图
func NewNovelCreationGraph(components *components.Components) (*NovelCreationGraph, error) {
	// 创建 Graph
	graph := compose.NewGraph[map[string]any, map[string]any]()

	// Node 1: 大纲生成
	graph.AddLambdaNode(
		"outline_generator",
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			// 调用大纲生成链
			output := map[string]any{
				"outline": "生成的详细大纲...",
				"stage":   "outline_completed",
			}
			// 合并输入
			for k, v := range input {
				output[k] = v
			}
			return output, nil
		}),
		compose.WithOutputKey("outline"),
	)

	// Node 2: 角色设定
	graph.AddLambdaNode(
		"character_creator",
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			output := map[string]any{
				"characters": "生成的角色设定...",
				"stage":      "characters_completed",
			}
			for k, v := range input {
				output[k] = v
			}
			return output, nil
		}),
		compose.WithInputKey("outline"),
	)

	// Node 3: 世界观设定
	graph.AddLambdaNode(
		"world_builder",
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			output := map[string]any{
				"world_settings": "生成的世界观...",
				"stage":          "world_completed",
			}
			for k, v := range input {
				output[k] = v
			}
			return output, nil
		}),
		compose.WithInputKey("outline"),
	)

	// Node 4: 章节生成（并行处理角色和世界观）
	graph.AddLambdaNode(
		"chapter_writer",
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			output := map[string]any{
				"content": "生成的章节内容...",
				"stage":   "chapter_completed",
			}
			for k, v := range input {
				output[k] = v
			}
			return output, nil
		}),
		compose.WithInputKeys([]string{"characters", "world_settings"}),
	)

	// Node 5: 内容审核
	graph.AddLambdaNode(
		"content_reviewer",
		compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			output := map[string]any{
				"reviewed": true,
				"feedback": "",
				"stage":    "review_completed",
			}
			for k, v := range input {
				output[k] = v
			}
			return output, nil
		}),
		compose.WithInputKey("content"),
	)

	// 添加边
	graph.AddEdge("outline_generator", "character_creator")
	graph.AddEdge("outline_generator", "world_builder")
	graph.AddEdge("character_creator", "chapter_writer")
	graph.AddEdge("world_builder", "chapter_writer")
	graph.AddEdge("chapter_writer", "content_reviewer")

	// 编译 Graph
	compiled, err := graph.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("编译 Graph 失败：%w", err)
	}

	return &NovelCreationGraph{graph: compiled}, nil
}

// Create 执行创作流程
func (g *NovelCreationGraph) Create(ctx context.Context, input map[string]any) (map[string]any, error) {
	return g.graph.Invoke(ctx, input)
}

// CreateStream 流式执行
func (g *NovelCreationGraph) CreateStream(ctx context.Context, input map[string]any) (*schema.StreamReader[map[string]any], error) {
	return g.graph.Stream(ctx, input)
}
