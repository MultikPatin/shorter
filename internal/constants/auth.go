package constants

import "time"

type userIDKey string

const (
	TokenExp               = time.Hour * 3
	CookieMaxAge           = 3600
	JwtSecret              = "your_secret_key"
	UserIDKey    userIDKey = "UserID"
)
