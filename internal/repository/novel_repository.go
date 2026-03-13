package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/cycling02/ai-novel-backend/internal/database"
	"github.com/cycling02/ai-novel-backend/internal/model"
)

// NovelRepository 小说数据访问层
type NovelRepository struct {
	db *database.PostgresDB
}

// NewNovelRepository 创建小说仓库
func NewNovelRepository(db *database.PostgresDB) *NovelRepository {
	return &NovelRepository{db: db}
}

// Create 创建小说
func (r *NovelRepository) Create(ctx context.Context, novel *model.Novel) error {
	novel.ID = uuid.New().String()
	novel.CreatedAt = time.Now()
	novel.UpdatedAt = time.Now()

	query := `INSERT INTO novels (id, user_id, title, description, genre, status, word_count, chapter_count, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(ctx, query,
		novel.ID, novel.UserID, novel.Title, novel.Description,
		novel.Genre, novel.Status, novel.WordCount, novel.ChapterCount,
		novel.CreatedAt, novel.UpdatedAt)

	return err
}

// GetByID 根据 ID 获取小说
func (r *NovelRepository) GetByID(ctx context.Context, id string) (*model.Novel, error) {
	query := `SELECT id, user_id, title, description, genre, status, word_count, chapter_count, created_at, updated_at
			  FROM novels WHERE id = $1`

	row, err := r.db.QueryRow(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}

	return mapToNovel(row), nil
}

// ListByUser 获取用户的小说列表
func (r *NovelRepository) ListByUser(ctx context.Context, userID string) ([]*model.Novel, error) {
	query := `SELECT id, user_id, title, description, genre, status, word_count, chapter_count, created_at, updated_at
			  FROM novels WHERE user_id = $1 ORDER BY updated_at DESC`

	results, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	novels := make([]*model.Novel, len(results))
	for i, row := range results {
		novels[i] = mapToNovel(row)
	}
	return novels, nil
}

// Update 更新小说
func (r *NovelRepository) Update(ctx context.Context, novel *model.Novel) error {
	novel.UpdatedAt = time.Now()

	query := `UPDATE novels SET title=$2, description=$3, genre=$4, status=$5, 
			  word_count=$6, chapter_count=$7, updated_at=$8 WHERE id=$1`

	_, err := r.db.Exec(ctx, query,
		novel.ID, novel.Title, novel.Description, novel.Genre,
		novel.Status, novel.WordCount, novel.ChapterCount, novel.UpdatedAt)

	return err
}

// Delete 删除小说
func (r *NovelRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM novels WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdateWordCount 更新字数统计
func (r *NovelRepository) UpdateWordCount(ctx context.Context, id string, wordCount, chapterCount int) error {
	query := `UPDATE novels SET word_count=$2, chapter_count=$3, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id, wordCount, chapterCount)
	return err
}

func mapToNovel(row map[string]interface{}) *model.Novel {
	return &model.Novel{
		ID:           row["id"].(string),
		UserID:       row["user_id"].(string),
		Title:        row["title"].(string),
		Description:  getString(row["description"]),
		Genre:        getString(row["genre"]),
		Status:       getString(row["status"]),
		WordCount:    getInt(row["word_count"]),
		ChapterCount: getInt(row["chapter_count"]),
		CreatedAt:    getTime(row["created_at"]),
		UpdatedAt:    getTime(row["updated_at"]),
	}
}
