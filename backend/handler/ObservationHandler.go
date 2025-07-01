package handler

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
	"github.com/gin-gonic/gin"
)

type ObservationHandler struct {
	Repo       repo.ObservationRepo
	ReportRepo repo.ReportRepo
}

func NewObservationHandler(obsRepo repo.ObservationRepo, reportRepo repo.ReportRepo) *ObservationHandler {
	return &ObservationHandler{
		Repo:       obsRepo,
		ReportRepo: reportRepo,
	}
}

func (h *ObservationHandler) PostUserObservation(c *gin.Context) {
	shortId := c.Param("id")

	report, err := h.Repo.FindByShortId(c.Request.Context(), shortId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "denuncia nao encontrada"})
		return
	}

	var body struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "conteúdo inválido"})
		return
	}

	obs := models.Observation{
		ReportID: report.ID,
		Content:  body.Content,
		Author:   "usuário",
	}

	if err := h.Repo.Create(c.Request.Context(), &obs); err != nil {
		log.Println("Erro ao salvar observação: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"succesS": false, "error": "Erro ao salvar observação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ObservationHandler) PostAdminObservation(c *gin.Context) {
	id := c.Param("id")

	// report, err := h.Repo.FindByReportId(c.Request.Context(), id)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "denuncia nao encontrada"})
	// 	return
	// }

	var body struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "conteúdo inválido"})
		return
	}

	obs := models.Observation{
		ReportID: id,
		Content:  body.Content,
		Author:   "admin",
	}

	if err := h.Repo.Create(c.Request.Context(), &obs); err != nil {
		log.Println("Erro ao salvar observação: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"succesS": false, "error": "Erro ao salvar observação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ObservationHandler) GetObservations(c *gin.Context) {
	paramId := c.Param("id")

	isUUID := isValidUUID(paramId)

	var reportId string

	if isUUID {
		reportId = paramId
	} else {
		report, err := h.Repo.FindByShortId(c.Request.Context(), paramId)
		if err != nil || report == nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "erro ao buscar observações"})
			return
		}
		reportId = report.ID
	}

	obss, err := h.Repo.FindByReportId(c.Request.Context(), reportId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "erro ao buscar observacoes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "observations": obss})
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return r.MatchString(uuid)
}
