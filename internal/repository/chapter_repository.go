package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/cycling02/ai-novel-backend/internal/database"
	"github.com/cycling02/ai-novel-backend/internal/model"
)

// ChapterRepository 章节数据访问层
type ChapterRepository struct {
	db *database.PostgresDB
}

// NewChapterRepository 创建章节仓库
func NewChapterRepository(db *database.PostgresDB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

// Create 创建章节
func (r *ChapterRepository) Create(ctx context.Context, chapter *model.Chapter) error {
	chapter.ID = uuid.New().String()
	chapter.CreatedAt = time.Now()
	chapter.UpdatedAt = time.Now()

	query := `INSERT INTO chapters (id, novel_id, title, content, "order", word_count, status, is_ai_generated, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(ctx, query,
		chapter.ID, chapter.NovelID, chapter.Title, chapter.Content,
		chapter.Order, chapter.WordCount, chapter.Status, chapter.IsAIGenerated,
		chapter.CreatedAt, chapter.UpdatedAt)

	return err
}

// GetByID 根据 ID 获取章节
func (r *ChapterRepository) GetByID(ctx context.Context, id string) (*model.Chapter, error) {
	query := `SELECT id, novel_id, title, content, "order", word_count, status, is_ai_generated, created_at, updated_at
			  FROM chapters WHERE id = $1`

	row, err := r.db.QueryRow(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}

	return mapToChapter(row), nil
}

// ListByNovel 获取小说的章节列表
func (r *ChapterRepository) ListByNovel(ctx context.Context, novelID string) ([]*model.Chapter, error) {
	query := `SELECT id, novel_id, title, content, "order", word_count, status, is_ai_generated, created_at, updated_at
			  FROM chapters WHERE novel_id = $1 ORDER BY "order" ASC`

	results, err := r.db.Query(ctx, query, novelID)
	if err != nil {
		return nil, err
	}

	chapters := make([]*model.Chapter, len(results))
	for i, row := range results {
		chapters[i] = mapToChapter(row)
	}
	return chapters, nil
}

// GetLastChapter 获取最后一章
func (r *ChapterRepository) GetLastChapter(ctx context.Context, novelID string) (*model.Chapter, error) {
	query := `SELECT id, novel_id, title, content, "order", word_count, status, is_ai_generated, created_at, updated_at
			  FROM chapters WHERE novel_id = $1 ORDER BY "order" DESC LIMIT 1`

	row, err := r.db.QueryRow(ctx, query, novelID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	return mapToChapter(row), nil
}

// GetByOrder 根据序号获取章节
func (r *ChapterRepository) GetByOrder(ctx context.Context, novelID string, order int) (*model.Chapter, error) {
	query := `SELECT id, novel_id, title, content, "order", word_count, status, is_ai_generated, created_at, updated_at
			  FROM chapters WHERE novel_id = $1 AND "order" = $2`

	row, err := r.db.QueryRow(ctx, query, novelID, order)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}

	return mapToChapter(row), nil
}

// Update 更新章节
func (r *ChapterRepository) Update(ctx context.Context, chapter *model.Chapter) error {
	chapter.UpdatedAt = time.Now()

	query := `UPDATE chapters SET title=$2, content=$3, "order"=$4, word_count=$5, 
			  status=$6, is_ai_generated=$7, updated_at=$8 WHERE id=$1`

	_, err := r.db.Exec(ctx, query,
		chapter.ID, chapter.Title, chapter.Content, chapter.Order,
		chapter.WordCount, chapter.Status, chapter.IsAIGenerated, chapter.UpdatedAt)

	return err
}

// Delete 删除章节
func (r *ChapterRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM chapters WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetSummary 获取章节摘要（用于 AI 上下文）
func (r *ChapterRepository) GetSummary(ctx context.Context, novelID string, limit int) (string, error) {
	if limit <= 0 {
		limit = 3
	}
	query := `SELECT content FROM chapters WHERE novel_id = $1 ORDER BY "order" DESC LIMIT $2`
	results, err := r.db.Query(ctx, query, novelID, limit)
	if err != nil {
		return "", err
	}

	summary := ""
	for i := len(results) - 1; i >= 0; i-- {
		row := results[i]
		if content, ok := row["content"].(string); ok {
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			summary += content + "\n\n"
		}
	}

	return summary, nil
}

func mapToChapter(row map[string]interface{}) *model.Chapter {
	return &model.Chapter{
		ID:            row["id"].(string),
		NovelID:       row["novel_id"].(string),
		Title:         row["title"].(string),
		Content:       getString(row["content"]),
		Order:         getInt(row["order"]),
		WordCount:     getInt(row["word_count"]),
		Status:        getString(row["status"]),
		IsAIGenerated: getBool(row["is_ai_generated"]),
		CreatedAt:     getTime(row["created_at"]),
		UpdatedAt:     getTime(row["updated_at"]),
	}
}

func getBool(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
