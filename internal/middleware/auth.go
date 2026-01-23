package middleware

import (
	"context"
	"errors"
	"net/http"

	"Backeven/internal/app/ds" // Импорт ds для доступа к UserRole
	"Backeven/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(secretKey string, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// ============================
		// 1. ЛОГИРУЕМ ВСЕ COOKIE
		// ============================
		cookies := ctx.Request.Cookies()
		for _, c := range cookies {
			logrus.Warnf("COOKIE RECEIVED: %s = %s", c.Name, c.Value)
		}

		// ============================
		// 2. ЛОГИРУЕМ Authorization
		// ============================
		authHeader := ctx.GetHeader("Authorization")
		logrus.Warnf("AUTH HEADER: %s", authHeader)

		// ============================
		// 3. Пробуем прочитать session_token напрямую
		// ============================
		sessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			logrus.Warn("NO session_token cookie")
		} else {
			logrus.Warnf("session_token cookie FOUND: %s", sessionToken)
		}

		// ============================
		// 4. Извлекаем токен через ExtractToken
		// ============================
		tokenString := service.ExtractToken(ctx)
		logrus.Warnf("ExtractToken() returned: %s", tokenString)

		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Неавторизован. Отсутствует токен.",
			})
			return
		}

		// ============================
		// 5. Проверяем Redis
		// ============================
		logrus.Warnf("REDIS CHECK: GET %s", tokenString)
		val, err := rdb.Get(context.Background(), tokenString).Result()

		if err != nil {
			logrus.Warnf("REDIS ERROR: %v", err)
		} else {
			logrus.Warnf("REDIS VALUE: %s", val)
		}

		if err == nil && val == "blacklist" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Токен недействителен (выход из системы).",
			})
			return
		}

		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.Error("Redis Error in AuthMiddleware:", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// ============================
		// 6. Парсим JWT
		// ============================
		claims, err := service.ParseJWT(tokenString, secretKey)
		if err != nil {
			logrus.Warnf("JWT PARSE ERROR: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Недействительный токен.",
			})
			return
		}

		// ============================
		// 7. Устанавливаем данные в контекст
		// ============================
		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.Role)
		ctx.Set("tokenString", tokenString)

		logrus.Warnf("AUTH OK: userID=%d role=%s", claims.UserID, claims.Role)

		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) int {
	userID, ok := ctx.Get("userID")
	if !ok {
		return 0
	}
	return userID.(int)
}

func IsModerator(ctx *gin.Context) bool {
	isModerator, ok := ctx.Get("isModerator")
	if !ok {
		return false
	}
	return isModerator.(bool)
}

func GetTokenString(ctx *gin.Context) string {
	token, ok := ctx.Get("tokenString")
	if !ok {
		return ""
	}
	return token.(string)
}

func GetRole(ctx *gin.Context) ds.UserRole {
	role, ok := ctx.Get("userRole")
	if !ok {
		return ds.RoleGuest
	}
	return role.(ds.UserRole)
}
