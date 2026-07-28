package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseyanez/dummies-app/internal/database"
	"github.com/joseyanez/dummies-app/internal/products/handlers"
	"github.com/joseyanez/dummies-app/internal/products/repositories"
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

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	r.Get("/admin/products/create", productHandler.CreateProduct)
	r.Get("/admin/products", productHandler.GetProducts)
	r.Post("/admin/products", productHandler.SaveProduct)
	r.Get("/admin/products/{id}", productHandler.ViewProduct)
	r.Get("/admin/products/edit/{id}", productHandler.EditProductPage)
	r.Put("/admin/products/{id}", productHandler.EditProduct)
	r.Delete("/admin/products/{id}", productHandler.DeleteProduct)

	http.ListenAndServe(":8080", r)
}
