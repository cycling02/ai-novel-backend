package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// WorldSettingHandler 世界观设定 HTTP 处理器
type WorldSettingHandler struct {
	worldSettingService *service.WorldSettingService
}

// NewWorldSettingHandler 创建世界观设定处理器
func NewWorldSettingHandler(worldSettingService *service.WorldSettingService) *WorldSettingHandler {
	return &WorldSettingHandler{worldSettingService: worldSettingService}
}

// CreateWorldSetting 创建设定
// POST /api/v1/novels/:id/world-settings
func (h *WorldSettingHandler) CreateWorldSetting(c *gin.Context) {
	var req model.WorldSetting
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novelID := c.Param("id")
	req.NovelID = novelID

	if err := h.worldSettingService.CreateWorldSetting(c.Request.Context(), &req); err != nil {
		if err == service.ErrNovelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// ListWorldSettings 获取设定列表
// GET /api/v1/novels/:id/world-settings
func (h *WorldSettingHandler) ListWorldSettings(c *gin.Context) {
	novelID := c.Param("id")

	settings, err := h.worldSettingService.ListWorldSettings(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if settings == nil {
		settings = []model.WorldSetting{}
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateWorldSetting 更新设定
// PUT /api/v1/novels/:novelId/world-settings/:id
func (h *WorldSettingHandler) UpdateWorldSetting(c *gin.Context) {
	var req model.WorldSetting
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("id")

	if err := h.worldSettingService.UpdateWorldSetting(c.Request.Context(), &req); err != nil {
		if err == service.ErrWorldSettingNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "设定不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteWorldSetting 删除设定
// DELETE /api/v1/novels/:novelId/world-settings/:id
func (h *WorldSettingHandler) DeleteWorldSetting(c *gin.Context) {
	id := c.Param("id")

	if err := h.worldSettingService.DeleteWorldSetting(c.Request.Context(), id); err != nil {
		if err == service.ErrWorldSettingNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "设定不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设定删除成功"})
}

// GetWorldSetting 获取设定详情
// GET /api/v1/novels/:novelId/world-settings/:id
func (h *WorldSettingHandler) GetWorldSetting(c *gin.Context) {
	id := c.Param("id")

	setting, err := h.worldSettingService.GetWorldSetting(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设定不存在"})
		return
	}

	c.JSON(http.StatusOK, setting)
}