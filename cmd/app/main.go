package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	authHdl "github.com/joseyanez/dummies-app/internal/auth/handlers"
	"github.com/joseyanez/dummies-app/internal/auth/middlewares"
	authRepo "github.com/joseyanez/dummies-app/internal/auth/repositories"
	authServ "github.com/joseyanez/dummies-app/internal/auth/services"
	"github.com/joseyanez/dummies-app/internal/database"
	homeHdl "github.com/joseyanez/dummies-app/internal/home/handlers"
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
	homeHandler := homeHdl.NewHomeHandler(productService)
	authHandler := authHdl.NewAuthHandler(*authService)

	authMiddleware := middlewares.NewAuthMiddleware(*sessionRepository)

	r.Get("/", homeHandler.HomePage)
	r.Get("/products", homeHandler.ProductsPage)
	r.Get("/products/{id}", homeHandler.ProductDetailPage)

	r.Route("/admin", func(r chi.Router) {
		r.Use(authMiddleware.Auth)
		r.Get("/products", productHandler.List)
		r.Post("/products", productHandler.Create)
		r.Get("/products/create", productHandler.New)
		r.Get("/products/edit/{id}", productHandler.Edit)
		r.Put("/products/{id}", productHandler.Update)
		r.Get("/products/{id}", productHandler.Show)
		r.Delete("/products/{id}", productHandler.Delete)
		r.Post("/logout", authHandler.Logout)
	})

	r.With(authMiddleware.Guest).Get("/admin/login", authHandler.LoginPage)
	r.With(authMiddleware.Guest).Post("/admin/login", authHandler.Login)

	http.ListenAndServe(":8080", r)
}
