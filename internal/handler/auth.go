package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Ошибка валидации полей", map[string]interface{}{"error": err.Error()}))
		return
	}

	err := h.svc.Register(c.Request.Context(), input)
	if err != nil {
		// Обрабатываем несовпадение паролей
		if err.Error() == "PASSWORDS_DO_NOT_MATCH" {
			c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Пароли не совпадают", nil))
			return
		}
		// Обрабатываем уже существующий email
		if err.Error() == "USER_EXISTS" {
			c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Пользователь с таким Email уже зарегистрирован", nil))
			return
		}
		
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось завершить регистрацию", nil))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешно завершена"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Некорректный логин или пароль", nil))
		return
	}

	access, refresh, err := h.svc.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Неверный Email или пароль", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input model.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Токен обновления обязателен", nil))
		return
	}

	access, refresh, err := h.svc.RefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Невалидный или просроченный сессионный токен", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input model.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Email указан неверно", nil))
		return
	}

	_ = h.svc.ResetPassword(c.Request.Context(), input.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Инструкция по восстановлению отправлена на Email"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID") 
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Пользователь не идентифицирован", nil))
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Ошибка обработки сессии", nil))
		return
	}
	
	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NewAPIError("NOT_FOUND", "Профиль пользователя не найден", nil))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Пользователь не идентифицирован", nil))
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Ошибка обработки сессии", nil))
		return
	}

	var input model.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Ошибка изменения профиля", nil))
		return
	}

	if err := h.svc.UpdateProfile(c.Request.Context(), userID, input); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось обновить данные", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно изменен"})
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewAPIError("UNAUTHORIZED", "Пользователь не идентифицирован", nil))
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Ошибка обработки сессии", nil))
		return
	}

	var input model.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewAPIError("VALIDATION_ERROR", "Пароль слишком короткий (минимум 6 символов)", nil))
		return
	}

	if err := h.svc.UpdatePassword(c.Request.Context(), userID, input); err != nil {
		if err.Error() == "INCORRECT_OLD_PASSWORD" {
			c.JSON(http.StatusBadRequest, model.NewAPIError("BAD_REQUEST", "Текущий пароль введен неверно", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewAPIError("SERVER_ERROR", "Не удалось обновить пароль", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}