package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabsfranca/mensagensAnonimasRH/database"
	"github.com/gabsfranca/mensagensAnonimasRH/handler"
	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gabsfranca/mensagensAnonimasRH/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	fmt.Println("rotas registradas: ")
	for _, route := range r.Routes() {
		fmt.Printf(" %s %s\n", route.Method, route.Path)
	}

	allowedOrigins := map[string]bool{
		"http://localhost:3000":   true,
		"http://172.23.96.1:3000": true,
	}

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// c.Writer.Header().Set("Access-Control-Allow-Origin", ["http://localhost:3000", "http://172.23.96.1:3000"]) // Em desenvolvimento
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	storageService := storage.NewLocalStorage("./uploads")

	db, err := database.NewPostgresConnection("localhost", "postgres", "senha", "denuncias", "5432")

	if err != nil {
		log.Fatal("Erro ao conectar no bacno: ", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Erro nas migrações", err)
	}

	reportRepo := repo.NewGormReportRepo(db)

	mediaRepo := repo.NewGormMediaRepo(db)

	reportService := service.NewReportService(reportRepo)

	reportHandler := handler.NewReportHandler(reportService)

	anonymousMessageHander := handler.NewAnonymousMessageHandler(reportRepo, mediaRepo, storageService)

	r.POST("/send-anonymous-message", func(c *gin.Context) {
		anonymousMessageHander.Handle(c)
	})

	grp := r.Group("/messages")
	{
		// GET  /messages               -> lista todas
		grp.GET("", reportHandler.GetAll) // corresponde a GET /messages
		// GET  /messages/:id           -> detalhes da mensagem
		grp.GET("/:id", reportHandler.GetByID)
		// PATCH /messages/:id/status   -> atualiza status
		grp.PATCH("/:id/status", reportHandler.PatchStatus)
		// POST  /messages/:id/obs      -> adiciona observação
		grp.POST("/:id/obs", reportHandler.PostObs)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Mensagem anônima",
		})
	})

	r.GET("/reports/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("req recebida para reports/%s/status\n", id)

		var report models.Report
		result := db.First(&report, "id = ?", id)

		if result.Error != nil {
			log.Printf("erro ao buscar dados de report %s: %v\n", id, result.Error)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "denuncia nao encontrada",
			})
			return
		}

		log.Printf("report econtrado: %s, status: %s\n", report.ID, report.Status)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  report.Status,
			"id":      report.ID,
		})

	})

	r.Static("/foto", "./frontend")

	r.NoRoute(func(c *gin.Context) {
		log.Printf("rota nao encontrada: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "rota nao encontrada",
		})
	})

	fmt.Println("todas as rotas registradas:")
	for _, route := range r.Routes() {
		fmt.Printf(" %s %s\n", route.Method, route.Path)
	}

	r.Static("/media", "./uploads")

	return r
}
