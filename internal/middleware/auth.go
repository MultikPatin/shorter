package middleware

import (
	"context"
	"errors"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"net/http"
	"time"
)

var UserService interfaces.UsersService

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie == nil {
			userID, err := UserService.Login()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = setJWTCookie(w, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			tokenStr := cookie.Value
			claims, err := verifyJWT(tokenStr)
			if err != nil {
				w.Header().Set("content-type", constants.TextContentType)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), constants.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func verifyJWT(tokenStr string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(constants.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func setJWTCookie(w http.ResponseWriter, userID int64) error {
	tokenStr, err := generateJWT(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    tokenStr,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   constants.CookieMaxAge,
	}

	http.SetCookie(w, &cookie)
	return nil
}

func generateJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(constants.TokenExp)
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(constants.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
