package service

import (
	"context"
	"time"

	"github.com/cycling02/ai-novel-backend/internal/eino/agents"
	"github.com/cycling02/ai-novel-backend/internal/eino/chains"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

// NovelService 小说业务服务
type NovelService struct {
	novelRepo       *repository.NovelRepository
	chapterRepo     *repository.ChapterRepository
	characterRepo   *repository.CharacterRepository
	worldRepo       *repository.WorldSettingRepository
	knowledgeRepo   *repository.KnowledgeRepository
	
	// Eino 组件
	chapterChain    *chains.ChapterGenerateChain
	outlineChain    *chains.OutlineExpandChain
	plotSuggestChain *chains.PlotSuggestChain
	editChain       *chains.ContentEditChain
	
	// 多 Agent 编排器
	orchestrator    *agents.MultiAgentOrchestrator
}

// NewNovelService 创建小说服务
func NewNovelService(
	novelRepo *repository.NovelRepository,
	chapterRepo *repository.ChapterRepository,
	characterRepo *repository.CharacterRepository,
	worldRepo *repository.WorldSettingRepository,
	knowledgeRepo *repository.KnowledgeRepository,
	chapterChain *chains.ChapterGenerateChain,
	outlineChain *chains.OutlineExpandChain,
	plotSuggestChain *chains.PlotSuggestChain,
	editChain *chains.ContentEditChain,
	orchestrator *agents.MultiAgentOrchestrator,
) *NovelService {
	return &NovelService{
		novelRepo:        novelRepo,
		chapterRepo:      chapterRepo,
		characterRepo:    characterRepo,
		worldRepo:        worldRepo,
		knowledgeRepo:    knowledgeRepo,
		chapterChain:     chapterChain,
		outlineChain:     outlineChain,
		plotSuggestChain: plotSuggestChain,
		editChain:        editChain,
		orchestrator:     orchestrator,
	}
}

// CreateNovel 创建小说
func (s *NovelService) CreateNovel(ctx context.Context, userID, title, description, genre string) (*model.Novel, error) {
	novel := &model.Novel{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Genre:        genre,
		Status:       "drafting",
		WordCount:    0,
		ChapterCount: 0,
	}

	if err := s.novelRepo.Create(ctx, novel); err != nil {
		return nil, err
	}

	return novel, nil
}

// GetNovel 获取小说详情
func (s *NovelService) GetNovel(ctx context.Context, id string) (*model.Novel, error) {
	return s.novelRepo.GetByID(ctx, id)
}

// ListNovels 获取小说列表
func (s *NovelService) ListNovels(ctx context.Context, userID string) ([]*model.Novel, error) {
	return s.novelRepo.ListByUser(ctx, userID)
}

// UpdateNovel 更新小说
func (s *NovelService) UpdateNovel(ctx context.Context, novel *model.Novel) error {
	return s.novelRepo.Update(ctx, novel)
}

// DeleteNovel 删除小说
func (s *NovelService) DeleteNovel(ctx context.Context, id string) error {
	return s.novelRepo.Delete(ctx, id)
}

// GenerateChapter AI 生成章节
func (s *NovelService) GenerateChapter(ctx context.Context, req *model.GenerationRequest) (*model.GenerationResponse, error) {
	// 获取小说信息
	novel, err := s.novelRepo.GetByID(ctx, req.NovelID)
	if err != nil {
		return nil, err
	}

	// 获取前文内容
	prevContent := req.PrevContent
	if prevContent == "" {
		lastChapter, _ := s.chapterRepo.GetLastChapter(ctx, req.NovelID)
		if lastChapter != nil {
			prevContent = lastChapter.Content
		}
	}

	// 获取世界观和角色信息
	worldSettings := "" // TODO: 从数据库获取
	characters := ""    // TODO: 从数据库获取

	// 调用 Eino Chain 生成
	input := map[string]any{
		"novel_id":       req.NovelID,
		"novel_title":    novel.Title,
		"genre":          novel.Genre,
		"chapter_title":  req.ChapterTitle,
		"outline":        req.Outline,
		"prev_content":   prevContent,
		"world_settings": worldSettings,
		"characters":     characters,
	}

	startTime := time.Now()
	content, err := s.chapterChain.Generate(ctx, input)
	if err != nil {
		return nil, err
	}

	return &model.GenerationResponse{
		Content:   content,
		WordCount: len([]rune(content)),
		Metadata: model.GenerationMetadata{
			Model:       "deepseek-chat",
			DurationMs:  time.Since(startTime).Milliseconds(),
			TokensUsed:  0, // TODO: 从 API 响应获取
			Confidence:  0.9,
			IsStreaming: false,
		},
	}, nil
}

// GenerateOutline AI 生成大纲
func (s *NovelService) GenerateOutline(ctx context.Context, novelID, briefOutline string) (string, error) {
	novel, err := s.novelRepo.GetByID(ctx, novelID)
	if err != nil {
		return "", err
	}

	input := map[string]any{
		"novel_title":   novel.Title,
		"genre":         novel.Genre,
		"brief_outline": briefOutline,
	}

	return s.outlineChain.Expand(ctx, input)
}

// SuggestPlot AI 情节建议
func (s *NovelService) SuggestPlot(ctx context.Context, novelID string) ([]model.SuggestionItem, error) {
	novel, err := s.novelRepo.GetByID(ctx, novelID)
	if err != nil {
		return nil, err
	}

	summary, _ := s.chapterRepo.GetSummary(ctx, novelID, 3)

	input := map[string]any{
		"novel_title":    novel.Title,
		"genre":          novel.Genre,
		"summary":        summary,
		"world_settings": "", // TODO: 获取世界观
	}

	suggestions, err := s.plotSuggestChain.Suggest(ctx, input)
	if err != nil {
		return nil, err
	}

	// 转换为结构体
	result := make([]model.SuggestionItem, len(suggestions))
	for i, s := range suggestions {
		result[i] = model.SuggestionItem{
			Content:   s,
			Reason:    "",
			Confidence: 0.8,
			Risk:      "",
		}
	}

	return result, nil
}

// EditContent AI 润色编辑
func (s *NovelService) EditContent(ctx context.Context, content string) (map[string]any, error) {
	input := map[string]any{
		"content": content,
	}

	return s.editChain.Edit(ctx, input)
}
