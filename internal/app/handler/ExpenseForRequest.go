package handler

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateExpenseInReq
// @Summary Обновить трату в заявке
// @Description Обновляет количество и тип для выбранной траты в текущем черновике заявки.
// @Tags Домен м-м
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security SessionCookie
// @Param id path int true "ID траты"
// @Param request body ds.UpdateRequestExpenseDTO true "Количество и тип"
// @Success 204 "Успешное обновление"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат ID или данных"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /expense-calc/{id} [put]
func (h *Handler) UpdateExpenseInReq(ctx *gin.Context) {
	userid := middleware.GetUserID(ctx)

	expenseidstr := ctx.Param("id")
	expenseid, err := strconv.Atoi(expenseidstr)

	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req ds.UpdateRequestExpenseDTO
	if err := ctx.BindJSON(&req); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateExpenseForReq(userid, expenseid, req.AmountService, req.TypeSpend)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, gin.H{})

}

// DeleteFromCalc
// @Summary Удалить трату из заявки
// @Description Удаляет трату из текущего черновика заявки.
// @Tags Домен м-м
// @Produce json
// @Param id path int true "ID траты для удаления"
// @Success 204 "Успешное удаление"
// @Failure 400 {object} handler.ErrorResponse "Неверный формат ID"
// @Failure 401 {object} handler.ErrorResponse "Неавторизован"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /expense-calc/{id} [delete]
func (h *Handler) DeleteFromCalc(ctx *gin.Context) {
	userid := middleware.GetUserID(ctx)

	expenseidstr := ctx.Param("id")
	expenseID, err := strconv.Atoi(expenseidstr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteFromCalc(expenseID, userid)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, gin.H{})
}
