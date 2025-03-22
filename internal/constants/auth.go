package constants

import "time"

const (
	TokenExp     = time.Hour * 3
	CookieMaxAge = 3600
	JwtSecret    = "your_secret_key"
)
