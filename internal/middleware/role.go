package middleware

import (
	"net/http"

	"Backeven/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func RequireModerator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if GetRole(ctx) != ds.RoleModerator {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Доступ запрещен. Требуются права модератора."})
			return
		}
		ctx.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := GetRole(ctx)
		if role != ds.RoleCreator && role != ds.RoleModerator {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Требуется авторизация."})
			return
		}
		ctx.Next()
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{
			"https://10.205.157.61:3000",
			"https://10.205.157.61:4173",
		}

		origin := c.Request.Header.Get("Origin")

		// Проверяем, есть ли origin в списке
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
