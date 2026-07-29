package repositories

import "database/sql"

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
