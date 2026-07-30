package services

import "github.com/joseyanez/dummies-app/internal/auth/models"

type CreateAdminRequest struct {
	Name     string
	Password string
}

type LoginAdminRequest struct {
	Name     string
	Password string
}

type LoginResult struct {
	Errors  map[string]string
	Session *models.Session
}
