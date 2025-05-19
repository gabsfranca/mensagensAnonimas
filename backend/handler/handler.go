package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gabsfranca/mensagensAnonimasRH/storage"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	report := models.Report{
		Message:   form.Content,
		Status:    models.Recebido,
		CreatedAt: form.TimeStamp,
	}

	if err := h.ReportRepo.Create(c.Request.Context(), &report); err != nil {
		log.Println("Erro ao salvar denúncia:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Falha ao salvar denúncia"})
		return
	}

	mediaURLs := h.Service.SaveMediaFiles(form.Files)

	for _, url := range mediaURLs {
		m := models.Media{
			ReportId:  report.ID,
			URL:       url,
			Type:      models.Image,
			CreatedAt: time.Now(),
		}

		if err := h.MediaRepo.Create(c.Request.Context(), &m); err != nil {
			log.Println("falha ao salvar midia: ", err)
		}
	}

	updated, _ := h.ReportRepo.FindByID(c.Request.Context(), report.ID)

	c.JSON(http.StatusOK, updated)
}
