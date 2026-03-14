package service

import (
	"context"
	"errors"

	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

var (
	ErrCharacterNotFound = errors.New("角色不存在")
	ErrNovelNotFound    = errors.New("小说不存在")
)

type CharacterService struct {
	characterRepo *repository.CharacterRepository
	novelRepo     *repository.NovelRepository
}

func NewCharacterService(characterRepo *repository.CharacterRepository, novelRepo *repository.NovelRepository) *CharacterService {
	return &CharacterService{
		characterRepo: characterRepo,
		novelRepo:     novelRepo,
	}
}

// CreateCharacter 创建角色
func (s *CharacterService) CreateCharacter(ctx context.Context, character *model.Character) error {
	// 验证小说是否存在
	novel, err := s.novelRepo.GetByID(ctx, character.NovelID)
	if err != nil {
		return ErrNovelNotFound
	}
	if novel == nil {
		return ErrNovelNotFound
	}
	return s.characterRepo.Create(character)
}

// GetCharacter 获取角色
func (s *CharacterService) GetCharacter(ctx context.Context, id string) (*model.Character, error) {
	return s.characterRepo.GetByID(id)
}

// ListCharacters 列出角色
func (s *CharacterService) ListCharacters(ctx context.Context, novelID string) ([]model.Character, error) {
	return s.characterRepo.GetByNovelID(novelID)
}

// UpdateCharacter 更新角色
func (s *CharacterService) UpdateCharacter(ctx context.Context, character *model.Character) error {
	existing, err := s.characterRepo.GetByID(character.ID)
	if err != nil {
		return ErrCharacterNotFound
	}
	if existing == nil {
		return ErrCharacterNotFound
	}
	return s.characterRepo.Update(character)
}

// DeleteCharacter 删除角色
func (s *CharacterService) DeleteCharacter(ctx context.Context, id string) error {
	return s.characterRepo.Delete(id)
}