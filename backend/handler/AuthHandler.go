package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/middleware"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[WARN] Erro na validação de registro: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro ao registrar usuário: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao registrar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registro criado"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[WARN] Erro na validação de login: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos", "details": err.Error()})
		return
	}

	admin, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro no login: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha incorretos"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(middleware.JWTSecret))
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if err != nil {
			log.Printf("[ERROR] Erro ao gerar token: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao processar login"})
		return
	}

	c.SetCookie(
		"auth_token",
		tokenString,
		int(time.Hour*24*30/time.Second),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login realizado com sucesso",
		"token":   tokenString,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",   //path
		"",    //domain
		false, //secure
		true,  //httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado com sucesso"})
}
