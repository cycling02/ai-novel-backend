package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/config"
	"github.com/cycling02/ai-novel-backend/internal/database"
	"github.com/cycling02/ai-novel-backend/internal/eino/agents"
	"github.com/cycling02/ai-novel-backend/internal/eino/chains"
	"github.com/cycling02/ai-novel-backend/internal/eino/components"
	"github.com/cycling02/ai-novel-backend/internal/handler"
	"github.com/cycling02/ai-novel-backend/internal/repository"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// Server HTTP 服务器
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	config     *config.Config
	db         *database.PostgresDB
}

// NewServer 创建服务器
func NewServer(cfg *config.Config, db *database.PostgresDB) (*Server, error) {
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	// 初始化 Eino 组件
	einoComponents, err := components.InitComponents(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 Eino 组件失败：%w", err)
	}

	// 初始化 Eino Chains
	chapterChain, err := chains.NewChapterGenerateChain(einoComponents)
	if err != nil {
		return nil, fmt.Errorf("创建章节生成链失败：%w", err)
	}

	outlineChain, err := chains.NewOutlineExpandChain(einoComponents)
	if err != nil {
		return nil, fmt.Errorf("创建大纲链失败：%w", err)
	}

	plotSuggestChain, err := chains.NewPlotSuggestChain(einoComponents)
	if err != nil {
		return nil, fmt.Errorf("创建情节建议链失败：%w", err)
	}

	editChain, err := chains.NewContentEditChain(einoComponents)
	if err != nil {
		return nil, fmt.Errorf("创建编辑链失败：%w", err)
	}

	// 初始化多 Agent 编排器
	orchestrator := agents.NewMultiAgentOrchestrator(einoComponents)

	// 初始化仓库
	novelRepo := repository.NewNovelRepository(db)
	chapterRepo := repository.NewChapterRepository(db)

	// 初始化服务
	novelService := service.NewNovelService(
		novelRepo,
		chapterRepo,
		nil, // TODO: character repo
		nil, // TODO: world repo
		nil, // TODO: knowledge repo
		chapterChain,
		outlineChain,
		plotSuggestChain,
		editChain,
		orchestrator,
	)

	// 初始化处理器
	novelHandler := handler.NewNovelHandler(novelService)
	healthHandler := handler.NewHealthHandler(db.Ping)

	// 路由配置
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Ready)
	router.GET("/health/live", healthHandler.Live)

	// API v1 路由
	api := router.Group("/api/v1")
	{
		// 小说管理
		novels := api.Group("/novels")
		{
			novels.POST("", novelHandler.CreateNovel)
			novels.GET("", novelHandler.ListNovels)
			novels.GET("/:id", novelHandler.GetNovel)
			novels.PUT("/:id", novelHandler.UpdateNovel)
			novels.DELETE("/:id", novelHandler.DeleteNovel)

			// AI 创作功能
			novels.POST("/:id/generate", novelHandler.GenerateChapter)
			novels.POST("/:id/outline", novelHandler.GenerateOutline)
			novels.POST("/:id/suggest", novelHandler.SuggestPlot)
		}

		// 内容编辑
		api.POST("/content/edit", novelHandler.EditContent)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		config:     cfg,
		db:         db,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown 优雅关闭
func (s *Server) Shutdown(ctx context.Context) error {
	if s.db != nil {
		s.db.Close()
	}
	return s.httpServer.Shutdown(ctx)
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
