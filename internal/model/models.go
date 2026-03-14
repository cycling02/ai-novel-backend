package model

import "time"

// Novel 小说
type Novel struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"` // 玄幻/都市/历史/科幻/言情等
	Status      string    `json:"status"` // drafting, publishing, completed
	WordCount   int       `json:"word_count"`
	ChapterCount int      `json:"chapter_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Chapter 章节
type Chapter struct {
	ID         string    `json:"id"`
	NovelID    string    `json:"novel_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Order      int       `json:"order"`
	WordCount  int       `json:"word_count"`
	Status     string    `json:"status"` // draft, published, archived
	IsAIGenerated bool   `json:"is_ai_generated"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Character 角色
type Character struct {
	ID          string            `json:"id"`
	NovelID     string            `json:"novel_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Role        string            `json:"role"` // protagonist, antagonist, supporting
	Gender      string            `json:"gender"`
	Age         string            `json:"age"`
	Attributes  map[string]string `json:"attributes"` // 性格、外貌、能力等
	Relations   map[string]string `json:"relations"`  // 与其他角色的关系
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// WorldSetting 世界观设定
type WorldSetting struct {
	ID        string    `json:"id"`
	NovelID   string    `json:"novel_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"` // magic_system, geography, history, society, culture, other
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Knowledge 知识库条目（用于向量检索）
type Knowledge struct {
	ID        string    `json:"id"`
	NovelID   string    `json:"novel_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // character, setting, plot, other
	VectorID  string    `json:"vector_id"` // Pinecone 向量 ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Outline 大纲
type Outline struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"` // 详细大纲内容
	Structure   []OutlineNode `json:"structure"` // 结构化大纲
	IsAIGenerated bool     `json:"is_ai_generated"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OutlineNode 大纲节点
type OutlineNode struct {
	Arc       string   `json:"arc"`       // 故事篇章
	Chapters  []int    `json:"chapters"`  // 章节范围
	Summary   string   `json:"summary"`   // 篇章摘要
	KeyEvents []string `json:"key_events"` // 关键事件
}

// PlotSuggestion 情节建议
type PlotSuggestion struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	ChapterID   string    `json:"chapter_id,omitempty"`
	Suggestions []SuggestionItem `json:"suggestions"`
	CreatedAt   time.Time `json:"created_at"`
}

// SuggestionItem 单个建议
type SuggestionItem struct {
	Content   string  `json:"content"`    // 建议内容
	Reason    string  `json:"reason"`     // 推荐理由
	Confidence float32 `json:"confidence"` // 置信度
	Risk      string  `json:"risk"`       // 可能风险
}



// GenerationRequest AI 生成请求
type GenerationRequest struct {
	NovelID      string `json:"novel_id"`
	ChapterTitle string `json:"chapter_title"`
	Outline      string `json:"outline"`
	PrevContent  string `json:"prev_content,omitempty"`
	Style        string `json:"style,omitempty"`        // 文风
	Tone         string `json:"tone,omitempty"`         // 语调
	WordCount    int    `json:"word_count,omitempty"`   // 目标字数
}

// GenerationResponse AI 生成响应
type GenerationResponse struct {
	Content   string            `json:"content"`
	WordCount int               `json:"word_count"`
	Metadata  GenerationMetadata `json:"metadata"`
}

// GenerationMetadata 生成元数据
type GenerationMetadata struct {
	Model       string  `json:"model"`
	DurationMs  int64   `json:"duration_ms"`
	TokensUsed  int     `json:"tokens_used"`
	Confidence  float32 `json:"confidence"`
	IsStreaming bool    `json:"is_streaming"`
}
