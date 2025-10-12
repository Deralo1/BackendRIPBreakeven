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
	router.GET("/CostService/:ExpenseID", h.GetExpenseByID)
	router.GET("/breakevencalc/:BreakevenRequestID", h.GetBreakeven)
	// POST
	router.POST("/breakevencalc/add-expense", h.AddExpenseToCalc)
	router.POST("/breakevencalc/delete-calc", h.DeleteBreakEvenCalc)
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
