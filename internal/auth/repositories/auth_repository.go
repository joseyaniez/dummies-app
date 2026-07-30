package repositories

import (
	"database/sql"

	"github.com/joseyanez/dummies-app/internal/auth/models"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) SaveAdmin(name, hashedPassword string) error {
	query := `
	  INSERT INTO admins(name, password) VALUES (?, ?)
	`
	_, err := r.DB.Exec(query, name, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) FindByName(name string) (*models.Admin, error) {
	query := `
	  SELECT id, name, password FROM admins WHERE name = ?
	`

	var admin models.Admin
	row := r.DB.QueryRow(query, name)
	err := row.Scan(
		&admin.Id,
		&admin.Name,
		&admin.Password,
	)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
