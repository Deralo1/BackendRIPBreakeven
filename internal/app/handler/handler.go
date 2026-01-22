package handler

import (
	"Backeven/internal/app/repository"
	"Backeven/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
	BucketName  string
	RedisClient *redis.Client
	SecretKey   string
	HostName    string
	JWTDur      time.Duration
}

func NewHandler(r *repository.Repository, minioClient *minio.Client, bucketName string, rdb *redis.Client, secretKey, hostName string, jwtDur time.Duration) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: minioClient,
		BucketName:  bucketName,
		RedisClient: rdb,
		SecretKey:   secretKey,
		HostName:    hostName,
		JWTDur:      jwtDur,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Получаем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "./resources")
}

// errorhandler для удобного вывода ошибок
func (h *Handler) errorhandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"description": err.Error(),
	})
}

func (h *Handler) successResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"data": data,
	})
}

// RegisterHandler Функция в которой мы отдельно регистрируем маршруты,чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
	api := router.Group("api/v1")
	{
		// Траты
		api.GET("/expenses", h.GetAllExpense)
		api.GET("/expenses/:id", h.GetExpenseByID)

		//для заявки(калькулятора точки безубыточности)
		api.GET("/breakeven/calc", h.GetBreakevenCalcInfo)

		// пользователи
		api.POST("/user/register", h.RegisterUser)
		api.POST("/user/login", h.LoginUser)

	}

	auth := router.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(h.SecretKey, h.RedisClient))
	{
		auth.GET("/user/profile", h.GetProfile)
		auth.PUT("/user/profile", h.UpdateUserProf)
		auth.POST("/user/logout", h.LogoutUser)

		//для заявки(калькулятора точки безубыточности
		auth.GET("/breakeven", h.GetBreakEvenList)
		auth.GET("/breakeven/:id", h.GetBreakeven)
		auth.PUT("/breakeven/:id", h.PutBreakEven)
		auth.PUT("/breakeven/:id/form", h.FormBreakEvenCalc)
		auth.DELETE("breakeven/:id", h.DeleteBreakEvenCalc)

		auth.POST("/expenses/add-to-calc/:id", h.AddExpenseToCalc)

		// м-м
		auth.PUT("/expense-calc/:id", h.UpdateExpenseInReq)
		auth.DELETE("/expense-calc/:id", h.DeleteFromCalc)
	}

	moderator := auth.Group("")
	moderator.Use(middleware.RequireModerator())
	{
		moderator.POST("/expenses", h.CreateExpense)
		moderator.PUT("/expenses/:id", h.UpdateExpense)
		moderator.DELETE("/expenses/:id", h.DeleteExpense)
		moderator.POST("/expenses/:id/image", h.UploadExpenseImage)

		moderator.PUT("/breakeven/:id/process", h.ProccessBreakEven)

	}
}
