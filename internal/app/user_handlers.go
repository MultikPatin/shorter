package app

import (
	"github.com/golang-jwt/jwt/v4"
	"main/internal/constants"
	"main/internal/interfaces"
	"main/internal/models"
	"net/http"
	"time"
)

func NewUsersHandlers(s interfaces.UsersService) *UsersHandlers {
	return &UsersHandlers{
		usersService: s,
	}
}

type UsersHandlers struct {
	usersService interfaces.UsersService
}

func (h *UsersHandlers) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.usersService.Login(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenStr, err := generateJWT(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    tokenStr,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   constants.CookieMaxAge,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
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
