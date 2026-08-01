package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joseyanez/dummies-app/internal/auth/repositories"
)

type AuthMIddleware struct {
	sessionRepository repositories.SessionRepository
}

func NewAuthMiddleware(sessionRepository repositories.SessionRepository) *AuthMIddleware {
	return &AuthMIddleware{
		sessionRepository: sessionRepository,
	}
}

func (m *AuthMIddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request after: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("Response before: %s %s\n", r.Method, r.URL.Path)
	})
}

func (m *AuthMIddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			log.Printf("Error retrieving session cookie: %v", err)
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		session, err := m.sessionRepository.GetSessionById(cookie.Value)
		if err != nil {
			log.Printf("Error retrieving session from repository: %v", err)
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			err := m.sessionRepository.DeleteSessionsByAdminId(session.AdminID)
			if err != nil {
				log.Printf("Error deleting expired session from repository: %v", err)
			}
			log.Printf("Session cookie has expired: %v", session.ExpiresAt)
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
