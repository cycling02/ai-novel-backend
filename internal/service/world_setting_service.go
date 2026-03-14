package service

import (
	"context"
	"errors"

	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

var (
	ErrWorldSettingNotFound = errors.New("世界观设定不存在")
)

type WorldSettingService struct {
	worldSettingRepo *repository.WorldSettingRepository
	novelRepo        *repository.NovelRepository
}

func NewWorldSettingService(worldSettingRepo *repository.WorldSettingRepository, novelRepo *repository.NovelRepository) *WorldSettingService {
	return &WorldSettingService{
		worldSettingRepo: worldSettingRepo,
		novelRepo:        novelRepo,
	}
}

// CreateWorldSetting 创建设定
func (s *WorldSettingService) CreateWorldSetting(ctx context.Context, ws *model.WorldSetting) error {
	// 验证小说是否存在
	novel, err := s.novelRepo.GetByID(ctx, ws.NovelID)
	if err != nil {
		return ErrNovelNotFound
	}
	if novel == nil {
		return ErrNovelNotFound
	}
	return s.worldSettingRepo.Create(ws)
}

// GetWorldSetting 获取设定
func (s *WorldSettingService) GetWorldSetting(ctx context.Context, id string) (*model.WorldSetting, error) {
	return s.worldSettingRepo.GetByID(id)
}

// ListWorldSettings 列出设定
func (s *WorldSettingService) ListWorldSettings(ctx context.Context, novelID string) ([]model.WorldSetting, error) {
	return s.worldSettingRepo.GetByNovelID(novelID)
}

// ListWorldSettingsByCategory 按分类列出设定
func (s *WorldSettingService) ListWorldSettingsByCategory(ctx context.Context, novelID, category string) ([]model.WorldSetting, error) {
	return s.worldSettingRepo.GetByCategory(novelID, category)
}

// UpdateWorldSetting 更新设定
func (s *WorldSettingService) UpdateWorldSetting(ctx context.Context, ws *model.WorldSetting) error {
	existing, err := s.worldSettingRepo.GetByID(ws.ID)
	if err != nil {
		return ErrWorldSettingNotFound
	}
	if existing == nil {
		return ErrWorldSettingNotFound
	}
	return s.worldSettingRepo.Update(ws)
}

// DeleteWorldSetting 删除设定
func (s *WorldSettingService) DeleteWorldSetting(ctx context.Context, id string) error {
	return s.worldSettingRepo.Delete(id)
}