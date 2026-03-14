package repository

import (
	"context"
	"fmt"

	"github.com/cycling02/ai-novel-backend/internal/database"
	"github.com/cycling02/ai-novel-backend/internal/model"
)

// KnowledgeRepository 知识库仓库
type KnowledgeRepository struct {
	db *database.PostgresDB
}

// NewKnowledgeRepository 创建知识库仓库
func NewKnowledgeRepository(db *database.PostgresDB) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

// Create 创建知识条目
func (r *KnowledgeRepository) Create(ctx context.Context, knowledge *model.Knowledge) error {
	query := `
		INSERT INTO knowledge (id, novel_id, title, content, type, vector_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		knowledge.ID,
		knowledge.NovelID,
		knowledge.Title,
		knowledge.Content,
		knowledge.Type,
		knowledge.VectorID,
		knowledge.CreatedAt,
		knowledge.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建知识条目失败：%w", err)
	}
	return nil
}

// GetByID 根据 ID 获取
func (r *KnowledgeRepository) GetByID(ctx context.Context, id string) (*model.Knowledge, error) {
	query := `
		SELECT id, novel_id, title, content, type, vector_id, created_at, updated_at
		FROM knowledge
		WHERE id = $1
	`
	row, err := r.db.QueryRow(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("获取知识条目失败：%w", err)
	}
	if row == nil {
		return nil, nil
	}

	knowledge := &model.Knowledge{
		ID:        row["id"].(string),
		NovelID:   row["novel_id"].(string),
		Title:     getString(row["title"]),
		Content:   getString(row["content"]),
		Type:      getString(row["type"]),
		VectorID:  getString(row["vector_id"]),
		CreatedAt: getTime(row["created_at"]),
		UpdatedAt: getTime(row["updated_at"]),
	}
	return knowledge, nil
}

// GetByNovelID 根据小说 ID 获取知识列表
func (r *KnowledgeRepository) GetByNovelID(ctx context.Context, novelID string) ([]*model.Knowledge, error) {
	query := `
		SELECT id, novel_id, title, content, type, vector_id, created_at, updated_at
		FROM knowledge
		WHERE novel_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, novelID)
	if err != nil {
		return nil, fmt.Errorf("查询知识列表失败：%w", err)
	}

	var results []*model.Knowledge
	for _, row := range rows {
		knowledge := &model.Knowledge{
			ID:        row["id"].(string),
			NovelID:   row["novel_id"].(string),
			Title:     getString(row["title"]),
			Content:   getString(row["content"]),
			Type:      getString(row["type"]),
			VectorID:  getString(row["vector_id"]),
			CreatedAt: getTime(row["created_at"]),
			UpdatedAt: getTime(row["updated_at"]),
		}
		results = append(results, knowledge)
	}
	return results, nil
}

// GetByType 根据类型获取
func (r *KnowledgeRepository) GetByType(ctx context.Context, novelID, knowledgeType string) ([]*model.Knowledge, error) {
	query := `
		SELECT id, novel_id, title, content, type, vector_id, created_at, updated_at
		FROM knowledge
		WHERE novel_id = $1 AND type = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, novelID, knowledgeType)
	if err != nil {
		return nil, fmt.Errorf("查询知识列表失败：%w", err)
	}

	var results []*model.Knowledge
	for _, row := range rows {
		knowledge := &model.Knowledge{
			ID:        row["id"].(string),
			NovelID:   row["novel_id"].(string),
			Title:     getString(row["title"]),
			Content:   getString(row["content"]),
			Type:      getString(row["type"]),
			VectorID:  getString(row["vector_id"]),
			CreatedAt: getTime(row["created_at"]),
			UpdatedAt: getTime(row["updated_at"]),
		}
		results = append(results, knowledge)
	}
	return results, nil
}

// Update 更新知识条目
func (r *KnowledgeRepository) Update(ctx context.Context, knowledge *model.Knowledge) error {
	query := `
		UPDATE knowledge
		SET title = $1, content = $2, type = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		knowledge.Title,
		knowledge.Content,
		knowledge.Type,
		knowledge.UpdatedAt,
		knowledge.ID,
	)
	if err != nil {
		return fmt.Errorf("更新知识条目失败：%w", err)
	}
	return nil
}

// Delete 删除知识条目
func (r *KnowledgeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM knowledge WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("删除知识条目失败：%w", err)
	}
	return nil
}