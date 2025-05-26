package models

import "github.com/golang-jwt/jwt/v4"

// User represents a user entity within the system.
type User struct {
	ID       int64  // Unique identifier for the user.
	Username string // Username for identification.
}

// Claims extends the standard JWT claims with a custom user ID field.
type Claims struct {
	UserID               int64 `json:"userId"` // Custom claim carrying the user ID.
	jwt.RegisteredClaims       // Embedding the standard JWT claims.
}
