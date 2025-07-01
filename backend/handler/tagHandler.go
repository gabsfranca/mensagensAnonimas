package handler

import (
	"net/http"

	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	Service service.TagService
}

func NewTagHandler(ts service.TagService) *TagHandler {
	return &TagHandler{Service: ts}
}

func (h *TagHandler) GetAvailableTags(c *gin.Context) {
	tags, err := h.Service.GetAvailableTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) RemoveTagFromMessage(c *gin.Context) {
	messageId := c.Param("messageId")
	tagId := c.Param("tagId")

	if err := h.Service.RemoveTagFromMessage(c.Request.Context(), messageId, tagId); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao remover tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "tag removida"})
}

func (h *TagHandler) GetTagStats(c *gin.Context) {
	stats, err := h.Service.CountReportsByTag(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar stats das tags"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
