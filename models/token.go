package models

import (
	"errors"
	"litstore/api/config"
	"litstore/api/models/enums"
	"litstore/api/utils"
	"time"

	"github.com/google/uuid"
)

type ActionToken struct {
	Base
	UserID     *uuid.UUID       `gorm:"type:uuid;index"`
	TokenHash  string           `json:"-" gorm:"unique;not null"`
	ActionType enums.ActionType `json:"action_type" gorm:"not null"`
	ExpiresAt  time.Time        `json:"expires_at" gorm:"not null;index"`
	UsedAt     *time.Time       `gorm:"index" json:"used_at"`
}

func GenerateActionToken(userID uuid.UUID, actionType enums.ActionType) (*ActionToken, string, error) {
	token, err := utils.GenerateToken()

	if err != nil {
		return nil, "", err
	}

	secret, err := utils.ReadHMACSecret(config.HMACSecretPath)

	if err != nil {
		return nil, "", err
	}

	tokenHash := utils.ComputeHMACToken(secret, token)

	return &ActionToken{
		UserID:     &userID,
		ActionType: actionType,
		TokenHash:  string(tokenHash),
		ExpiresAt:  time.Now().Add(enums.ExpirationForAction(actionType)),
	}, token, nil
}

func VerifyActionToken(token *ActionToken) error {
	if time.Now().After(token.ExpiresAt) {
		return errors.New("token has expired")
	}

	if token.UsedAt != nil {
		return errors.New("token has already been used")
	}

	return nil
}
