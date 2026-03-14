package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/service"
)

// StreamHandler 流式输出 HTTP 处理器
type StreamHandler struct {
	novelService *service.NovelService
}

// NewStreamHandler 创建流式处理器
func NewStreamHandler(novelService *service.NovelService) *StreamHandler {
	return &StreamHandler{novelService: novelService}
}

// GenerateChapterStream 流式生成章节
// GET /api/v1/novels/:id/generate/stream
func (h *StreamHandler) GenerateChapterStream(c *gin.Context) {
	novelID := c.Param("id")

	var req model.GenerateStreamRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.NovelID = novelID

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 创建 flush writer
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "流式输出不支持"})
		return
	}

	// 发送开始事件
	h.sendEvent(c, "start", "", map[string]string{"novel_id": novelID})
	flusher.Flush()

	// 模拟流式生成（实际需要接入 LLM 流式 API）
	chunks := []string{
		"【第",
		"一章】",
		" 开始的",
		" 章节",
		" 内容",
		"...",
	}

	for i, chunk := range chunks {
		h.sendEvent(c, "chunk", chunk, map[string]interface{}{
			"index":    i,
			"progress": float64(i+1) / float64(len(chunks)) * 100,
		})
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}

	// 发送完成事件
	fullContent := "这是完整的章节内容。" // 实际应该是累积的内容
	h.sendEvent(c, "done", fullContent, map[string]interface{}{
		"word_count": len(fullContent),
	})
	flusher.Flush()
}

// GenerateChapterStreamPOST POST 版本流式生成
// POST /api/v1/novels/:id/generate/stream
func (h *StreamHandler) GenerateChapterStreamPOST(c *gin.Context) {
	novelID := c.Param("id")

	var req model.GenerateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.NovelID = novelID

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "流式输出不支持"})
		return
	}

	// 发送开始事件
	h.sendEvent(c, "start", "", map[string]string{"novel_id": novelID})
	flusher.Flush()

	// TODO: 接入实际的 LLM 流式输出
	// 这里模拟流式输出
	content := "这是 AI 生成的流式章节内容。\n\n"
	words := []string{"在", "一个", "遥远", "的", "大陆", "上", "，", "存在", "着", "一个", "神秘", "的", "王国", "..."}

	for i, word := range words {
		h.sendEvent(c, "delta", word, map[string]interface{}{
			"index":     i,
			"is_first":  i == 0,
			"is_last":   i == len(words)-1,
			"timestamp": time.Now().Unix(),
		})
		flusher.Flush()
		time.Sleep(50 * time.Millisecond)
	}

	// 完成
	h.sendEvent(c, "done", content, map[string]interface{}{
		"word_count":    len(content),
		"duration_ms":   1500,
		"model":         "MiniMax-M2.5",
	})
	flusher.Flush()
}

// BatchGenerateStream 批量生成（流式）
// POST /api/v1/novels/:id/generate/batch
func (h *StreamHandler) BatchGenerateStream(c *gin.Context) {
	var req model.BatchGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置 SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "流式输出不支持"})
		return
	}

	// 发送开始
	h.sendEvent(c, "batch_start", "", map[string]interface{}{
		"count":      req.Count,
		"word_count": req.WordCount,
	})
	flusher.Flush()

	// 模拟批量生成
	for i := 0; i < req.Count; i++ {
		chapterNum := i + 1
		h.sendEvent(c, "chapter_start", "", map[string]interface{}{
			"chapter":    chapterNum,
			"title":      fmt.Sprintf("%s 第%d章", req.StartTitle, chapterNum),
		})
		flusher.Flush()

		// 模拟章节内容生成
		words := []string{"第", fmt.Sprintf("%d", chapterNum), "章", "内容", "生成", "中", "..."}
		for j, word := range words {
			h.sendEvent(c, "chunk", word, map[string]interface{}{
				"chapter": chapterNum,
				"index":   j,
			})
			flusher.Flush()
			time.Sleep(30 * time.Millisecond)
		}

		h.sendEvent(c, "chapter_done", "", map[string]interface{}{
			"chapter":   chapterNum,
			"word_count": 500, // 模拟字数
		})
		flusher.Flush()
	}

	// 批量完成
	h.sendEvent(c, "batch_done", "", map[string]interface{}{
		"total_chapters": req.Count,
		"total_words":    req.Count * req.WordCount,
	})
	flusher.Flush()
}

// sendEvent 发送 SSE 事件
func (h *StreamHandler) sendEvent(c *gin.Context, eventType, content string, meta interface{}) {
	event := model.StreamEvent{
		Type:    eventType,
		Content: content,
		Meta:    meta,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	c.Writer.Write([]byte("data: "))
	c.Writer.Write(data)
	c.Writer.Write([]byte("\n\n"))
}