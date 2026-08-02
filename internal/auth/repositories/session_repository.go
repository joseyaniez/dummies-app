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

func (r *SessionRepository) GetSessionById(sessionToken string) (*models.Session, error) {
	query := `
	  SELECT id, admin_id, token, expires_at FROM sessions WHERE token = ?
	`
	row := r.DB.QueryRow(query, sessionToken)

	session := &models.Session{}
	err := row.Scan(&session.ID, &session.AdminID, &session.Token, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) DeleteSessionsByAdminId(adminId string) error {
	query := `
	  DELETE FROM sessions WHERE admin_id = ?
	`
	_, err := r.DB.Exec(query, adminId)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) DeleteSessionsByToken(token string) error {
	query := `
	  DELETE FROM sessions WHERE token = ?
	`
	_, err := r.DB.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
