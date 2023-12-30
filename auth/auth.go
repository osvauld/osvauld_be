package auth

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"math/big"
	"osvauld/config"
	"osvauld/infra/logger"
	"time"

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
	expirationTime := time.Now().Add(10 * time.Hour)
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

func VerifySignature(signatureStr string, publicKeyStr string, challengeStr string, userId uuid.UUID) (string, error) {
	// Decode the base64 encoded public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		logger.Errorf("Failed to decode base64 public key: %v", err)
	}

	// Parse the ECDSA public key
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		logger.Errorf("Failed to parse ECDSA public key: %v", err)
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Errorf("Public key is not of type *ecdsa.PublicKey")
	}

	// Decode the base64 encoded signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		logger.Errorf("Failed to decode base64 signature: %v", err)
	}

	// Assuming the signature is in ASN.1 DER format
	r := big.NewInt(0).SetBytes(signatureBytes[:len(signatureBytes)/2])
	s := big.NewInt(0).SetBytes(signatureBytes[len(signatureBytes)/2:])

	// Hash the challenge text
	hashed := sha256.Sum256([]byte(challengeStr))

	// Verify the signature
	valid := ecdsa.Verify(ecdsaPubKey, hashed[:], r, s)
	logger.Debugf("Signature verification result: %v", valid)
	if valid {
		return GenerateToken("test", userId)
	}
	return "", nil
}
