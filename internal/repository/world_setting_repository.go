package repository

import (
	"github.com/cycling02/ai-novel-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorldSettingRepository struct {
	db *gorm.DB
}

func NewWorldSettingRepository(db *gorm.DB) *WorldSettingRepository {
	return &WorldSettingRepository{db: db}
}

// Create 创建设定
func (r *WorldSettingRepository) Create(ws *model.WorldSetting) error {
	if ws.ID == "" {
		ws.ID = uuid.New().String()
	}
	return r.db.Create(ws).Error
}

// GetByID 根据 ID 获取设定
func (r *WorldSettingRepository) GetByID(id string) (*model.WorldSetting, error) {
	var ws model.WorldSetting
	err := r.db.First(&ws, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

// GetByNovelID 根据小说 ID 获取所有设定
func (r *WorldSettingRepository) GetByNovelID(novelID string) ([]model.WorldSetting, error) {
	var settings []model.WorldSetting
	err := r.db.Where("novel_id = ?", novelID).Order("category, created_at DESC").Find(&settings).Error
	return settings, err
}

// GetByCategory 按分类获取设定
func (r *WorldSettingRepository) GetByCategory(novelID, category string) ([]model.WorldSetting, error) {
	var settings []model.WorldSetting
	err := r.db.Where("novel_id = ? AND category = ?", novelID, category).Find(&settings).Error
	return settings, err
}

// Update 更新设定
func (r *WorldSettingRepository) Update(ws *model.WorldSetting) error {
	return r.db.Save(ws).Error
}

// Delete 删除设定
func (r *WorldSettingRepository) Delete(id string) error {
	return r.db.Delete(&model.WorldSetting{}, "id = ?", id).Error
}

// CountByNovelID 统计小说设定数
func (r *WorldSettingRepository) CountByNovelID(novelID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.WorldSetting{}).Where("novel_id = ?", novelID).Count(&count).Error
	return count, err
}