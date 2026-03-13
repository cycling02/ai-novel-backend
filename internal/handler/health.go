package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	startTime time.Time
	dbCheck   func() error
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(dbCheck func() error) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		dbCheck:   dbCheck,
	}
}

// Health 整体健康检查
// GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	status := "healthy"
	services := make(map[string]interface{})

	// 检查数据库
	if h.dbCheck != nil {
		if err := h.dbCheck(); err != nil {
			status = "unhealthy"
			services["database"] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			services["database"] = gin.H{"status": "healthy"}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(h.startTime).String(),
		"services":  services,
	})
}

// Ready 就绪检查
// GET /health/ready
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.dbCheck != nil {
		if err := h.dbCheck(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// Live 存活检查
// GET /health/live
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}
