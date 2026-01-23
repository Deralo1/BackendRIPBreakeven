package handler

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/middleware"
	"Backeven/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Регистрация нового пользователя
// @Tags Домен пользователя
// @Accept json
// @Produce json
// @Param request body ds.ChangeUserDTO true "Данные для регистрации (логин и пароль)"
// @Success 204 "Успешная регистрация"
// @Failure 400 {object} handler.ErrorResponse  "Неверный формат данных"
// @Failure 500 {object} handler.ErrorResponse  "Ошибка сервера или пользователь с таким логином уже существует"
// @Router /user/register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input ds.ChangeUserDTO

	if err := ctx.BindJSON(&input); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	if input.Login == "" || input.Password == "" {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("введите логин и пароль"))
		return
	}

	err := h.Repository.RegisterUser(input)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary Аутентификация пользователя
// @Tags Домен пользователя
// @Accept  json
// @Produce json
// @Param user body ds.ChangeUserDTO true "Данные для входа"
// @Success 200 {object} ds.AuthResponseDTO "Успешный вход"
// @Failure 400 {object} handler.ErrorResponse "Неверный запрос"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Router /user/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	logrus.Warn("LOGIN: вход в LoginUser()")

	var input ds.ChangeUserDTO
	if err := ctx.BindJSON(&input); err != nil {
		logrus.Error("LOGIN: ошибка BindJSON:", err)
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Warnf("LOGIN: получены данные: login=%s", input.Login)

	if input.Login == "" || input.Password == "" {
		logrus.Error("LOGIN: логин или пароль пустые")
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("логин и пароль обязательны"))
		return
	}

	userDTO, err := h.Repository.LoginUser(input.Login, input.Password)
	if err != nil {
		logrus.Warn("LOGIN: неверный логин или пароль")
		h.errorhandler(ctx, http.StatusUnauthorized, err)
		return
	}

	logrus.Warnf("LOGIN: пользователь найден: userID=%d role=%s", userDTO.UserId, userDTO.Role)

	tokenString, expTime, err := service.GenerateJWT(
		userDTO.UserId,
		userDTO.Role,
		h.SecretKey,
		h.JWTDur,
	)
	if err != nil {
		logrus.Error("LOGIN: ошибка генерации токена:", err)
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка генерации токена: %w", err))
		return
	}

	logrus.Warnf("LOGIN: токен сгенерирован: %s", tokenString)
	logrus.Warnf("LOGIN: токен истекает в: %s", expTime)

	// ============================
	// 1. Записываем токен в Redis
	// ============================
	ttl := time.Until(expTime)
	logrus.Warnf("LOGIN: TTL токена = %v", ttl)

	logrus.Warnf("LOGIN: попытка записи токена в Redis: key=%s value=active ttl=%v", tokenString, ttl)
	if err := h.Repository.SaveToken(ctx, tokenString, ttl); err != nil {
		logrus.Error("LOGIN: ошибка записи токена в Redis:", err)
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка Redis"))
		return
	}

	logrus.Warn("LOGIN: токен успешно записан в Redis")

	// ============================
	// 2. Устанавливаем cookie
	// ============================
	logrus.Warnf("LOGIN: установка cookie session_token=%s", tokenString)

	ctx.SetCookie(
		"session_token",
		tokenString,
		int(ttl.Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	logrus.Warn("LOGIN: cookie успешно установлено")

	// ============================
	// 3. Возвращаем ответ
	// ============================
	logrus.Warn("LOGIN: отправка успешного ответа клиенту")

	h.successResponse(ctx, ds.AuthResponseDTO{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expTime.Unix(),
	})
}

// GetProfile
// @Summary Получить профиль
// @Tags Домен пользователя
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Success 200 {object} ds.UserDTO "Данные пользователя"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 404 {object} handler.ErrorResponse  "Пользователь не найден (редкий случай)"
// @Router /user/profile [get]
func (h *Handler) GetProfile(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	role := middleware.GetRole(ctx) // ← роль из JWT

	userdto, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}

	userdto.Role = role // ← вот это ключевой момент

	h.successResponse(ctx, userdto)
}

// UpdateUserProf
// @Summary Обновить профиль
// @Tags Домен пользователя
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param request body ds.ChangeUserDTO true "Новые данные (логин/пароль)"
// @Success 200 {object} ds.UserDTO "Успешное обновление"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат данных"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /user/profile [put]
func (h *Handler) UpdateUserProf(ctx *gin.Context) {

	userID := middleware.GetUserID(ctx)
	var userUpdate ds.ChangeUserDTO

	if err := ctx.BindJSON(&userUpdate); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	updateduserDto, err := h.Repository.UpdateUser(userID, userUpdate)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, updateduserDto)
}

// LogoutUser
// @Summary Выход из системы
// @Tags Домен пользователя
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Success 204 "Успешный выход"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован (отсутствует токен)"
// @Failure 500 {object} handler.ErrorResponse "Ошибка Redis/сервера"
// @Router /user/logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
	tokenString := service.ExtractToken(ctx)
	if tokenString == "" {
		h.errorhandler(ctx, http.StatusUnauthorized, fmt.Errorf("отсутствует токен для выхода"))
		return
	}

	claims, err := service.ParseJWT(tokenString, h.SecretKey)
	if err != nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	remainingDur := time.Until(claims.ExpiresAt.Time)

	if remainingDur > 0 {
		err := h.Repository.AddToBlacklist(ctx, tokenString, remainingDur)
		if err != nil {
			h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка при добавлении в блеклист: %w", err))
			return
		}
	} else {
		logrus.Warnf("Попытка выхода с просроченным токеном. TTL: %v", remainingDur)
	}

	ctx.SetCookie("session_token", "", -1, "/", h.HostName, false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Выход выполнен успешно"})
}
