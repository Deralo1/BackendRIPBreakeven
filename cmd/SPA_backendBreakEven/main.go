package main

import (
	"fmt"

	"Backeven/internal/app/config"
	"Backeven/internal/app/dsn"
	"Backeven/internal/app/handler"
	"Backeven/internal/app/repository"
	"Backeven/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal("Failed to load config:%v", err)
	}
	//Инициализация minio клиента
	minioclient, err := config.NewMinioClient(cfg.Minio)
	if err != nil {
		logrus.Fatal("Failed to inizialize Minio client: %v", err)
	}
	hand := handler.NewHandler(rep, minioclient, cfg.Minio.BucketName)
	// инициалазиция обработчика
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
