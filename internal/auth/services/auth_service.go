package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/joseyanez/dummies-app/internal/auth/models"
	"github.com/joseyanez/dummies-app/internal/auth/repositories"
	"github.com/joseyanez/dummies-app/internal/auth/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository    repositories.AuthRepository
	sessionRepository repositories.SessionRepository
}

func NewAuthService(authRepo repositories.AuthRepository, sessionRepo repositories.SessionRepository) *AuthService {
	return &AuthService{
		authRepository:    authRepo,
		sessionRepository: sessionRepo,
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

func (s *AuthService) LoginUser(request LoginAdminRequest) (*LoginResult, error) {
	// Buscar usuario
	admin, err := s.authRepository.FindByName(request.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &LoginResult{
				Errors:  map[string]string{"form": "Usuario o contraseña incorrectos"},
				Session: nil,
			}, nil
		}
		return nil, err
	}
	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password))
	if err != nil {
		return &LoginResult{
			Errors:  map[string]string{"form": "Usuario o contraseña incorrectos"},
			Session: nil,
		}, nil
	}

	// Generar token
	token, err := util.GenerateToken()
	if err != nil {
		return nil, err
	}
	// Crear sesión con el token
	session := models.Session{
		AdminID:   admin.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	// Guardar sesión en la base de datos
	err = s.sessionRepository.SaveSession(&session)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Errors:  nil,
		Session: &session,
	}, nil
}

func (s *AuthService) LogoutUser(token string) error {
	err := s.sessionRepository.DeleteSessionsByToken(token)
	if err != nil {
		return err
	}
	return nil
}
