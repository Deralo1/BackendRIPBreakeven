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
	calc, err := h.Repository.GetBreakevenCalcByID(uint(id))
	if err != nil {
		// если заявка не найдена или удалена
		h.errorhandler(ctx, http.StatusNotFound, err)

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
func (h *Handler) AddExpenseToCalc(ctx *gin.Context) {
	userID := 1 // хардкод ид пользователя
	//Получаем ID из формы
	expenseIDstr := ctx.PostForm("ExpenseID")
	expenseID, err := strconv.Atoi(expenseIDstr)
	if err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}
	// Добавляем в калькулятор
	err = h.Repository.AddExpenseToCalc(userID, expenseID)
	if err != nil {
		if err.Error() == "услуга уже добавлена в калькулятор" {
			ctx.Redirect(http.StatusFound, "/Nalogimain") // здесь под вопросов куда редиректить
		} else {
			h.errorhandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Redirect(http.StatusFound, "/Nalogimain")
}
