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
	"github.com/cycling02/ai-novel-backend/internal/middleware"
	"github.com/cycling02/ai-novel-backend/internal/repository"
	"github.com/cycling02/ai-novel-backend/internal/service"
	"gorm.io/gorm"
)

// Server HTTP 服务器
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	config     *config.Config
	db         *database.PostgresDB
}

// NewServer 创建服务器
func NewServer(cfg *config.Config, db *database.PostgresDB, gormDB *gorm.DB) (*Server, error) {
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

	// ==================== 初始化仓库 ====================
	// 使用 PostgresDB 的仓库
	novelRepo := repository.NewNovelRepository(db)
	chapterRepo := repository.NewChapterRepository(db)

	// 使用 GORM 的仓库
	userRepo := repository.NewUserRepository(gormDB)
	characterRepo := repository.NewCharacterRepository(gormDB)
	worldSettingRepo := repository.NewWorldSettingRepository(gormDB)
	chapterVersionRepo := repository.NewChapterVersionRepository(gormDB)

	// ==================== 初始化服务 ====================
	novelService := service.NewNovelService(
		novelRepo,
		chapterRepo,
		nil,
		nil,
		nil,
		chapterChain,
		outlineChain,
		plotSuggestChain,
		editChain,
		orchestrator,
	)

	authService := service.NewAuthService(userRepo)
	characterService := service.NewCharacterService(characterRepo, novelRepo)
	worldSettingService := service.NewWorldSettingService(worldSettingRepo, novelRepo)
	versionService := service.NewVersionService(chapterVersionRepo, chapterRepo)
	exportService := service.NewExportService(novelRepo, chapterRepo)

	// ==================== 初始化处理器 ====================
	novelHandler := handler.NewNovelHandler(novelService)
	healthHandler := handler.NewHealthHandler(db.Ping)
	authHandler := handler.NewAuthHandler(authService)
	characterHandler := handler.NewCharacterHandler(characterService)
	worldSettingHandler := handler.NewWorldSettingHandler(worldSettingService)
	versionHandler := handler.NewVersionHandler(versionService)
	exportHandler := handler.NewExportHandler(exportService)
	streamHandler := handler.NewStreamHandler(novelService)

	// ==================== 路由配置 ====================
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Ready)
	router.GET("/health/live", healthHandler.Live)

	// API v1 路由
	api := router.Group("/api/v1")
	{
		// 认证路由（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(authService))
		{
			// 小说相关
			novels := protected.Group("/novels")
			{
				novels.GET("", novelHandler.ListNovels)
				novels.POST("", novelHandler.CreateNovel)
				novels.GET("/:id", novelHandler.GetNovel)
				novels.PUT("/:id", novelHandler.UpdateNovel)
				novels.DELETE("/:id", novelHandler.DeleteNovel)

				// AI 生成
				novels.POST("/:id/generate", novelHandler.GenerateChapter)
				novels.POST("/:id/outline", novelHandler.GenerateOutline)
				novels.POST("/:id/suggest", novelHandler.SuggestPlot)

				// 角色
				novels.POST("/:id/characters", characterHandler.CreateCharacter)
				novels.GET("/:id/characters", characterHandler.ListCharacters)
				novels.PUT("/:novelId/characters/:id", characterHandler.UpdateCharacter)
				novels.DELETE("/:novelId/characters/:id", characterHandler.DeleteCharacter)
				novels.GET("/:novelId/characters/:id", characterHandler.GetCharacter)

				// 世界观设定
				novels.POST("/:id/world-settings", worldSettingHandler.CreateWorldSetting)
				novels.GET("/:id/world-settings", worldSettingHandler.ListWorldSettings)
				novels.PUT("/:novelId/world-settings/:id", worldSettingHandler.UpdateWorldSetting)
				novels.DELETE("/:novelId/world-settings/:id", worldSettingHandler.DeleteWorldSetting)
				novels.GET("/:novelId/world-settings/:id", worldSettingHandler.GetWorldSetting)

				// 导出
				novels.POST("/:id/export", exportHandler.Export)
			}

			// 章节版本管理
			chapters := protected.Group("/chapters")
			{
				chapters.POST("/:chapterId/versions", versionHandler.CreateVersion)
				chapters.GET("/:chapterId/versions", versionHandler.ListVersions)
				chapters.GET("/:chapterId/versions/:versionNum", versionHandler.GetVersion)
				chapters.POST("/:chapterId/versions/:versionNum/rollback", versionHandler.Rollback)
			}

			// 内容编辑
			protected.POST("/content/edit", novelHandler.EditContent)

			// 流式生成
			protected.GET("/novels/:id/generate/stream", streamHandler.GenerateChapterStream)
			protected.POST("/novels/:id/generate/stream", streamHandler.GenerateChapterStreamPOST)
			protected.POST("/novels/:id/generate/batch", streamHandler.BatchGenerateStream)

			// 当前用户
			protected.GET("/auth/me", authHandler.GetCurrentUser)
		}
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		router: router,
		config: cfg,
		db:     db,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}