package handler

import (
	"Backeven/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllExpense(ctx *gin.Context) {
	var services []ds.ExpenseDTO
	var err error

	searchQuery := ctx.Query("BreakenevSearch") // получаем значение из поля поиска

	if searchQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		services, err = h.Repository.GetAllExpense()
	} else {
		services, err = h.Repository.GetExpenseByTitle(searchQuery) // в ином случае ищем заказ по заголовку
	}
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, services)
}

func (h *Handler) GetExpenseByID(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	service, err := h.Repository.GetExpenseByID(id)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successResponse(ctx, service)
}
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Expense sucessfully deleted",
	})
}

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
func (h *Handler) AddExpenseToCalc(ctx *gin.Context) {
	userID := h.GetCurrentUserId()
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
	ctx.JSON(200, gin.H{
		"message": "Added to Calc",
	})
}
