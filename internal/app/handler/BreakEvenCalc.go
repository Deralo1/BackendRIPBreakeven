package handler

import (
	"Backeven/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBreakevenCalcInfo(ctx *gin.Context) { // Get "иконки" корзины
	userID := h.GetCurrentUserId()

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
func (h *Handler) GetBreakEvenList(ctx *gin.Context) { // Get списка кроме удаленных и чероновых с фильтрацией по диапозону даты формирования и статусу
	startdatectx := ctx.Query("StartDate")
	enddatectx := ctx.Query("EndDate")
	statusctx := ctx.Query("Status")

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
	breakevenrequest, err := h.Repository.GetListCalcByDateAndStatus(statusctx, startdate, enddate)
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
	ctx.JSON(http.StatusOK, gin.H{})
}
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
	ctx.JSON(http.StatusOK, gin.H{})
}
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
	ctx.JSON(http.StatusOK, gin.H{})
}
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
	moderatorID := h.GetCurrentModeratorId()

	breakevenDTO, err := h.Repository.ProccessBreakEvenRequest(uint(id), moderatorID, action)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, breakevenDTO)
}
