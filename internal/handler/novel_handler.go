package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// NovelHandler 小说 HTTP 处理器
type NovelHandler struct {
	novelService *service.NovelService
}

// NewNovelHandler 创建小说处理器
func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{novelService: novelService}
}

// CreateNovel 创建小说
// POST /api/v1/novels
func (h *NovelHandler) CreateNovel(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous" // TODO: 实现认证
	}

	novel, err := h.novelService.CreateNovel(c.Request.Context(), userID, req.Title, req.Description, req.Genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, novel)
}

// GetNovel 获取小说详情
// GET /api/v1/novels/:id
func (h *NovelHandler) GetNovel(c *gin.Context) {
	id := c.Param("id")
	novel, err := h.novelService.GetNovel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "小说不存在"})
		return
	}

	c.JSON(http.StatusOK, novel)
}

// ListNovels 获取小说列表
// GET /api/v1/novels
func (h *NovelHandler) ListNovels(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	novels, err := h.novelService.ListNovels(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"novels": novels, "total": len(novels)})
}

// UpdateNovel 更新小说
// PUT /api/v1/novels/:id
func (h *NovelHandler) UpdateNovel(c *gin.Context) {
	id := c.Param("id")
	var req model.Novel

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	if err := h.novelService.UpdateNovel(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteNovel 删除小说
// DELETE /api/v1/novels/:id
func (h *NovelHandler) DeleteNovel(c *gin.Context) {
	id := c.Param("id")
	if err := h.novelService.DeleteNovel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GenerateChapter AI 生成章节
// POST /api/v1/novels/:id/generate
func (h *NovelHandler) GenerateChapter(c *gin.Context) {
	novelID := c.Param("id")
	var req model.GenerationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.NovelID = novelID

	resp, err := h.novelService.GenerateChapter(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GenerateOutline AI 生成大纲
// POST /api/v1/novels/:id/outline
func (h *NovelHandler) GenerateOutline(c *gin.Context) {
	novelID := c.Param("id")
	var req struct {
		BriefOutline string `json:"brief_outline" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outline, err := h.novelService.GenerateOutline(c.Request.Context(), novelID, req.BriefOutline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"outline": outline})
}

// SuggestPlot AI 情节建议
// POST /api/v1/novels/:id/suggest
func (h *NovelHandler) SuggestPlot(c *gin.Context) {
	novelID := c.Param("id")

	suggestions, err := h.novelService.SuggestPlot(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// EditContent AI 润色编辑
// POST /api/v1/content/edit
func (h *NovelHandler) EditContent(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.novelService.EditContent(c.Request.Context(), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
