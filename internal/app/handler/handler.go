package handler

import (
	"Backeven/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция в которой мы отдельно регистрируем маршруты,чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/Nalogimain", h.GetAllExpense)
	router.GET("/CostService/:id", h.GetExpenseByID)
	//router.GET("/breakevencalc", h.GetBreakeven)
}

// Получаем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorhandler для удобного вывода ошибок
func (h *Handler) errorhandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

/*
	func (h *Handler) GetAllServices(ctx *gin.Context) {
		var services []repository.Service
		var err error

		searchQuery := ctx.Query("BreakenevSearch") // получаем значение из поля поиска
		if searchQuery == "" {                      // если поле поиска пусто, то просто получаем из репозитория все записи
			services, err = h.Repository.GetCostsService()
			if err != nil {
				logrus.Error(err)
			}
		} else {
			services, err = h.Repository.GetServicesByTitle(searchQuery) // в ином случае ищем заказ по заголовку
			if err != nil {
				logrus.Error(err)
			}
		}

		ctx.HTML(http.StatusOK, "mainpageNalogi.html", gin.H{
			"time":            time.Now().Format("15:04:05"),
			"services":        services,
			"BreakenevSearch": searchQuery, // передаем введенный запрос обратно на страницу
			// в ином случае оно будет очищаться при нажатии на кнопку
		})
	}
*/
/*func (h *Handler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	service, err := h.Repository.GetService(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "CostProductInfo.html", gin.H{
		"service": service,
	})
}
*/
/*func (h *Handler) GetBreakeven(ctx *gin.Context) {
	CalcRequest := h.Repository.GetCalcServices()

	ctx.HTML(http.StatusOK, "breakevenCalc.html", gin.H{
		"CalcRequest": CalcRequest,
	})
}*/
