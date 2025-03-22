package middleware

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	jwt "github.com/golang-jwt/jwt/v4"
//	"net/http"
//	"time"
//)
//
//const tokenExp = time.Hour * 3
//const jwtSecret = "your_secret_key"
//
//type User struct {
//	ID       int
//	Username string
//}
//
//type Claims struct {
//	UserID int `json:"userId"`
//	jwt.RegisteredClaims
//}
//
//func Authentication(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		cookie, err := r.Cookie("access_token")
//		if err != nil || cookie == nil {
//			http.Redirect(w, r, "/login", http.StatusSeeOther)
//			return
//		}
//
//		tokenStr := cookie.Value
//		claims, err := verifyJWT(tokenStr)
//		if err != nil {
//			http.Redirect(w, r, "/login", http.StatusSeeOther)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
//
//func verifyJWT(tokenStr string) (*Claims, error) {
//	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return []byte(jwtSecret), nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := token.Claims.(*Claims)
//	if !ok || !token.Valid {
//		return nil, errors.New("invalid token")
//	}
//
//	return claims, nil
//}
//
//func generateJWT(userID int) (string, error) {
//	expirationTime := time.Now().Add(1 * time.Hour)
//	claims := &Claims{
//		UserID: userID,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(jwtSecret)
//	if err != nil {
//		return "", err
//	}
//
//	return tokenString, nil
//}
//
//func setJWTCookie(w http.ResponseWriter, userID int) error {
//	tokenStr, err := generateJWT(userID)
//	if err != nil {
//		return err
//	}
//
//	cookie := http.Cookie{
//		Name:     "access_token",
//		Value:    tokenStr,
//		Path:     "/",
//		HttpOnly: true,
//		MaxAge:   3600, // Время жизни куки в секундах (здесь 1 час)
//	}
//
//	http.SetCookie(w, &cookie)
//	return nil
//}
