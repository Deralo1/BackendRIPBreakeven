package main

import (
	_ "Backeven/docs"
	"Backeven/internal/app/config"
	"Backeven/internal/app/dsn"
	"Backeven/internal/app/handler"
	"Backeven/internal/app/repository"
	"Backeven/internal/middleware"
	"Backeven/internal/pkg"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Загружаем конфиг
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// 2. Инициализируем Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: "password", // если есть
		DB:       conf.Redis.DB,
	})

	if _, err = rdb.Ping(context.Background()).Result(); err != nil {
		logrus.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	logrus.Info("Успешное подключение к Redis.")

	// 3. Инициализируем Gin
	router := gin.Default()

	// 4. Подключаем middleware гостевой сессии
	router.Use(middleware.GuestSessionMiddleware(rdb))

	// 5. Остальная инициализация
	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString, rdb)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	minioClient, err := config.NewMinioClient(conf.Minio)
	if err != nil {
		logrus.Fatal("Failed to initialize Minio client: ", err)
	}

	jwtDuration, err := time.ParseDuration(conf.JWT.ExpiresIn)
	if err != nil {
		logrus.Fatalf("Invalid JWT expiration duration in config: %v", err)
	}

	hand := handler.NewHandler(
		rep,
		minioClient,
		conf.Minio.BucketName,
		rdb,
		conf.JWT.SecretKey,
		conf.ServiceHost,
		jwtDuration,
	)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
