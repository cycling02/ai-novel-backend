package service

import (
	"context"
	"errors"
	"time"

	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

var (
	ErrChapterNotFound  = errors.New("章节不存在")
	ErrVersionNotFound  = errors.New("版本不存在")
	ErrInvalidVersion   = errors.New("无效的版本号")
)

type VersionService struct {
	versionRepo *repository.ChapterVersionRepository
	chapterRepo *repository.ChapterRepository
}

func NewVersionService(versionRepo *repository.ChapterVersionRepository, chapterRepo *repository.ChapterRepository) *VersionService {
	return &VersionService{
		versionRepo: versionRepo,
		chapterRepo: chapterRepo,
	}
}

// CreateVersion 创建新版本（保存当前章节内容为新版本）
func (s *VersionService) CreateVersion(ctx context.Context, chapterID, content, title, changeDesc, changeType, createdBy string) (*model.ChapterVersion, error) {
	// 获取章节
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, ErrChapterNotFound
	}
	if chapter == nil {
		return nil, ErrChapterNotFound
	}

	// 获取最新版本号
	latestVersion, err := s.versionRepo.GetLatestVersion(chapterID)
	versionNum := 1
	if err == nil && latestVersion != nil {
		versionNum = latestVersion.Version + 1
	}

	// 计算字数
	wordCount := len([]rune(content))

	// 创建新版本
	version := &model.ChapterVersion{
		ChapterID:  chapterID,
		Version:    versionNum,
		Content:    content,
		Title:      title,
		WordCount:  wordCount,
		ChangeDesc: changeDesc,
		ChangeType: changeType,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	if err := s.versionRepo.Create(version); err != nil {
		return nil, err
	}

	return version, nil
}

// GetVersions 获取章节所有版本
func (s *VersionService) GetVersions(ctx context.Context, chapterID string) ([]model.ChapterVersion, error) {
	return s.versionRepo.GetByChapterID(chapterID)
}

// GetVersion 获取指定版本
func (s *VersionService) GetVersion(ctx context.Context, chapterID string, versionNum int) (*model.ChapterVersion, error) {
	version, err := s.versionRepo.GetVersionByNumber(chapterID, versionNum)
	if err != nil {
		return nil, ErrVersionNotFound
	}
	return version, nil
}

// Rollback 回滚到指定版本
func (s *VersionService) Rollback(ctx context.Context, chapterID string, versionNum int, userID string) (*model.ChapterVersion, error) {
	// 获取要回滚的版本
	targetVersion, err := s.versionRepo.GetVersionByNumber(chapterID, versionNum)
	if err != nil {
		return nil, ErrVersionNotFound
	}

	// 获取当前章节
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, ErrChapterNotFound
	}

	// 创建新版本（保存当前内容）
	_, err = s.CreateVersion(ctx, chapterID, chapter.Content, chapter.Title, "回滚前备份", "rollback", userID)
	if err != nil {
		return nil, err
	}

	// 更新章节内容
	chapter.Content = targetVersion.Content
	chapter.WordCount = len([]rune(targetVersion.Content))
	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		return nil, err
	}

	// 创建回滚后的新版本
	newVersion, err := s.CreateVersion(ctx, chapterID, targetVersion.Content, targetVersion.Title, "从版本回滚", "rollback", userID)
	if err != nil {
		return nil, err
	}

	return newVersion, nil
}

// CompareVersions 对比两个版本
func (s *VersionService) CompareVersions(ctx context.Context, chapterID string, v1, v2 int) (string, string, error) {
	version1, err := s.versionRepo.GetVersionByNumber(chapterID, v1)
	if err != nil {
		return "", "", ErrVersionNotFound
	}

	version2, err := s.versionRepo.GetVersionByNumber(chapterID, v2)
	if err != nil {
		return "", "", ErrVersionNotFound
	}

	return version1.Content, version2.Content, nil
}