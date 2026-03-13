package database

import (
	"context"
	"log"
)

// Migrate 执行数据库迁移
func Migrate(db *PostgresDB) error {
	ctx := context.Background()

	migrations := []string{
		// 用户表
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			nickname VARCHAR(100) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			avatar VARCHAR(500) DEFAULT '',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 小说表
		`CREATE TABLE IF NOT EXISTS novels (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			genre VARCHAR(100),
			status VARCHAR(50) DEFAULT 'drafting',
			word_count INTEGER DEFAULT 0,
			chapter_count INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 章节表
		`CREATE TABLE IF NOT EXISTS chapters (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			"order" INTEGER NOT NULL,
			word_count INTEGER DEFAULT 0,
			status VARCHAR(50) DEFAULT 'draft',
			is_ai_generated BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 角色表
		`CREATE TABLE IF NOT EXISTS characters (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			role VARCHAR(50),
			gender VARCHAR(20),
			age VARCHAR(50),
			attributes JSONB DEFAULT '{}',
			relations JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 世界观设定表
		`CREATE TABLE IF NOT EXISTS world_settings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			category VARCHAR(100),
			tags TEXT[] DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 知识库表（向量检索）
		`CREATE TABLE IF NOT EXISTS knowledge (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			type VARCHAR(50),
			vector_id VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 大纲表
		`CREATE TABLE IF NOT EXISTS outlines (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			structure JSONB DEFAULT '[]',
			is_ai_generated BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 情节建议表
		`CREATE TABLE IF NOT EXISTS plot_suggestions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			novel_id UUID REFERENCES novels(id) ON DELETE CASCADE,
			chapter_id UUID,
			suggestions JSONB DEFAULT '[]',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// 创建索引
		`CREATE INDEX IF NOT EXISTS idx_chapters_novel_id ON chapters(novel_id)`,
		`CREATE INDEX IF NOT EXISTS idx_chapters_order ON chapters(novel_id, "order")`,
		`CREATE INDEX IF NOT EXISTS idx_characters_novel_id ON characters(novel_id)`,
		`CREATE INDEX IF NOT EXISTS idx_world_settings_novel_id ON world_settings(novel_id)`,
		`CREATE INDEX IF NOT EXISTS idx_knowledge_novel_id ON knowledge(novel_id)`,
		`CREATE INDEX IF NOT EXISTS idx_knowledge_type ON knowledge(type)`,
		`CREATE INDEX IF NOT EXISTS idx_outlines_novel_id ON outlines(novel_id)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(ctx, migration); err != nil {
			log.Printf("警告：迁移 %d 失败：%v", i, err)
			return err
		}
	}

	log.Println("✅ 数据库迁移完成")
	return nil
}
