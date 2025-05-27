package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWTSecret = "chave_secreta_muito_segura_para_ambiente_de_producao"

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		tokenCookie, err := c.Cookie("auth_token")

		if err == nil && tokenCookie != "" {
			tokenStr = tokenCookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenStr = parts[1]
				}
			}
		}

		if tokenStr == "" {
			log.Println("Token nao fornecido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "nao autorizado"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo de assinatura invalido: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil {
			log.Printf("erro ao avaliar token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			adminID, ok := claims["sub"]
			if !ok {
				log.Println("ID do admin não encontrado no token")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido"})
				c.Abort()
				return
			}

			c.Set("adminID", adminID)
			c.Set("adminEmail", claims["email"])

			c.Next()
		} else {
			log.Println("token invalido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "nao autorizado"})
			c.Abort()
			return
		}
	}
}

func GetCurrentAdmin(c *gin.Context) (uint, error) {
	adminID, exists := c.Get("adminID")
	if !exists {
		return 0, errors.New("admin nao autenticado")
	}

	var id uint
	switch v := adminID.(type) {
	case float64:
		id = uint(v)
	case int:
		id = uint(v)
	case uint:
		id = v
	default:
		return 0, errors.New("id do admin em formato inválido")
	}

	return id, nil
}
