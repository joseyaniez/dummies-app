package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	authHdl "github.com/joseyanez/dummies-app/internal/auth/handlers"
	authRepo "github.com/joseyanez/dummies-app/internal/auth/repositories"
	authServ "github.com/joseyanez/dummies-app/internal/auth/services"
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
	authRepository := authRepo.NewAuthRepository(db)
	sessionRepository := authRepo.NewSessionRepository(db)

	productService := services.NewProductService(productRepository)
	authService := authServ.NewAuthService(*authRepository, *sessionRepository)

	productHandler := handlers.NewProductHandler(productService)
	authHandler := authHdl.NewAuthHandler(*authService)

	r.Get("/admin/products/create", productHandler.New)
	r.Get("/admin/products", productHandler.List)
	r.Post("/admin/products", productHandler.Create)
	r.Get("/admin/products/{id}", productHandler.Show)
	r.Get("/admin/products/edit/{id}", productHandler.Edit)
	r.Put("/admin/products/{id}", productHandler.Update)
	r.Delete("/admin/products/{id}", productHandler.Delete)
	r.Get("/admin/login", authHandler.Login)

	r.Get("/products", productHandler.PublicList)

	http.ListenAndServe(":8080", r)
}
