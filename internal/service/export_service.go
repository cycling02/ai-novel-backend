package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

type ExportService struct {
	novelRepo   *repository.NovelRepository
	chapterRepo *repository.ChapterRepository
}

func NewExportService(novelRepo *repository.NovelRepository, chapterRepo *repository.ChapterRepository) *ExportService {
	return &ExportService{
		novelRepo:   novelRepo,
		chapterRepo: chapterRepo,
	}
}

// Export 导出小说
func (s *ExportService) Export(ctx context.Context, req *model.ExportRequest) (*model.ExportResponse, error) {
	// 获取小说信息
	novel, err := s.novelRepo.GetByID(ctx, req.NovelID)
	if err != nil || novel == nil {
		return nil, ErrNovelNotFound
	}

	// 获取章节列表
	chapters, err := s.chapterRepo.ListByNovel(ctx, req.NovelID)
	if err != nil {
		return nil, err
	}

	// 过滤章节范围
	if req.StartChap > 0 || req.EndChap > 0 {
		var filtered []*model.Chapter
		for _, ch := range chapters {
			if req.StartChap > 0 && ch.Order < req.StartChap {
				continue
			}
			if req.EndChap > 0 && ch.Order > req.EndChap {
				continue
			}
			filtered = append(filtered, ch)
		}
		chapters = filtered
	}

	// 按格式导出
	var content string
	var fileName string

	switch req.Format {
	case "txt":
		content, fileName = s.exportToTXT(novel, chapters)
	case "markdown", "md":
		content, fileName = s.exportToMarkdown(novel, chapters)
	case "epub":
		fileName = s.exportToEPUB(novel, chapters)
	case "pdf":
		fileName = s.exportToPDF(novel, chapters)
	default:
		content, fileName = s.exportToTXT(novel, chapters)
	}

	return &model.ExportResponse{
		FileName: fileName,
		Content:  content,
		Size:     int64(len(content)),
	}, nil
}

// exportToTXT 导出为 TXT
func (s *ExportService) exportToTXT(novel *model.Novel, chapters []*model.Chapter) (string, string) {
	var sb strings.Builder

	// 标题
	sb.WriteString(novel.Title)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	// 简介
	if novel.Description != "" {
		sb.WriteString("【简介】\n")
		sb.WriteString(novel.Description)
		sb.WriteString("\n\n")
	}

	// 章节
	for _, ch := range chapters {
		sb.WriteString(fmt.Sprintf("第%d章 %s\n", ch.Order, ch.Title))
		sb.WriteString(strings.Repeat("-", 30))
		sb.WriteString("\n")
		sb.WriteString(ch.Content)
		sb.WriteString("\n\n")
	}

	fileName := fmt.Sprintf("%s_%s.txt", novel.Title, time.Now().Format("20060102"))
	return sb.String(), fileName
}

// exportToMarkdown 导出为 Markdown
func (s *ExportService) exportToMarkdown(novel *model.Novel, chapters []*model.Chapter) (string, string) {
	var sb strings.Builder

	// 标题
	sb.WriteString("# ")
	sb.WriteString(novel.Title)
	sb.WriteString("\n\n")

	// 元信息
	sb.WriteString("---")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("**类型**: %s\n", novel.Genre))
	sb.WriteString(fmt.Sprintf("**字数**: %d\n", novel.WordCount))
	sb.WriteString(fmt.Sprintf("**章节数**: %d\n", novel.ChapterCount))
	sb.WriteString("---\n\n")

	// 简介
	if novel.Description != "" {
		sb.WriteString("## 简介\n\n")
		sb.WriteString(novel.Description)
		sb.WriteString("\n\n")
	}

	// 章节
	for _, ch := range chapters {
		sb.WriteString(fmt.Sprintf("## 第%d章 %s\n\n", ch.Order, ch.Title))
		sb.WriteString(ch.Content)
		sb.WriteString("\n\n")
	}

	fileName := fmt.Sprintf("%s_%s.md", novel.Title, time.Now().Format("20060102"))
	return sb.String(), fileName
}

// exportToEPUB 导出为 EPUB（返回文件名，实际生成需要库支持）
func (s *ExportService) exportToEPUB(novel *model.Novel, chapters []*model.Chapter) string {
	// TODO: 使用第三方库如 github.com/BurntSushi/toml 生成 EPUB
	// 这里只返回文件名
	return fmt.Sprintf("%s_%s.epub", novel.Title, time.Now().Format("20060102"))
}

// exportToPDF 导出为 PDF（返回文件名，实际生成需要库支持）
func (s *ExportService) exportToPDF(novel *model.Novel, chapters []*model.Chapter) string {
	// TODO: 使用第三方库如 github.com/jung-kurt/gofpdf 生成 PDF
	// 这里只返回文件名
	return fmt.Sprintf("%s_%s.pdf", novel.Title, time.Now().Format("20060102"))
}

// ExportOutline 导出大纲
func (s *ExportService) ExportOutline(ctx context.Context, novelID string) (string, error) {
	novel, err := s.novelRepo.GetByID(ctx, novelID)
	if err != nil || novel == nil {
		return "", ErrNovelNotFound
	}

	// TODO: 获取大纲数据并导出
	return "", nil
}