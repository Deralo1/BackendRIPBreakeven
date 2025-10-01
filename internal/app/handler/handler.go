package handler

import (
	"Backeven/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error

	searchQuery := ctx.Query("BreakenevSearch") // получаем значение из поля поиска
	if searchQuery == "" {                      // если поле поиска пусто, то просто получаем из репозитория все записи
		orders, err = h.Repository.GetCostsCards()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "mainpageNalogi.html", gin.H{
		"time":   time.Now().Format("15:04:05"),
		"orders": orders,
		"query":  searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}
func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "CostProductInfo.html", gin.H{
		"order": order,
	})
}

func (h *Handler) GetBreakeven(ctx *gin.Context) {
	orders, err := h.Repository.GetCostsCards()
	if err != nil {
		logrus.Error(err)
	}

	// фильтруем только необходимые карточки: Аренда склада и Сырье
	var filtered []repository.Order
	for _, o := range orders {
		if o.Title == "Аренда склада" || o.Title == "Сырье" {
			filtered = append(filtered, o)
		}
	}

	ctx.HTML(http.StatusOK, "breakevenCalc.html", gin.H{
		"orders": filtered,
	})
}
