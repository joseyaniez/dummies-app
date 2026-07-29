package services

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
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
	// Crear el nuevo usuario
}
