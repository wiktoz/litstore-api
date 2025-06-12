package enums

import "time"

type ActionType string

const (
	PasswordReset     ActionType = "password_reset"
	EmailVerification ActionType = "email_verification"
	AdminAction       ActionType = "admin_action"
)

var ActionExpirations = map[ActionType]time.Duration{
	PasswordReset:     30 * time.Minute,
	EmailVerification: 24 * time.Hour,
	AdminAction:       10 * time.Minute,
}

func ExpirationForAction(a ActionType) time.Duration {
	if exp, ok := ActionExpirations[a]; ok {
		return exp
	}

	return 15 * time.Minute
}
