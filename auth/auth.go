package auth

import (
	"time"

	"osvauld/config"
	"osvauld/infra/logger"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	UserID   string `json:"token"`
}

// GenerateToken creates a JWT token for authenticated users.
func GenerateToken(username string, id uuid.UUID) (string, error) {
	jwtSecret := config.GetJWTSecret()
	expirationTime := time.Now().Add(10 * time.Hour) // Token expires after 1 hour
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		UserID: id.String(),
	}
	jwtSecretKey := []byte(jwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logger.Errorf(err.Error())
		return "", errors.Wrap(err, "GenerateToken failed")
	}
	return tokenString, nil
}
