package handler

import (
	"Backeven/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllExpense(ctx *gin.Context) {
	var services []ds.Expense
	var err error
	UserID := 1
	searchQuery := ctx.Query("BreakenevSearch") // получаем значение из поля поиска
	if searchQuery == "" {                      // если поле поиска пусто, то просто получаем из репозитория все записи
		services, err = h.Repository.GetAllExpense()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		services, err = h.Repository.GetExpenseByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}
	calc, err := h.Repository.GetBreakevenCalc(UserID)
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "mainpageNalogi.html", gin.H{
		"time":            time.Now().Format("15:04:05"),
		"services":        services,
		"BreakenevSearch": searchQuery, // передаем введенный запрос обратно на страницу
		"calc_count":      h.Repository.GetExpensesInCalcCount(UserID),
		"CalcRequest":     calc,
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}
func (h *Handler) GetExpenseByID(ctx *gin.Context) {
	idStr := ctx.Param("ExpenseID") // получаем id заказа из урла
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	service, err := h.Repository.GetExpenseByID(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "CostProductInfo.html", gin.H{
		"service": service,
	})
}

/*func (h *Handler) GetBreakeven(ctx *gin.Context) {
	CalcRequest := h.Repository.GetCalcServices()

	ctx.HTML(http.StatusOK, "breakevenCalc.html", gin.H{
		"CalcRequest": CalcRequest,
	})
}*/
