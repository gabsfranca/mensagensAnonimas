package router

import (
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
	"github.com/go-errors/errors"
)

func SetupRouter() *gin.Engine {
	log.Println("[INFO] Inicializando roteador...")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	allowedOrigin := config.GetEnvVar("URL")
	allowedOrigins := map[string]bool{
		allowedOrigin:          true,
		"http://192.168.0.183": true,
	}

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

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

	log.Println("[INFO] Inicializando serviço de armazenamento local...")
	storageService := storage.NewLocalStorage("./uploads")

	log.Println("[INFO] Conectando ao banco de dados...")
	db, err := database.NewPostgresConnection(
		config.GetEnvVar("DB_HOST"),
		config.GetEnvVar("DB_USER"),
		config.GetEnvVar("DB_PASSWORD"),
		config.GetEnvVar("DB_NAME"),
		config.GetEnvVar("DB_PORT"),
	)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Fatalf("[FATAL] Erro ao conectar no banco: %v\nStacktrace:\n%s", err, stackErr.Stack())

		}
	}

	log.Println("[INFO] Executando migrações do banco de dados...")
	if err := database.AutoMigrate(db); err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Fatalf("[FATAL] Erro nas migrações: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
	}

	log.Println("[INFO] Inicializando repositórios...")
	reportRepo := repo.NewGormReportRepo(db)
	mediaRepo := repo.NewGormMediaRepo(db)
	adminRepo := repo.NewGormAdminRepo(db)

	if reportRepo == nil || mediaRepo == nil || adminRepo == nil {
		log.Println("[ERROR] Falha ao instanciar repositórios")
	} else {
		log.Println("[INFO] Repositórios inicializados com sucesso")
	}

	log.Println("[INFO] Inicializando serviços...")
	reportService := service.NewReportService(reportRepo)
	authService := service.NewAuthService(adminRepo)

	if reportService == nil || authService == nil {
		log.Println("[ERROR] Falha ao instanciar serviços")
	} else {
		log.Println("[INFO] Serviços inicializados com sucesso")
	}

	log.Println("[INFO] Inicializando handlers...")
	reportHandler := handler.NewReportHandler(reportService)
	anonymousMessageHander := handler.NewAnonymousMessageHandler(reportRepo, mediaRepo, storageService)
	authHandler := handler.NewAuthHandler(authService)

	if reportHandler == nil || anonymousMessageHander == nil || authHandler == nil {
		log.Println("[ERROR] Falha ao instanciar handlers")
	} else {
		log.Println("[INFO] Handlers inicializados com sucesso")
	}

	// ===== ROTAS ESTÁTICAS PRIMEIRO =====
	r.Static("/media", "./uploads")
	r.Static("/foto", "./frontend")

	// ===== ROTAS DA API =====
	r.POST("/send-anonymous-message", func(c *gin.Context) {
		anonymousMessageHander.Handle(c)
	})

	// r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	adminGroup := r.Group("/messages")
	adminGroup.Use(middleware.RequireAuth())
	{
		adminGroup.GET("", reportHandler.GetAll)
		adminGroup.GET("/:id", reportHandler.GetByID)
		adminGroup.PATCH("/:id/status", reportHandler.PatchStatus)
		adminGroup.POST("/:id/obs", reportHandler.PostObs)
	}

	r.GET("/messages/:id/obs", reportHandler.GetObs)

	r.GET("/reports/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[INFO] Requisição recebida para /reports/%s/status", id)

		var report models.Report
		result := db.First(&report, "id = ?", id)

		if result.Error != nil {
			stackErr := errors.Wrap(result.Error, 0)
			log.Printf("[ERROR] Falha ao buscar denúncia %s: %v\nStacktrace:\n%s", id, result.Error, stackErr.Stack())
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "denúncia não encontrada",
			})
			return
		}

		log.Printf("[INFO] Denúncia encontrada: %s, status: %s", report.ID, report.Status)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  report.Status,
			"id":      report.ID,
		})
	})

	// ===== MIDDLEWARE DE FALLBACK POR ÚLTIMO =====
	r.NoRoute(func(c *gin.Context) {
		log.Printf("[WARN] Rota não encontrada: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "rota não encontrada",
		})
	})

	log.Println("[INFO] Todas as rotas registradas:")
	for _, route := range r.Routes() {
		log.Printf(" %s %s", route.Method, route.Path)
	}

	log.Println("[INFO] Roteador configurado com sucesso.")
	return r
}
