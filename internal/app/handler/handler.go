package handler

import (
	"Backeven/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
	BucketName  string
}

func NewHandler(r *repository.Repository, minioClient *minio.Client, bucketName string) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: minioClient,
		BucketName:  bucketName,
	}
}

// Получаем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
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

func (h *Handler) successResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"data": data,
	})
}
func (h *Handler) GetCurrentUserId() int {
	return 1
}
func (h *Handler) GetCurrentModeratorId() int {
	return 2
}

// RegisterHandler Функция в которой мы отдельно регистрируем маршруты,чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("api/v1")
	{
		// Траты
		api.GET("/expenses", h.GetAllExpense)
		api.GET("/expenses/:id", h.GetExpenseByID)
		api.POST("/expenses", h.CreateExpense)
		api.PUT("/expenses/:id", h.UpdateExpense)
		api.DELETE("/expenses/:id", h.DeleteExpense)
		api.POST("/expenses/:id/image", h.UploadExpenseImage)
		api.POST("/expenses/add-to-calc/:id", h.AddExpenseToCalc)

		//для заявки(калькулятора точки безубыточности)
		api.GET("/breakeven/calc", h.GetBreakevenCalcInfo)
		api.GET("/breakeven", h.GetBreakEvenList)
		api.GET("/breakeven/:id", h.GetExpenseByID)
		api.PUT("/breakeven/:id", h.PutBreakEven)
		api.PUT("/breakeven/:id/form", h.FormBreakEvenCalc)
		api.PUT("/breakeven/:id/process", h.ProccessBreakEven)
		api.DELETE("breakeven/:id", h.DeleteBreakEvenCalc)
	}
}
