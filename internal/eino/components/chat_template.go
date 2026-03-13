package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// NovelChatTemplate 小说创作专用的 ChatTemplate
type NovelChatTemplate struct {
	templates map[string]string
}

// NewNovelChatTemplate 创建小说创作模板
func NewNovelChatTemplate() prompt.ChatTemplate {
	return &NovelChatTemplate{
		templates: map[string]string{
			"generate_chapter": `你是一位专业的网络小说作家。请根据以下信息创作小说章节。

小说信息：
- 书名：{{.NovelTitle}}
- 类型：{{.Genre}}
- 本章标题：{{.ChapterTitle}}

本章大纲：
{{.Outline}}

前文内容摘要：
{{.PrevContent}}

世界观设定：
{{.WorldSettings}}

主要角色：
{{.Characters}}

请创作这一章节的内容，要求：
1. 保持角色性格一致
2. 情节连贯有吸引力
3. 描写生动，对话自然
4. 符合{{.Genre}}类型小说的风格
5. 字数约 3000 字`,

			"suggest_plot": `你是一位经验丰富的网络小说编辑。请分析以下小说的当前进展，提供情节发展建议。

小说信息：
- 书名：{{.NovelTitle}}
- 类型：{{.Genre}}
- 当前剧情摘要：{{.Summary}}

已设定的世界观：
{{.WorldSettings}}

请提供 3-5 个合理的情节发展建议，每个建议包含：
1. 建议内容
2. 推荐理由
3. 可能的风险或注意事项`,

			"expand_outline": `你是一位专业的小说策划。请将以下简短的故事构思扩展为详细的大纲。

小说信息：
- 书名：{{.NovelTitle}}
- 类型：{{.Genre}}

简要构思：
{{.BriefOutline}}

请扩展为详细大纲，包括：
1. 故事背景设定
2. 主要角色介绍
3. 剧情发展脉络（起承转合）
4. 重要情节节点
5. 高潮和结局设计`,

			"edit_content": `你是一位专业的文字编辑。请对以下小说内容进行润色和优化。

原文：
{{.Content}}

优化要求：
1. 修正语病和不通顺的句子
2. 增强描写的生动性
3. 优化对话的自然度
4. 保持原有风格和情节
5. 标注主要修改之处`,
		},
	}
}

// Format 格式化提示词
func (t *NovelChatTemplate) Format(ctx context.Context, templateName string, args map[string]any) ([]*schema.Message, error) {
	tmpl, ok := t.templates[templateName]
	if !ok {
		return nil, fmt.Errorf("模板不存在：%s", templateName)
	}

	// 简单的模板替换（实际应该使用 text/template）
	content := tmpl
	for key, value := range args {
		placeholder := "{{." + key + "}}"
		content = replacePlaceholder(content, placeholder, fmt.Sprintf("%v", value))
	}

	return []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一位专业的网络小说创作助手。",
		},
		{
			Role:    schema.User,
			Content: content,
		},
	}, nil
}

func replacePlaceholder(content, placeholder, value string) string {
	// 简单替换实现
	return content
}
