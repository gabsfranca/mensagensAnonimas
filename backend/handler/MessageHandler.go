package handler

import (
	"log"
	"net/http"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

type ReportHandler struct {
	Service service.ReportService
}

func NewReportHandler(rs service.ReportService) *ReportHandler {
	return &ReportHandler{Service: rs}
}

func (h *ReportHandler) GetAll(c *gin.Context) {
	reports, err := h.Service.GetAll(c.Request.Context())
	if err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[ERROR] GetAll: %v\nStacktrace:\n%s", err, stack.Stack())

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar mensagens"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	report, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[WARN] GetByID ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())

		}
		c.JSON(http.StatusNotFound, gin.H{"error": "mensagem não encontrada"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetObs(c *gin.Context) {
	id := c.Param("id")
	obs, err := h.Service.GetObs(c.Request.Context(), id)
	if err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[WARN] GetObs ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())
		}
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Sem observacao",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"obs":     obs,
	})
}

func (h *ReportHandler) PatchStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status models.Status `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[WARN] PatchStatus invalid body ID=%s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "status inválido"})
		return
	}
	if err := h.Service.ChangeStatus(c.Request.Context(), id, body.Status); err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[ERROR] PatchStatus change error ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao atualizar status"})
		return
	}
	report, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[ERROR] PatchStatus get updated ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao recuperar registro atualizado"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) PostObs(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Obs string `json:"obs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[WARN] PostObs invalid body ID=%s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "obs obrigatória"})
		return
	}
	if err := h.Service.AddObs(c.Request.Context(), id, body.Obs); err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[ERROR] PostObs add error ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao adicionar obs"})
		return
	}
	report, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		stack := errors.Wrap(err, 0)
		if stack != nil {
			log.Printf("[ERROR] PostObs get updated ID=%s: %v\nStacktrace:\n%s", id, err, stack.Stack())

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao recuperar registro atualizado"})
		return
	}
	c.JSON(http.StatusOK, report)
}
