package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gabsfranca/mensagensAnonimasRH/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

type AnonymousMessageHandler struct {
	ReportRepo repo.ReportRepo
	MediaRepo  repo.MediaRepo
	Storage    storage.StorageService
	Service    *service.AnonymousService
}

func NewAnonymousMessageHandler(rr repo.ReportRepo, mr repo.MediaRepo, storage storage.StorageService) *AnonymousMessageHandler {
	return &AnonymousMessageHandler{
		ReportRepo: rr,
		MediaRepo:  mr,
		Storage:    storage,
		Service:    service.NewAnonymousService(storage),
	}
}

func (h *AnonymousMessageHandler) Handle(c *gin.Context) {
	form, err := service.ParseAndValidateForm(c)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[WARN] Erro na validação do formulário: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	report := models.Report{
		Message:   form.Content,
		Status:    models.Recebido,
		CreatedAt: form.TimeStamp,
	}

	if err := h.ReportRepo.Create(c.Request.Context(), &report); err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro ao salvar denúncia: %v\nStacktrace:\n%s", err, stackErr.Stack())

		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Falha ao salvar denúncia"})
		return
	}

	mediaURLs := h.Service.SaveMediaFiles(form.Files)

	baseUrl := config.GetEnvVar("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://" + c.Request.Host
	}

	for _, media := range mediaURLs {
		m := models.Media{
			ReportId:  report.ID,
			URL:       media.URL,
			Type:      media.Type,
			CreatedAt: time.Now(),
		}

		if err := h.MediaRepo.Create(c.Request.Context(), &m); err != nil {
			stackErr := errors.Wrap(err, 0)
			if stackErr != nil {
				log.Printf("[ERROR] Falha ao salvar mídia: %v\nStacktrace:\n%s", err, stackErr.Stack())
			}
		}
	}

	updated, err := h.ReportRepo.FindByIdWithMedia(c.Request.Context(), report.ID)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro ao buscar denúncia após criação: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Erro ao buscar denúncia criada"})
		return
	}

	baseUrl = config.GetEnvVar("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://" + c.Request.Host
	}

	type MediaResponse struct {
		ID        string    `json:"id"`
		ReportId  string    `json:"reportId"`
		URL       string    `json:"url"`
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"createdAt"`
	}

	mediaResponses := make([]MediaResponse, len(updated.Media))
	for i, m := range updated.Media {
		mediaResponses[i] = MediaResponse{
			ID:        m.ID,
			ReportId:  m.ReportId,
			URL:       baseUrl + m.URL,
			Type:      string(m.Type),
			CreatedAt: m.CreatedAt,
		}
	}

	response := struct {
		ID        string          `json:"id"`
		Message   string          `json:"message"`
		Status    string          `json:"status"`
		CreatedAt time.Time       `json:"createdAt"`
		Media     []MediaResponse `json:"media"`
	}{
		ID:        updated.ID,
		Message:   updated.Message,
		Status:    string(updated.Status),
		CreatedAt: updated.CreatedAt,
		Media:     mediaResponses,
	}

	c.JSON(http.StatusOK, response)
}
