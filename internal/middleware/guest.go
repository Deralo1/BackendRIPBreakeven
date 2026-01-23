package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GuestSessionMiddleware(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("guest_session")

		if err != nil || cookie == "" {
			// создаём новую сессию
			sessionID := uuid.New().String()
			now := time.Now().Unix()

			// сохраняем в Redis
			redis.HSet(ctx, "guest:"+sessionID, "created_at", now)
			redis.Expire(ctx, "guest:"+sessionID, 20*time.Minute)

			// отправляем куку
			ctx.SetCookie("guest_session", sessionID, 60*60*24*30, "/", "", true, true)

			ctx.Set("guest_session", sessionID)
		} else {
			// обновляем TTL
			redis.Expire(ctx, "guest:"+cookie, 20*time.Minute)
			ctx.Set("guest_session", cookie)
		}

		ctx.Next()
	}
}
