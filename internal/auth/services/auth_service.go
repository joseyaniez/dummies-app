package services

import (
	"github.com/joseyanez/dummies-app/internal/auth/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepo,
	}
}

func (s *AuthService) CreateNewUser(request *CreateAdminRequest) (map[string]string, error) {
	validationErrors := make(map[string]string)
	if len(request.Name) < 8 {
		validationErrors["name"] = "El nombre debe tener al menos 8 caracteres"
	}
	if len(request.Password) < 8 {
		validationErrors["password"] = "La contraseña debe tener al menos 8 caracteres"
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	// Hacer un hash de la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordString := string(hash)

	// Crear el nuevo usuario
	err = s.authRepository.SaveAdmin(request.Name, passwordString)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
