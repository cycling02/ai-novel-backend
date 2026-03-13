package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB PostgreSQL 数据库连接
type PostgresDB struct {
	Pool *pgxpool.Pool
}

// NewPostgresDB 创建数据库连接
func NewPostgresDB(url string) (*PostgresDB, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("创建连接池失败：%w", err)
	}

	// 测试连接
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("连接数据库失败：%w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

// Close 关闭连接池
func (db *PostgresDB) Close() {
	db.Pool.Close()
}

// Ping 测试连接
func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Query 查询
func (db *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, value := range values {
			row[string(rows.FieldDescriptions()[i].Name)] = value
		}
		results = append(results, row)
	}

	return results, nil
}

// Exec 执行
func (db *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// QueryRow 查询单行
func (db *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	results, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}
