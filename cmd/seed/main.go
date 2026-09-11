package main

import (
	"log"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joseyanez/dummies-app/internal/auth/repositories"
	"github.com/joseyanez/dummies-app/internal/auth/services"
	"github.com/joseyanez/dummies-app/internal/database"
	prodRepo "github.com/joseyanez/dummies-app/internal/products/repositories"
)

func main() {
	db, err := database.OpenSQLiteDatabase("data/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	authRepo := repositories.NewAuthRepository(db)
	productRepo := prodRepo.NewProductRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	authService := services.NewAuthService(*authRepo, *sessionRepo)

	adminReq := services.CreateAdminRequest{
		Name:     "Gamer64XD",
		Password: "joseito99",
	}

	_, err = authService.CreateNewUser(&adminReq)
	if err != nil {
		log.Fatalf("Error to create admin: %v", err)
	}
	log.Println("Admin ", adminReq.Name+"created sucessfully")

	log.Println("Try to create 500 products with same image.")
	for range 0 {
		id, err := productRepo.SaveProduct(
			gofakeit.Name(),
			gofakeit.Paragraph(),
			gofakeit.Price(0.5, 250.0),
		)
		if err != nil {
			continue
		}
		idText := strconv.Itoa(id)
		productRepo.SaveProductImageFilenames(idText, []string{"sapito2.png"})
	}
}
