package handler

import (
	"net/http"
	"strings"
	"ozinse-backend/internal/config"
	"ozinse-backend/internal/model"
	"ozinse-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Отсутствует заголовок Authorization", nil))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewAPIError("INVALID_TOKEN", "Некорректный формат токена", nil))
			return
		}

		claims, err := jwt.ParseToken(parts[1], cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewAPIError("INVALID_TOKEN", "Токен недействителен или просрочен", nil))
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
