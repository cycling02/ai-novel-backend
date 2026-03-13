package agents

import (
	"context"

	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
)

// AgentType Agent 类型
type AgentType string

const (
	PlanningAgent  AgentType = "planning"
	WritingAgent   AgentType = "writing"
	EditingAgent   AgentType = "editing"
	ResearchAgent  AgentType = "research"
)

// NovelAgent 小说创作 Agent
type NovelAgent struct {
	agentType  AgentType
	components *components.Components
	systemPrompt string
}

// NewNovelAgent 创建小说创作 Agent
func NewNovelAgent(agentType AgentType, components *components.Components) *NovelAgent {
	prompts := map[AgentType]string{
		PlanningAgent: `你是一位专业的小说策划编辑。你的职责是：
1. 帮助作者规划小说整体结构
2. 设计剧情发展脉络
3. 确保故事逻辑自洽
4. 提供创作建议和指导`,

		WritingAgent: `你是一位才华横溢的网络小说作家。你的职责是：
1. 根据大纲创作精彩章节
2. 塑造生动的角色形象
3. 编写引人入胜的情节
4. 保持文风一致且富有感染力`,

		EditingAgent: `你是一位经验丰富的文字编辑。你的职责是：
1. 审核稿件质量
2. 润色文字表达
3. 修正逻辑漏洞
4. 提升阅读体验`,

		ResearchAgent: `你是一位专业的资料研究员。你的职责是：
1. 检索相关知识库
2. 整理世界观设定
3. 核实背景资料
4. 提供创作素材`,
	}

	return &NovelAgent{
		agentType:    agentType,
		components:   components,
		systemPrompt: prompts[agentType],
	}
}

// Chat 对话
func (a *NovelAgent) Chat(ctx context.Context, messages []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
	// 添加系统提示
	allMessages := append([]*schema.Message{
		{
			Role:    schema.System,
			Content: a.systemPrompt,
		},
	}, messages...)

	// 调用 ChatModel
	return a.components.ChatModel.Generate(ctx, allMessages)
}

// Stream 流式对话
func (a *NovelAgent) Stream(ctx context.Context, messages []*schema.Message, opts ...agent.AgentOption) (
	*schema.StreamReader[*schema.Message], error) {
	allMessages := append([]*schema.Message{
		{
			Role:    schema.System,
			Content: a.systemPrompt,
		},
	}, messages...)

	return a.components.ChatModel.Stream(ctx, allMessages)
}

// Transfer 转移给其他 Agent
func (a *NovelAgent) Transfer(targetType AgentType, context map[string]any) map[string]any {
	context["from_agent"] = a.agentType
	context["to_agent"] = targetType
	return context
}

// MultiAgentOrchestrator 多 Agent 编排器
type MultiAgentOrchestrator struct {
	agents map[AgentType]*NovelAgent
}

// NewMultiAgentOrchestrator 创建多 Agent 编排器
func NewMultiAgentOrchestrator(components *components.Components) *MultiAgentOrchestrator {
	return &MultiAgentOrchestrator{
		agents: map[AgentType]*NovelAgent{
			PlanningAgent:  NewNovelAgent(PlanningAgent, components),
			WritingAgent:   NewNovelAgent(WritingAgent, components),
			EditingAgent:   NewNovelAgent(EditingAgent, components),
			ResearchAgent:  NewNovelAgent(ResearchAgent, components),
		},
	}
}

// GetAgent 获取指定 Agent
func (o *MultiAgentOrchestrator) GetAgent(agentType AgentType) *NovelAgent {
	return o.agents[agentType]
}

// ExecuteWorkflow 执行工作流
func (o *MultiAgentOrchestrator) ExecuteWorkflow(ctx context.Context, workflow []AgentType, input map[string]any) (map[string]any, error) {
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: input["prompt"].(string),
		},
	}

	var lastResult *schema.Message
	var err error

	for _, agentType := range workflow {
		agent := o.agents[agentType]
		lastResult, err = agent.Chat(ctx, messages)
		if err != nil {
			return nil, err
		}

		// 将结果作为下一轮的输入
		messages = append(messages, lastResult)
	}

	return map[string]any{
		"result": lastResult.Content,
	}, nil
}
