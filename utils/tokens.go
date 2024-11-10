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
	"github.com/google/uuid"
)

const (
	JwtKeyPath       string        = "/app/keys/ec.pem"
	JwtPublicKeyPath string        = "/app/keys/ec-pub.pem"
	JwtAccessExp     time.Duration = time.Minute * 15   // 15 minutes
	JwtRefreshExp    time.Duration = time.Hour * 24 * 7 // 1 week
	CsrfExp          time.Duration = time.Minute * 15   // 15 minutes
)

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
	privateKey, err := readECPrivateKey(JwtKeyPath)

	if err != nil {
		return "", err
	}

	var duration time.Duration

	if tokenType == "access" {
		duration = JwtAccessExp
	} else if tokenType == "refresh" {
		duration = JwtRefreshExp
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
	publicKey, err := readECPublicKey(JwtPublicKeyPath)

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

func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
