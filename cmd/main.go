package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseyanez/dummies-app/internal/products/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	productHandler := handlers.NewProductHandler()

	r.Get("/products", productHandler.GetProducts)

	http.ListenAndServe(":8080", r)
}
