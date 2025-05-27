package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/middleware"
	"github.com/gabsfranca/mensagensAnonimasRH/service"
	"github.com/gin-gonic/gin"
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
		log.Printf("Erro na validação: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados invalidos"})
		return
	}

	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		log.Printf("erro ao registrar usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao registrar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registro criado"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("erro na validacao: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados invalidos", "details": err.Error()})
		return
	}

	admin, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("erro no login: %v", err)
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
		log.Printf("erro ao gerar token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao processar login"})
		return
	}

	c.SetCookie(
		"auth_token",
		tokenString,
		int(time.Hour*24*30/time.Second),
		"/",   //path
		"",    //domain
		false, //Secure
		true,  //HTTPonly
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
		-1, //max age negativo para expirar o cookie
		"/",
		"",
		false, //secure
		true,  //httponly
	)

	c.JSON(http.StatusOK, gin.H{"message": "logout realisado com sucesso"})
}
