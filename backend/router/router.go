package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
	"github.com/gabsfranca/mensagensAnonimasRH/database"
	"github.com/gabsfranca/mensagensAnonimasRH/handler"
	"github.com/gabsfranca/mensagensAnonimasRH/middleware"
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

	allowedOrigin := config.GetEnvVar("ALLOWED_ORIGIN")

	allowedOrigins := map[string]bool{
		allowedOrigin:          true,
		"http://192.168.0.183": true,
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

	// dbHost := config.GetEnvVar("DB_HOST")
	// dbUser := config.GetEnvVar("DB_USER")
	// dbPassword := config.GetEnvVar("DB_PASSWORD")
	// dbDatabase := config.GetEnvVar("DB_DATABASE")
	// if dbDatabase == "" {
	// 	dbDatabase = "denuncias"
	// }
	// dbPort := config.GetEnvVar("DB_PORT")

	//ESTANCIA BANCO DE DADOS
	db, err := database.NewPostgresConnection(
		config.GetEnvVar("DB_HOST"),
		config.GetEnvVar("DB_USER"),
		config.GetEnvVar("DB_PASSWORD"),
		config.GetEnvVar("DB_NAME"),
		config.GetEnvVar("DB_PORT"),
	)

	if err != nil {
		log.Fatal("Erro ao conectar no bacno: ", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Erro nas migrações", err)
	}

	//repare que um repo sempre depende de um db
	reportRepo := repo.NewGormReportRepo(db)
	mediaRepo := repo.NewGormMediaRepo(db)
	adminRepo := repo.NewGormAdminRepo(db)

	if reportRepo != nil && mediaRepo != nil && adminRepo != nil {
		log.Printf("repos ok")
	}

	//um service depende de um repo
	reportService := service.NewReportService(reportRepo)
	authService := service.NewAuthService(adminRepo)

	if reportService != nil && authService != nil {
		log.Printf("services ok")
	}

	//e um handler depende de um service
	reportHandler := handler.NewReportHandler(reportService)
	anonymousMessageHander := handler.NewAnonymousMessageHandler(reportRepo, mediaRepo, storageService)
	authHandler := handler.NewAuthHandler(authService)

	if reportHandler != nil && anonymousMessageHander != nil && authHandler != nil {
		log.Printf("handlers ok")
	}

	r.POST("/send-anonymous-message", func(c *gin.Context) {
		anonymousMessageHander.Handle(c)
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	adminGroup := r.Group("/messages")
	adminGroup.Use(middleware.RequireAuth())
	{
		// GET  /messages               -> lista todas
		adminGroup.GET("", reportHandler.GetAll) // corresponde a GET /messages
		// GET  /messages/:id           -> detalhes da mensagem
		adminGroup.GET("/:id", reportHandler.GetByID)
		// PATCH /messages/:id/status   -> atualiza status
		adminGroup.PATCH("/:id/status", reportHandler.PatchStatus)
		// POST  /messages/:id/obs      -> adiciona observação
		adminGroup.POST("/:id/obs", reportHandler.PostObs)
	}

	r.GET("/messages/:id/obs", reportHandler.GetObs)

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
