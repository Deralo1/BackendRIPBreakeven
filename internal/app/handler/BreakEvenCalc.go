package handler

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/middleware"
	"Backeven/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// GetBreakevenCalcInfo
// @Summary Получить информацию о черновике
// @Description Возвращает ID текущей черновой заявки и количество трат в ней.
// @Tags Домен заявки на подсчет точки
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Success 200 {object} object{BreakevenRequestID=int,expense_in_BreakEvenCount=int} "Информация о черновике"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /breakeven/calc [get]
func (h *Handler) GetBreakevenCalcInfo(ctx *gin.Context) { // Get "иконки" корзины
	tokenString := service.ExtractToken(ctx)

	var userID int = 0

	if tokenString != "" {
		val, redisErr := h.RedisClient.Get(context.Background(), tokenString).Result()

		if (redisErr == nil && val == "blacklist") || (redisErr != nil && !errors.Is(redisErr, redis.Nil)) {
			logrus.Warnf("Token is blacklisted or Redis error (Guest status): %v", redisErr)
		} else {
			claims, parseErr := service.ParseJWT(tokenString, h.SecretKey)

			if parseErr == nil {
				userID = claims.UserID
			}
		}
	}

	if userID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"BreakevenRequestID":        0,
			"expense_in_BreakEvenCount": 0,
		})
		return
	}

	CurrentCalcId, count, err := h.Repository.GetCalcInfo(userID)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	response := gin.H{
		"BreakevenRequestID":        CurrentCalcId,
		"expense_in_BreakEvenCount": count,
	}
	h.successResponse(ctx, response)
}

// GetBreakEvenList
// @Summary Получить список заявок
// @Description Для модератора - все заявки. Для создателя - только его заявки.
// @Tags Домен заявки на подсчет точки
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param status query string false "Фильтр по статусу ('черновик', 'сформирован', 'завершён', 'отклонён')"
// @Success 200 {array} []ds.BreakevenRequestDTO "Список заявок"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Router /breakeven [get]
func (h *Handler) GetBreakEvenList(ctx *gin.Context) { // Get списка кроме удаленных и черновых с фильтрацией по диапозону даты формирования и статусу
	startdatectx := ctx.Query("StartDate")
	enddatectx := ctx.Query("EndDate")
	statusctx := ctx.Query("Status")

	userID := middleware.GetUserID(ctx)
	userRole := middleware.GetRole(ctx)

	var startdate time.Time
	var enddate time.Time
	var err error
	if startdatectx != "" {
		startdate, err = time.Parse("2006-01-02", startdatectx)
		if err != nil {
			h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат даты. Используйте YYYY-MM-DD"))
			return
		}
	}
	if enddatectx != "" {
		enddate, err = time.Parse("2006-01-02", enddatectx)
		if err != nil {
			h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат даты. Используйте YYYY-MM-DD"))
			return
		}
	}
	breakevenrequest, err := h.Repository.GetListCalcByDateAndStatus(userID, userRole, statusctx, startdate, enddate)
	if err != nil {
		if err.Error() == "record not found" {
			h.successResponse(ctx, []ds.BreakevenRequestDTO{})
			return
		}
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, breakevenrequest)
}

// GetBreakeven
// @Summary Получить одну заявку по id
// @Tags Домен заявки на подсчет точки
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID заявки"
// @Success 200 {array} ds.BreakevenRequestDTO "Заявка"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Router /breakeven/{id} [get]
func (h *Handler) GetBreakeven(ctx *gin.Context) { // GET одна запись заявки
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	calcDTO, err := h.Repository.GetBreakevenCalcByID(int(id))
	if err != nil {
		if err.Error() == "record not found" {
			h.errorhandler(ctx, http.StatusNotFound, fmt.Errorf("заявка с ID %v не найдена", id))
			return
		}
		h.errorhandler(ctx, http.StatusInternalServerError, err)
	}
	h.successResponse(ctx, calcDTO)
}

// PutBreakEven
// @Summary Обновить черновик
// @Description Обновляет поле 'TextToAnalyse' черновой заявки.
// @Tags Домен заявки на подсчет точки
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID заявки (черновика)"
// @Param request body ds.UpdateRequestExpenseDTO true "Новые данные для подсчета"
// @Success 200 {object} ds.BreakevenRequestDTO "Успешное обновление"
// @Failure 400 {object} handler.ErrorResponse"Неверный формат ID или данных"
// @Failure 401 {object} handler.ErrorResponse"Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера/Не является черновиком"
// @Router /breakeven/{id} [put]
func (h *Handler) PutBreakEven(ctx *gin.Context) { // PUT изменение полей заявки по теме
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	var CalcUpdates ds.UpdateBreakEvenCalcDTO
	if err := ctx.BindJSON(&CalcUpdates); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateBreakEvenCalc(uint(id), CalcUpdates)

	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	updatedexpenseDTO, err := h.Repository.GetBreakevenCalcByID(int(id))
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("обновление успешно, но ошибка при получении данных для ответа: %w", err))
		return
	}

	h.successResponse(ctx, updatedexpenseDTO)
}

// FormBreakEvenCalc
// @Summary Отправить заявку на модерацию
// @Description Переводит статус черновика на 'на модерации'.
// @Tags Домен заявки на подсчет точки
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID заявки (черновика)"
// @Success 200 {object} ds.BreakevenRequestDTO "Успешная отправка"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат ID"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера/Не является черновиком"
// @Router /breakeven/{id}/form [put]
func (h *Handler) FormBreakEvenCalc(ctx *gin.Context) { // PUT сформировать создателем
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.FormByCreatorBreakEvenCalc(uint(id))
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	updatedexpenseDTO, err := h.Repository.GetBreakevenCalcByID(int(id))
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("обновление успешно, но ошибка при получении данных для ответа: %w", err))
		return
	}
	h.successResponse(ctx, updatedexpenseDTO)
}

// DeleteBreakEvenCalc
// @Summary Удалить черновик заявки
// @Description Удаляет заявку (только черновик).
// @Tags Домен заявки на подсчет точки
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID заявки (черновика)"
// @Success 204 "Успешное удаление"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат ID"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера/Не является черновиком"
// @Router /breakeven/{id} [delete]
func (h *Handler) DeleteBreakEvenCalc(ctx *gin.Context) { // DELETE заявки
	breakevenIDstr := ctx.Param("id")
	CalcID, err := strconv.Atoi(breakevenIDstr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	err = h.Repository.DeleteBreakEvenCalc(uint(CalcID))
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ProccessBreakEven godoc
// @Summary Завершить или отклонить заявку (Только для Модератора)
// @Tags Домен заявки на подсчет точки
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID заявки"
// @Param action query string true "Действие ('complete' или 'reject')"
// @Success 200 {object} ds.BreakevenRequestDTO "Обновленная заявка"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 403 {object} handler.ErrorResponse "Доступ запрещен (не модератор)"
// @Failure 404 {object} handler.ErrorResponse "Заявка не найдена"
// @Router /breakeven/{id}/process [put]
func (h *Handler) ProccessBreakEven(ctx *gin.Context) { // PUT какой-то сложный пут
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	action := ctx.Query("action")
	if action == "" {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("action не может быть пустым"))
		return
	}
	moderatorID := middleware.GetUserID(ctx)

	breakevenDTO, err := h.Repository.ProccessBreakEvenRequest(uint(id), moderatorID, action)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, breakevenDTO)
}
