package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseyanez/dummies-app/internal/database"
	"github.com/joseyanez/dummies-app/internal/products/handlers"
	"github.com/joseyanez/dummies-app/internal/products/services"
)

func main() {
	db, err := database.OpenSQLiteDatabase("data/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established successfully.")

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	productService := services.NewProductService()
	productHandler := handlers.NewProductHandler(productService)

	r.Get("/admin/products/create", productHandler.CreateProduct)
	r.Get("/admin/products", productHandler.GetProducts)
	r.Post("/admin/products", productHandler.SaveProduct)

	http.ListenAndServe(":8080", r)
}
