package handler

import (
	"Backeven/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateExpenseInReq(ctx *gin.Context) {
	userid := h.GetCurrentUserId()

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

func (h *Handler) DeleteFromCalc(ctx *gin.Context) {
	userid := h.GetCurrentUserId()

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
