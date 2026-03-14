package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// ExportHandler 导出 HTTP 处理器
type ExportHandler struct {
	exportService *service.ExportService
}

// NewExportHandler 创建导出处理器
func NewExportHandler(exportService *service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// Export 导出小说
// POST /api/v1/novels/:id/export
func (h *ExportHandler) Export(c *gin.Context) {
	novelID := c.Param("id")

	var req struct {
		Format    string `json:"format" binding:"required,oneof=txt epub markdown pdf"`
		StartChap int    `json:"start_chap"`
		EndChap   int    `json:"end_chap"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.exportService.Export(c.Request.Context(), &model.ExportRequest{
		NovelID:   novelID,
		Format:    req.Format,
		StartChap: req.StartChap,
		EndChap:   req.EndChap,
	})
	if err != nil {
		if err == service.ErrNovelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 根据格式设置响应头
	switch req.Format {
	case "txt", "markdown", "md":
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", result.FileName))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, result.Content)
	case "epub":
		c.JSON(http.StatusOK, gin.H{
			"message":  "EPUB 导出功能开发中",
			"fileName": result.FileName,
		})
	case "pdf":
		c.JSON(http.StatusOK, gin.H{
			"message":  "PDF 导出功能开发中",
			"fileName": result.FileName,
		})
	}
}