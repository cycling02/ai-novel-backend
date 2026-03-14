package repository

import (
	"github.com/cycling02/ai-novel-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharacterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// Create 创建角色
func (r *CharacterRepository) Create(character *model.Character) error {
	if character.ID == "" {
		character.ID = uuid.New().String()
	}
	return r.db.Create(character).Error
}

// GetByID 根据 ID 获取角色
func (r *CharacterRepository) GetByID(id string) (*model.Character, error) {
	var character model.Character
	err := r.db.First(&character, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// GetByNovelID 根据小说 ID 获取所有角色
func (r *CharacterRepository) GetByNovelID(novelID string) ([]model.Character, error) {
	var characters []model.Character
	err := r.db.Where("novel_id = ?", novelID).Order("created_at DESC").Find(&characters).Error
	return characters, err
}

// Update 更新角色
func (r *CharacterRepository) Update(character *model.Character) error {
	return r.db.Save(character).Error
}

// Delete 删除角色
func (r *CharacterRepository) Delete(id string) error {
	return r.db.Delete(&model.Character{}, "id = ?", id).Error
}

// CountByNovelID 统计小说角色数
func (r *CharacterRepository) CountByNovelID(novelID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Character{}).Where("novel_id = ?", novelID).Count(&count).Error
	return count, err
}