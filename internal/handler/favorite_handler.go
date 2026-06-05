package handler

import (
	"net/http"
	"strconv"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	svc service.FavoriteService
}

func NewFavoriteHandler(svc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

// GET /api/v1/favorites
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(int)

	projects, err := h.svc.GetFavorites(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить избранное", nil))
		return
	}

	c.JSON(http.StatusOK, projects)
}

// POST /api/v1/favorites/:project_id
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(int)

	projectIDStr := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Некорректный ID проекта", nil))
		return
	}

	if err := h.svc.AddFavorite(c.Request.Context(), userID, projectID); err != nil {
		if err.Error() == "PROJECT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, model.NewAPIError("NOT_FOUND", "Проект не найден", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось добавить в избранное", nil))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Проект добавлен в избранное"})
}

// DELETE /api/v1/favorites/:project_id
func (h *FavoriteHandler) DeleteFavorite(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(int)

	projectIDStr := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Некорректный ID проекта", nil))
		return
	}

	if err := h.svc.DeleteFavorite(c.Request.Context(), userID, projectID); err != nil {
		if err.Error() == "PROJECT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, model.NewAPIError("NOT_FOUND", "Проект не найден", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось удалить из избранного", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Проект удален из избранного"})
}