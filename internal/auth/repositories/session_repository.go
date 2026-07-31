package repositories

import (
	"database/sql"
	"fmt"

	"github.com/joseyanez/dummies-app/internal/auth/models"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

func (r *SessionRepository) SaveSession(session *models.Session) error {
	query := `
	  INSERT INTO sessions(admin_id, token, expires_at) VALUES (?, ?, ?)
	`
	result, err := r.DB.Exec(query, session.AdminID, session.Token, session.ExpiresAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	session.ID = fmt.Sprintf("%d", id)

	return nil
}
