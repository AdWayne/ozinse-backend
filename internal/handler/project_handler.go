package handler

import (
	"net/http"
	"strconv"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	svc service.ProjectService
}

func NewProjectHandler(svc service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// 1. GET /api/v1/projects — Список проектов с фильтрами
func (h *ProjectHandler) GetAll(c *gin.Context) {
	var filter model.ProjectFilter
	
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Некорректные параметры фильтрации", nil))
		return
	}

	projects, err := h.svc.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить список проектов", nil))
		return
	}

	c.JSON(http.StatusOK, projects)
}

// 2. GET /api/v1/projects/featured — Подборки для главного экрана
func (h *ProjectHandler) GetFeatured(c *gin.Context) {
	featured, err := h.svc.GetFeatured(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось загрузить подборки для главного экрана", nil))
		return
	}

	c.JSON(http.StatusOK, featured)
}

// 3. GET /api/v1/projects/:id — Детальная информация о проекте
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Некорректный ID проекта", nil))
		return
	}

	project, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "PROJECT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, model.NewAPIError("NOT_FOUND", "Проект не найден", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Ошибка при получении данных проекта", nil))
		return
	}

	c.JSON(http.StatusOK, project)
}

// 4. GET /api/v1/projects/:id/seasons — Сезоны и серии для сериала
func (h *ProjectHandler) GetSeasons(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Некорректный ID сериала", nil))
		return
	}

	seasons, err := h.svc.GetSeasons(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NewAPIError("NOT_FOUND", "Сезоны не найдены или этот проект является фильмом", nil))
		return
	}

	c.JSON(http.StatusOK, seasons)
}