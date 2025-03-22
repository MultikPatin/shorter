package app

import (
	"main/internal/interfaces"
	"net/http"
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
	//ctx := r.Context()

	//if r.Method != http.MethodGet {
	//	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//	return
	//}
	//
	//UserID, err := h.usersService.Login(ctx)
	//if err != nil {
	//	http.Error(w, "Origin not found", http.StatusNotFound)
	//	return
	//}
	//
	//w.Header().Set("content-type", jsonContentType)
	//w.Header().Set("Location", originLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
