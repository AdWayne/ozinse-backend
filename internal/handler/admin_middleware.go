package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository" 
	"github.com/gin-gonic/gin"
)

func AdminMiddleware(adminRepo repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Пользователь не идентифицирован", nil))
			return
		}

		roleName, err := adminRepo.GetUserRoleName(c.Request.Context(), userID.(int))
		if err != nil || roleName != "Администратор" {
			c.AbortWithStatusJSON(http.StatusForbidden, model.NewAPIError("FORBIDDEN", "Доступ только для администраторов", nil))
			return
		}
		
		c.Next()
	}
}