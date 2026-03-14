package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cycling02/ai-novel-backend/internal/config"
	"github.com/cycling02/ai-novel-backend/internal/database"
	"github.com/cycling02/ai-novel-backend/internal/server"
)

func main() {
	log.Println("🚀 AI Novel Backend 启动中...")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败：%v", err)
	}
	log.Println("✅ 配置加载完成")

	// 初始化 pgx 数据库
	log.Println("📦 连接数据库 (pgx)...")
	db, err := database.NewPostgresDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("❌ 连接数据库失败：%v", err)
	}
	defer db.Close()
	log.Println("✅ 数据库连接成功")

	// 初始化 GORM 数据库
	log.Println("📦 连接数据库 (GORM)...")
	gormDB, err := database.NewGormDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("❌ 连接 GORM 数据库失败：%v", err)
	}
	log.Println("✅ GORM 数据库连接成功")

	// 运行数据库迁移
	log.Println("🔧 执行数据库迁移...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ 数据库迁移失败：%v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 创建服务器
	log.Println("🌐 创建 HTTP 服务器...")
	srv, err := server.NewServer(cfg, db, gormDB)
	if err != nil {
		log.Fatalf("❌ 创建服务器失败：%v", err)
	}

	// 优雅关闭
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 启动服务器
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("📡 服务器启动在 %s", addr)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ 服务器错误：%v", err)
		}
	}()

	// 等待关闭信号
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("🛑 正在关闭服务器...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ 服务器关闭错误：%v", err)
	}
	log.Println("✅ 服务器已关闭")
}