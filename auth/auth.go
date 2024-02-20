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
	"golang.org/x/crypto/bcrypt"
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

func VerifySignature(signatureStr string, publicKeyStr string, challengeStr string) (bool, error) {
	// Decode the base64 encoded public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		logger.Errorf("Failed to decode base64 public key: %v", err)
		return false, err
	}

	// Parse the ECDSA public key
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return false, err
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("public key is not of type *ecdsa.PublicKey")
	}

	// Decode the base64 encoded signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		logger.Errorf("Failed to decode base64 signature: %v", err)
		return false, err
	}
	if len(signatureBytes) != 64 {
		return false, errors.New("invalid signature length")
	}
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])
	// Assuming the signature is in ASN.1 DER format

	// Hash the challenge text
	hashed := sha256.Sum256([]byte(challengeStr))

	// Verify the signature
	valid := ecdsa.Verify(ecdsaPubKey, hashed[:], r, s)
	logger.Debugf("Signature verification result: %v", valid)

	return valid, nil
}

func HashPassword(password string) (string, error) {
	// The second argument is the cost of hashing, which determines how much time is needed to calculate the hash.
	// The higher the cost, the more secure the hash, but the longer it will take to generate.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain password to check if they match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
