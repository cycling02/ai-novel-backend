package repository

import (
	"github.com/cycling02/ai-novel-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterVersionRepository struct {
	db *gorm.DB
}

func NewChapterVersionRepository(db *gorm.DB) *ChapterVersionRepository {
	return &ChapterVersionRepository{db: db}
}

// Create 创建版本
func (r *ChapterVersionRepository) Create(version *model.ChapterVersion) error {
	if version.ID == "" {
		version.ID = uuid.New().String()
	}
	return r.db.Create(version).Error
}

// GetByID 根据 ID 获取版本
func (r *ChapterVersionRepository) GetByID(id string) (*model.ChapterVersion, error) {
	var version model.ChapterVersion
	err := r.db.First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetByChapterID 获取章节所有版本
func (r *ChapterVersionRepository) GetByChapterID(chapterID string) ([]model.ChapterVersion, error) {
	var versions []model.ChapterVersion
	err := r.db.Where("chapter_id = ?", chapterID).Order("version DESC").Find(&versions).Error
	return versions, err
}

// GetLatestVersion 获取最新版本
func (r *ChapterVersionRepository) GetLatestVersion(chapterID string) (*model.ChapterVersion, error) {
	var version model.ChapterVersion
	err := r.db.Where("chapter_id = ?", chapterID).Order("version DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetVersionByNumber 获取指定版本
func (r *ChapterVersionRepository) GetVersionByNumber(chapterID string, versionNum int) (*model.ChapterVersion, error) {
	var version model.ChapterVersion
	err := r.db.Where("chapter_id = ? AND version = ?", chapterID, versionNum).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// DeleteByChapterID 删除章节所有版本
func (r *ChapterVersionRepository) DeleteByChapterID(chapterID string) error {
	return r.db.Where("chapter_id = ?", chapterID).Delete(&model.ChapterVersion{}).Error
}

// CountVersions 统计版本数
func (r *ChapterVersionRepository) CountVersions(chapterID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.ChapterVersion{}).Where("chapter_id = ?", chapterID).Count(&count).Error
	return count, err
}