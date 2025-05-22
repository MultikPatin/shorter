package constants

import "time"

// Type for storing the key of a user's ID in request contexts
type userIDKey string

// TokenExp defines the expiration duration of a JWT token (3 hours).
// Used to set the expiration period for authentication tokens.
const TokenExp = time.Hour * 3

// CookieMaxAge sets the maximum lifetime of cookies (3600 seconds), which equals one hour.
const CookieMaxAge = 3600

// JwtSecret holds the secret key used for signing JWT tokens.
// It is important to keep this key secure and avoid exposing it.
const JwtSecret = "your_secret_key"

// UserIDKey represents a unique identifier key for users stored in HTTP request contexts.
const UserIDKey userIDKey = "UserID"
