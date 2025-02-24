package middleware

import (
	"net/http"

	"github.com/NeManMarty1/todo-list/internal/config"
	"github.com/NeManMarty1/todo-list/internal/logger"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Log.WithContext(ctx).Warn("Отсутствует токен в заголовке Authorization")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется токен"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			logger.Log.WithContext(ctx).WithError(err).Warn("Недействительный токен")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			logger.Log.WithContext(ctx).Warn("Неверный формат токена")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		userID := int((*claims)["user_id"].(float64))
		c.Set("userID", userID)
		logger.Log.WithContext(ctx).WithFields(logrus.Fields{"user_id": userID}).Info("Токен проверен")
		c.Next()
	}
}
