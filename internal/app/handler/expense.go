package handler

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/middleware"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetAllExpense
// @Summary Получить список трат
// @Description Возвращает список всех существующих трат. Доступен публично.
// @Tags Домен трат
// @Produce json
// @Param searchbyexpensename query string false "Поиск по названию траты (частичное совпадение)"
// @Success 200 {object} ds.ExpenseDTO "Список трат"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /expenses [get]
func (h *Handler) GetAllExpense(ctx *gin.Context) {
	searchQuery := ctx.Query("BreakenevSearch")
	filterRecent := ctx.Query("recent") == "true"

	// 1. Получаем обычный список услуг
	var services []ds.ExpenseDTO
	var err error

	if searchQuery == "" {
		services, err = h.Repository.GetAllExpense()
	} else {
		services, err = h.Repository.GetExpenseByTitle(searchQuery)
	}

	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// 2. Если нужен фильтр по недавно просмотренным
	if filterRecent {
		sessionID := ctx.GetString("guest_session")

		// получаем ID просмотренных услуг
		viewedIDs, _ := h.Repository.GetRecentlyViewedIDs(ctx.Request.Context(), sessionID)

		if len(viewedIDs) == 0 {
			h.successResponse(ctx, []ds.ExpenseDTO{})
			return
		}

		// превращаем в set для быстрого поиска
		viewedSet := make(map[int]bool)
		for _, id := range viewedIDs {
			viewedSet[id] = true
		}

		// фильтруем список
		filtered := make([]ds.ExpenseDTO, 0)
		for _, s := range services {
			if viewedSet[s.ExpenseID] {
				filtered = append(filtered, s)
			}
		}

		h.successResponse(ctx, filtered)
		return
	}

	// 3. Если фильтр не включён — возвращаем обычный список
	h.successResponse(ctx, services)
}

func (h *Handler) AddViewedExpense(ctx *gin.Context) {
	sessionID := ctx.GetString("guest_session")
	expenseID := ctx.Param("id")

	key := "guest:" + sessionID + ":viewed"

	// удаляем дубликаты
	h.RedisClient.LRem(ctx, key, 0, expenseID)

	// добавляем в начало
	h.RedisClient.LPush(ctx, key, expenseID)

	// храним только последние 10
	h.RedisClient.LTrim(ctx, key, 0, 9)

	// обновляем TTL
	h.RedisClient.Expire(ctx, key, 20*time.Minute)
}

// GetExpenseByID
// @Summary Получить трату по ID
// @Description Возвращает информацию о конкретном трате. Доступен публично.
// @Tags Домен трат
// @Produce json
// @Param id path int true "ID траты"
// @Success 200 {object} ds.ExpenseDTO "Информация о трате"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} handler.ErrorResponse "Трата не найден"
// @Router /expenses/{id} [get]
func (h *Handler) GetExpenseByID(ctx *gin.Context) {
	h.AddViewedExpense(ctx)
	idStr := ctx.Param("id") // получаем id заказа из урла
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	service, err := h.Repository.GetExpenseByID(id)
	if err != nil {
		h.errorhandler(ctx, http.StatusNotFound, err)
		return
	}

	h.successResponse(ctx, service)
}

// CreateExpense
// @Summary Создать новую трату
// @Description Создает новую трату. Требуются права **Модератора**.
// @Tags Домен трат
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param request body ds.UpdateexpenseDTO true "Цена и описание траты"
// @Success 201 {object} ds.ExpenseDTO "Успешное создание траты"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат данных"
// @Failure 403 {object} handler.ErrorResponse "Доступ запрещен (не модератор)"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /expenses [post]
func (h *Handler) CreateExpense(ctx *gin.Context) {
	var expenseDTO ds.ExpenseDTO
	if err := ctx.BindJSON(&expenseDTO); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	CreatedExpense, err := h.Repository.CreateExpense(expenseDTO)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(201, gin.H{
		"data": CreatedExpense,
	})
}

// UpdateExpense
// @Summary Обновить трату
// @Description Обновляет цену и описание траты по ID. Требуются права **Модератора**.
// @Tags Домен трат
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID траты"
// @Param request body ds.UpdateexpenseDTO true "Новая цены и описание траты"
// @Success 200 {object} ds.ExpenseDTO "Успешное обновление"
// @Failure 400 {object} handler.ErrorResponse  "Неверный формат ID или данных"
// @Failure 403 {object} handler.ErrorResponse  "Доступ запрещен (не модератор)"
// @Failure 404 {object} handler.ErrorResponse  "Трата не найден"
// @Router /expenses/{id} [put]
func (h *Handler) UpdateExpense(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	var expenseUpdates ds.UpdateexpenseDTO
	if err := ctx.BindJSON(&expenseUpdates); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	expenseUpdate, err := h.Repository.Updateexpense(id, expenseUpdates)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, expenseUpdate)
}

// DeleteExpense
// @Summary Удалить Трату
// @Description Устанавливает флаг is_deleted = true для траты. Требуются права **Модератора**.
// @Tags Домен трат
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID траты"
// @Success 204 "Успешное удаление"
// @Failure 400 {object} handler.ErrorResponse"Неверный формат ID"
// @Failure 403 {object} handler.ErrorResponse "Доступ запрещен (не модератор)"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /expenses/{id} [delete]
func (h *Handler) DeleteExpense(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	expenseDTO, err := h.Repository.GetExpenseByID(id)
	if err != nil {
		h.errorhandler(ctx, http.StatusNotFound, err)
		return
	}
	if expenseDTO.ImageURL != "" && h.MinioClient != nil {
		err = h.deleteImageFromMinio(expenseDTO.ImageURL)
		if err != nil {
			logrus.Errorf("Failed to delete image from minio %d, %v", id, err)
		}
	}
	err = h.Repository.DeleteExpense(id)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// UploadExpenseImage
// @Summary Загрузить изображение траты
// @Description Загружает и обновляет изображение для траты по ID. Требуются права **Модератора**.
// @Tags Домен трат
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID траты"
// @Param file formData file true "Файл изображения"
// @Success 200 {object} ds.ExpenseDTO "Успешная загрузка, возвращает обновленный трату"
// @Failure 400 {object} handler.ErrorResponse "Ошибка загрузки/формата файла"
// @Failure 403 {object} handler.ErrorResponse "Доступ запрещен (не модератор)"
// @Failure 500 {object} handler.ErrorResponse "Ошибка Minio/сервера"
// @Router /expenses/{id}/image [post]
func (h *Handler) UploadExpenseImage(ctx *gin.Context) {
	if h.MinioClient == nil {
		h.errorhandler(ctx, http.StatusServiceUnavailable, fmt.Errorf("image storage service not configured"))
		return
	}

	if !h.checkMinioConnection() {
		h.errorhandler(ctx, http.StatusServiceUnavailable, fmt.Errorf("image storage service termporarity unvailable. Please try again later"))
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("expense id not found: %v", err))
		return
	}

	expenseDTO, err := h.Repository.GetExpenseByID(id)
	if err != nil {
		h.errorhandler(ctx, http.StatusNotFound, fmt.Errorf("expense not found: %v", err))
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("image file is required: %v", err))
		return
	}

	if !h.isValidImage(file) {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid format: allowed JPG, PNG, Gif, WebP"))
		return
	}

	if file.Size > 5*1024*1024 {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("image size too large. Maximum 5MB allowed"))
		return
	}

	if expenseDTO.ImageURL != "" {
		err = h.deleteImageFromMinio(expenseDTO.ImageURL)
		if err != nil {
			logrus.Warnf("Failed to delete old image frim Minio : %v", err)
		}
	}
	objectname := h.generateImageName(file.Filename, id)

	imageurl, err := h.uploadImageToMinio(file, objectname)

	if err != nil {
		logrus.Errorf("Failed to upload image to Minio: %v", err)
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to upload image: %v", err))
		return
	}
	updatedExpenseDTO, err := h.Repository.UpdateExpenseImage(id, imageurl)

	if err != nil {
		logrus.Errorf("Failed to update genre in database, rolling back image")
		h.deleteImageFromMinio(imageurl)
		h.errorhandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to update expense: %v", err))
		return
	}
	logrus.Infof("Succesfully updated expense %d with image %v", id, imageurl)
	h.successResponse(ctx, updatedExpenseDTO)
}

// AddExpenseToCalc (в handler/expense.go, но относится к Заявкам)
// @Summary Добавить трату в черновик заявки
// @Description Добавляет трату в текущую черновую заявку пользователя. Требуется **Авторизация**.
// @Tags Домен трат
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID траты для добавления"
// @Success 204 "Успешное добавление"
// @Failure 400 {object} handler.ErrorResponse"Неверный формат ID"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера (например, трата уже добавлен)"
// @Router /expenses/add-to-calc/{id} [post]
func (h *Handler) AddExpenseToCalc(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	//Получаем ID из формы
	expenseIDstr := ctx.Param("id")
	expenseID, err := strconv.Atoi(expenseIDstr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	// Добавляем в калькулятор
	err = h.Repository.AddExpenseToCalc(userID, expenseID)
	if err != nil {
		if err.Error() == "услуга уже добавлена в калькулятор" {
			h.errorhandler(ctx, http.StatusConflict, err)
		} else {
			h.errorhandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}
