package config

import (
	"time"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

const (
	JwtPrivateKeyPath string        = "/app/keys/ec.pem"
	JwtPublicKeyPath  string        = "/app/keys/ec-pub.pem"
	JwtAccessExpTime  time.Duration = time.Minute * 15   // 15 minutes
	JwtRefreshExpTime time.Duration = time.Hour * 24 * 7 // 1 week
	CsrfExpTime       time.Duration = time.Minute * 15   // 15 minutes

	JwtAccessName  string = "jwt_access_token"
	JwtRefreshName string = "jwt_refresh_token"
	CsrfName       string = "csrf_token"

	HMACSecretPath string = "/app/keys/hmac_secret.key"
)
