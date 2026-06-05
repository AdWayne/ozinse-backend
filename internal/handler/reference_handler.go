package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ReferenceHandler struct {
	svc service.ReferenceService
}

func NewReferenceHandler(svc service.ReferenceService) *ReferenceHandler {
	return &ReferenceHandler{svc: svc}
}

func (h *ReferenceHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить категории", nil))
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *ReferenceHandler) GetGenres(c *gin.Context) {
	genres, err := h.svc.GetGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить жанры", nil))
		return
	}
	c.JSON(http.StatusOK, genres)
}

func (h *ReferenceHandler) GetAgeRatings(c *gin.Context) {
	ratings, err := h.svc.GetAgeRatings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить возрастные рейтинги", nil))
		return
	}
	c.JSON(http.StatusOK, ratings)
}