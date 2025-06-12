package models

import "github.com/golang-jwt/jwt/v4"

// User represents a user entity within the system.
type User struct {
	Username string // Username for identification.
	ID       int64  // Unique identifier for the user.

}

// Claims extends the standard JWT claims with a custom user ID field.
type Claims struct {
	jwt.RegisteredClaims       // Embedding the standard JWT claims.
	UserID               int64 `json:"userId"` // Custom claim carrying the user ID.

}
