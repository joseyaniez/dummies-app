package handlers

import (
	"net/http"

	"github.com/joseyanez/dummies-app/internal/auth/services"
	"github.com/joseyanez/dummies-app/internal/auth/views/pages"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authServ services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authServ,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	pages.LoginPage(nil, nil).Render(r.Context(), w)
}
