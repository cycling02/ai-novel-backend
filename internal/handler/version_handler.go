package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// VersionHandler 版本控制 HTTP 处理器
type VersionHandler struct {
	versionService *service.VersionService
}

// NewVersionHandler 创建版本处理器
func NewVersionHandler(versionService *service.VersionService) *VersionHandler {
	return &VersionHandler{versionService: versionService}
}

// CreateVersion 创建新版本
// POST /api/v1/chapters/:chapterId/versions
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	chapterID := c.Param("chapterId")

	var req struct {
		Content    string `json:"content" binding:"required"`
		Title      string `json:"title"`
		ChangeDesc string `json:"change_desc"`
		ChangeType string `json:"change_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	if userID == nil {
		userID = "manual"
	}

	version, err := h.versionService.CreateVersion(
		c.Request.Context(),
		chapterID,
		req.Content,
		req.Title,
		req.ChangeDesc,
		req.ChangeType,
		userID.(string),
	)
	if err != nil {
		if err == service.ErrChapterNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, version)
}

// ListVersions 获取版本列表
// GET /api/v1/chapters/:chapterId/versions
func (h *VersionHandler) ListVersions(c *gin.Context) {
	chapterID := c.Param("chapterId")

	versions, err := h.versionService.GetVersions(c.Request.Context(), chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if versions == nil {
		versions = []model.ChapterVersion{}
	}

	c.JSON(http.StatusOK, versions)
}

// GetVersion 获取指定版本
// GET /api/v1/chapters/:chapterId/versions/:versionNum
func (h *VersionHandler) GetVersion(c *gin.Context) {
	chapterID := c.Param("chapterId")
	var versionNum int
	if _, err := fmt.Sscanf(c.Param("versionNum"), "%d", &versionNum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本号"})
		return
	}

	version, err := h.versionService.GetVersion(c.Request.Context(), chapterID, versionNum)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "版本不存在"})
		return
	}

	c.JSON(http.StatusOK, version)
}

// Rollback 回滚到指定版本
// POST /api/v1/chapters/:chapterId/versions/:versionNum/rollback
func (h *VersionHandler) Rollback(c *gin.Context) {
	chapterID := c.Param("chapterId")
	var versionNum int
	if _, err := fmt.Sscanf(c.Param("versionNum"), "%d", &versionNum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本号"})
		return
	}

	userID, _ := c.Get("user_id")
	if userID == nil {
		userID = "manual"
	}

	version, err := h.versionService.Rollback(c.Request.Context(), chapterID, versionNum, userID.(string))
	if err != nil {
		if err == service.ErrVersionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "版本不存在"})
			return
		}
		if err == service.ErrChapterNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, version)
}