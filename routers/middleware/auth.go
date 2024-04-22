package middleware

import (
	"errors"
	"net/http"
	"strings"

	"osvauld/auth"
	"osvauld/config"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with Bearer"})
			return
		}

		claims := &auth.Claims{}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(config.GetJWTSecret()), nil
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

		if err != nil {
			logger.Errorf(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			logger.Errorf("Parse UUID error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			return
		}

		// Token is valid, set the username in the context so it can be used by the route handler if needed
		c.Set("username", claims.Username)
		c.Set("userId", userID)

		c.Next() // Proceed to the route handler
	}
}
