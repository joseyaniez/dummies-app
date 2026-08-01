package middlewares

import (
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request after: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("Response before: %s %s\n", r.Method, r.URL.Path)
	})
}
