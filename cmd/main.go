package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseyanez/dummies-app/internal/database"
	"github.com/joseyanez/dummies-app/internal/products/handlers"
)

func main() {
	db, err := database.OpenSQLiteDatabase("data/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established successfully.")

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	productHandler := handlers.NewProductHandler()

	r.Get("/admin/products/create", productHandler.CreateProduct)
	r.Get("/admin/products", productHandler.GetProducts)
	r.Post("/admin/products", productHandler.SaveProduct)

	http.ListenAndServe(":8080", r)
}
