package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JwtKeyPath     string = "/app/keys/ec-pub.key"
	JwtAccessExp   uint   = 60 * 15          // 15 minutes
	JwtRefreshExp  uint   = 60 * 60 * 24 * 7 // 1 week
	CsrfAccessExp  uint   = 60 * 15          // 15 minutes
	CsrfRefreshExp uint   = 60 * 60 * 24 * 7 // 1 week
)

func readECPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM file")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func GenerateJWT(userID uint, tokenType string) (string, error) {
	privateKey, err := readECPrivateKey(JwtKeyPath)

	if err != nil {
		return "", err
	}

	var duration time.Duration

	if tokenType == "access" {
		duration = time.Duration(JwtAccessExp)
	} else if tokenType == "refresh" {
		duration = time.Duration(JwtRefreshExp)
	} else {
		return "", fmt.Errorf("incorect tokenType")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, claims)

	return token.SignedString(privateKey)
}

func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
