package tools

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// CharacterQueryTool 角色查询工具
type CharacterQueryTool struct{}

// NewCharacterQueryTool 创建角色查询工具
func NewCharacterQueryTool() tool.BaseTool {
	return &CharacterQueryTool{}
}

// Info 工具信息
func (t *CharacterQueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "character_query",
		Desc: "查询小说中的角色信息，包括姓名、性格、外貌、关系等",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"novel_id": {
				Desc:     "小说 ID",
				Required: true,
				Type:     schema.String,
			},
			"character_name": {
				Desc:     "角色名称（可选，不填则返回所有角色）",
				Required: false,
				Type:     schema.String,
			},
		}),
	}, nil
}

// Invoke 执行工具
func (t *CharacterQueryTool) Invoke(ctx context.Context, argsInJSON string) (resultInJSON string, err error) {
	var args struct {
		NovelID       string `json:"novel_id"`
		CharacterName string `json:"character_name"`
	}
	json.Unmarshal([]byte(argsInJSON), &args)

	// 实际实现应该查询数据库
	result := map[string]interface{}{
		"characters": []map[string]interface{}{
			{"name": "示例角色", "role": "主角", "description": "这是一个示例角色"},
		},
	}
	data, _ := json.Marshal(result)
	return string(data), nil
}

// WorldSettingQueryTool 世界观查询工具
type WorldSettingQueryTool struct{}

// NewWorldSettingQueryTool 创建世界观查询工具
func NewWorldSettingQueryTool() tool.BaseTool {
	return &WorldSettingQueryTool{}
}

func (t *WorldSettingQueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "world_setting_query",
		Desc: "查询小说的世界观设定，包括魔法体系、地理环境、历史背景、社会规则等",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"novel_id": {
				Desc:     "小说 ID",
				Required: true,
				Type:     schema.String,
			},
			"category": {
				Desc:     "设定类别（magic/geography/history/society/other）",
				Required: false,
				Type:     schema.String,
			},
			"keyword": {
				Desc:     "搜索关键词",
				Required: false,
				Type:     schema.String,
			},
		}),
	}, nil
}

func (t *WorldSettingQueryTool) Invoke(ctx context.Context, argsInJSON string) (resultInJSON string, err error) {
	var args struct {
		NovelID  string `json:"novel_id"`
		Category string `json:"category"`
		Keyword  string `json:"keyword"`
	}
	json.Unmarshal([]byte(argsInJSON), &args)

	result := map[string]interface{}{
		"settings": []map[string]interface{}{
			{"title": "魔法体系", "content": "这是一个示例设定"},
		},
	}
	data, _ := json.Marshal(result)
	return string(data), nil
}

// PlotOutlineTool 大纲生成工具
type PlotOutlineTool struct{}

// NewPlotOutlineTool 创建大纲生成工具
func NewPlotOutlineTool() tool.BaseTool {
	return &PlotOutlineTool{}
}

func (t *PlotOutlineTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "plot_outline",
		Desc: "根据当前剧情生成后续大纲建议",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"novel_id": {
				Desc:     "小说 ID",
				Required: true,
				Type:     schema.String,
			},
			"current_chapter": {
				Desc:     "当前章节序号",
				Required: false,
				Type:     schema.Integer,
			},
			"direction": {
				Desc:     "期望的剧情走向（action/romance/mystery/character_growth）",
				Required: false,
				Type:     schema.String,
			},
		}),
	}, nil
}

func (t *PlotOutlineTool) Invoke(ctx context.Context, argsInJSON string) (resultInJSON string, err error) {
	var args struct {
		NovelID        string `json:"novel_id"`
		CurrentChapter int    `json:"current_chapter"`
		Direction      string `json:"direction"`
	}
	json.Unmarshal([]byte(argsInJSON), &args)

	result := map[string]interface{}{
		"outline": "这是生成的大纲建议...",
		"options": []string{"选项 1", "选项 2", "选项 3"},
	}
	data, _ := json.Marshal(result)
	return string(data), nil
}

// StyleTransferTool 风格转换工具
type StyleTransferTool struct{}

// NewStyleTransferTool 创建风格转换工具
func NewStyleTransferTool() tool.BaseTool {
	return &StyleTransferTool{}
}

func (t *StyleTransferTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "style_transfer",
		Desc: "将文本转换为指定的写作风格",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "原始内容",
				Required: true,
				Type:     schema.String,
			},
			"style": {
				Desc:     "目标风格（wuxia/fantasy/romance/scifi/historical）",
				Required: true,
				Type:     schema.String,
			},
			"tone": {
				Desc:     "语调（serious/humorous/dramatic/calm）",
				Required: false,
				Type:     schema.String,
			},
		}),
	}, nil
}

func (t *StyleTransferTool) Invoke(ctx context.Context, argsInJSON string) (resultInJSON string, err error) {
	var args struct {
		Content string `json:"content"`
		Style   string `json:"style"`
		Tone    string `json:"tone"`
	}
	json.Unmarshal([]byte(argsInJSON), &args)

	result := map[string]interface{}{
		"original": args.Content,
		"styled":   "风格转换后的内容...",
		"style":    args.Style,
	}
	data, _ := json.Marshal(result)
	return string(data), nil
}
