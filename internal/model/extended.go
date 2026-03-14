package model

import (
	"time"
)

// ChapterVersion 章节版本
type ChapterVersion struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	ChapterID   string    `json:"chapter_id" gorm:"index;not null"`
	Version     int       `json:"version" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text"`
	Title       string    `json:"title"`
	WordCount   int       `json:"word_count"`
	ChangeDesc  string    `json:"change_desc"` // 变更描述
	ChangeType  string    `json:"change_type"` // ai_generate, manual_edit, rollback
	CreatedBy   string    `json:"created_by"`  // user_id or 'ai'
	CreatedAt   time.Time `json:"created_at"`
}

// VersionRequest 版本请求
type VersionRequest struct {
	ChapterID  string `json:"chapter_id"`
	ChangeDesc string `json:"change_desc"`
}

// BatchGenerateRequest 批量生成请求
type BatchGenerateRequest struct {
	StartTitle string `json:"start_title"` // 起始章节标题
	Count      int    `json:"count"`        // 生成章节数量
	WordCount  int    `json:"word_count"`   // 每章字数
	Style      string `json:"style"`        // 文风
	Tone       string `json:"tone"`         // 语调
}

// BatchGenerateResponse 批量生成响应
type BatchGenerateResponse struct {
	Chapters   []Chapter `json:"chapters"`
	TotalWords int       `json:"total_words"`
	Success    bool      `json:"success"`
	Error      string    `json:"error,omitempty"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	NovelID   string `json:"novel_id" binding:"required"`
	Format    string `json:"format" binding:"required,oneof=txt epub markdown pdf"` // 导出格式
	StartChap int    `json:"start_chap"`  // 起始章节（可选）
	EndChap   int    `json:"end_chap"`    // 结束章节（可选）
	Include   string `json:"include"`     // 包含内容: all, chapter_only, outline_only
}

// ExportResponse 导出响应
type ExportResponse struct {
	FileName string `json:"file_name"`
	Content  string `json:"content,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	Size     int64  `json:"size"`
}

// StreamEvent 流式事件
type StreamEvent struct {
	Type    string      `json:"type"` // start, chunk, error, done
	Content string      `json:"content,omitempty"`
	Delta   string      `json:"delta,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// GenerateStreamRequest 流式生成请求
type GenerateStreamRequest struct {
	NovelID      string `json:"novel_id"`
	ChapterTitle string `json:"chapter_title"`
	Outline      string `json:"outline"`
	PrevContent  string `json:"prev_content,omitempty"`
	Style        string `json:"style,omitempty"`
	Tone         string `json:"tone,omitempty"`
	WordCount    int    `json:"word_count,omitempty"`
}

// NovelOutline 小说大纲（完整版）
type NovelOutline struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid"`
	NovelID     string         `json:"novel_id" gorm:"index;not null"`
	Title       string         `json:"title"`
	Content     string         `json:"content" gorm:"type:text"` // 完整大纲内容
	Summary     string         `json:"summary"`                  // 一句话简介
	Genre       string         `json:"genre"`                    // 分类
	TargetWord  int            `json:"target_word"`             // 目标字数
	Chapters    int            `json:"chapters"`                // 目标章节数
	Characters  []string       `json:"characters"`               // 主要角色列表
	WorldSetup  []string       `json:"world_setup"`             // 世界观设定列表
	Structure   []OutlineNode  `json:"structure"`               // 结构化大纲
	IsAIGenerated bool        `json:"is_ai_generated"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// GenerateOutlineRequest 生成大纲请求
type GenerateOutlineRequest struct {
	Title       string   `json:"title" binding:"required"`
	Genre       string   `json:"genre"`
	Summary     string   `json:"summary"`      // 一句话简介
	Keywords    []string `json:"keywords"`      // 关键词
	TargetWord  int      `json:"target_word"`   // 目标字数
	TargetChapters int   `json:"target_chapters"` // 目标章节数
}

// GenerateOutlineResponse 生成大纲响应
type GenerateOutlineResponse struct {
	Outline    *NovelOutline `json:"outline"`
	DurationMs int64         `json:"duration_ms"`
}