package pkg

import (
	"fmt"

	"Backeven/internal/app/config"
	"Backeven/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server starting up...")

	a.Handler.RegisterHandler(a.Router)
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ZeroTier IP + порт
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)

	logrus.Infof("Server listening on %s (HTTPS)", serverAddress)

	// Запуск HTTPS
	err := a.Router.RunTLS(
		serverAddress,
		"certs/10.205.157.61.pem",
		"certs/10.205.157.61-key.pem",
	)
	if err != nil {
		logrus.Fatal("Failed to start HTTPS server: ", err)
	}

	logrus.Info("Server shut down")
}
