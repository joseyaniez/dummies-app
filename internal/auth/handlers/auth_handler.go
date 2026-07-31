package handlers

import (
	"log"
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

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	pages.LoginPage(nil, nil).Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error al parsear formulario: %v", err)
		errors := map[string]string{"form": "Hubo un error, intente de nuevo más tarde"}
		pages.LoginPage(nil, errors).Render(r.Context(), w)
		return
	}

	loginRequest := services.LoginAdminRequest{
		Name:     r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	loginResult, err := h.authService.LoginUser(loginRequest)
	if err != nil {
		log.Printf("Error en la autenticación: %v", err)
		errors := map[string]string{"form": "Hubo un error, intente de nuevo más tarde"}
		pages.LoginPage(nil, errors).Render(r.Context(), w)
		return
	}

	if len(loginResult.Errors) > 0 {
		values := map[string]string{"username": loginRequest.Name, "password": loginRequest.Password}
		pages.LoginPage(values, loginResult.Errors).Render(r.Context(), w)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    loginResult.Session.Token,
		Expires:  loginResult.Session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
