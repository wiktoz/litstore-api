package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"litstore/api/config"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	Name    string
	Value   string
	ExpTime time.Duration
}

func readECPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read file")
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		fmt.Print("type:", block.Type)
		return nil, fmt.Errorf("failed to decode PEM file")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func readECPublicKey(filename string) (*ecdsa.PublicKey, error) {
	keyData, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read file")
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM(keyData)

	if err != nil {
		return nil, fmt.Errorf("failed to parse key")
	}

	return publicKey, nil
}

func GenerateJWT(userID string, tokenType string) (string, error) {
	privateKey, err := readECPrivateKey(config.JwtPrivateKeyPath)

	if err != nil {
		return "", err
	}

	var duration time.Duration

	if tokenType == "access" {
		duration = config.JwtAccessExpTime
	} else if tokenType == "refresh" {
		duration = config.JwtRefreshExpTime
	} else {
		return "", fmt.Errorf("incorect tokenType")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"jti": uuid.New().String(),
		"exp": time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, claims)

	return token.SignedString(privateKey)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	publicKey, err := readECPublicKey(config.JwtPublicKeyPath)

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse/validate token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err) // Token is expired or invalid
	}

	return token, nil
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func IsBlacklisted(c *gin.Context, rds *redis.Client, token string) (bool, error) {
	result, err := rds.Get(c, token).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return true, err
	}

	if result == "revoked" {
		return true, nil
	}

	return false, nil
}

func RevokeToken(c *gin.Context, rds *redis.Client, token Token) error {
	var err error

	token.Value, err = c.Cookie(token.Name)

	if err != nil {
		return err
	}

	err = rds.Set(c, token.Value, "revoked", token.ExpTime).Err()

	if err != nil {
		return err
	}

	c.SetCookie(token.Name, "", -1, "/", "localhost", true, true)

	return nil
}
