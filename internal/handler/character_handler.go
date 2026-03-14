package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// CharacterHandler 角色 HTTP 处理器
type CharacterHandler struct {
	characterService *service.CharacterService
}

// NewCharacterHandler 创建角色处理器
func NewCharacterHandler(characterService *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

// CreateCharacter 创建角色
// POST /api/v1/novels/:id/characters
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req model.Character
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novelID := c.Param("id")
	req.NovelID = novelID

	if err := h.characterService.CreateCharacter(c.Request.Context(), &req); err != nil {
		if err == service.ErrNovelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// ListCharacters 获取角色列表
// GET /api/v1/novels/:id/characters
func (h *CharacterHandler) ListCharacters(c *gin.Context) {
	novelID := c.Param("id")

	characters, err := h.characterService.ListCharacters(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if characters == nil {
		characters = []model.Character{}
	}

	c.JSON(http.StatusOK, characters)
}

// UpdateCharacter 更新角色
// PUT /api/v1/novels/:novelId/characters/:id
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	var req model.Character
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("id")

	if err := h.characterService.UpdateCharacter(c.Request.Context(), &req); err != nil {
		if err == service.ErrCharacterNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteCharacter 删除角色
// DELETE /api/v1/novels/:novelId/characters/:id
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
	id := c.Param("id")

	if err := h.characterService.DeleteCharacter(c.Request.Context(), id); err != nil {
		if err == service.ErrCharacterNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色删除成功"})
}

// GetCharacter 获取角色详情
// GET /api/v1/novels/:novelId/characters/:id
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	id := c.Param("id")

	character, err := h.characterService.GetCharacter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, character)
}