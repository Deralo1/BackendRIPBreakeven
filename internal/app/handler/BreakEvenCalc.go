package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBreakeven(ctx *gin.Context) {
	idStr := ctx.Param("BreakevenRequestID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.GetBreakevenCalcByID(int(id))
	if err != nil {
		// если заявка не найдена или удалена
		ctx.Redirect(http.StatusFound, "/Nalogimain")

		return
	}
	if len(calc.RequestExpense) == 0 {
		ctx.Redirect(http.StatusFound, "/Nalogimain")
		return
	}
	ctx.HTML(http.StatusOK, "breakevenCalc.html", gin.H{
		"CalcRequest": calc,
	})
}
func (h *Handler) DeleteBreakEvenCalc(ctx *gin.Context) {
	breakevenIDstr := ctx.PostForm("BreakevenRequestID")
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
	ctx.Redirect(http.StatusFound, "/Nalogimain")
}
