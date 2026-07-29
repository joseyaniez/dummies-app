package main

import (
	"log"

	"github.com/joseyanez/dummies-app/internal/auth/repositories"
	"github.com/joseyanez/dummies-app/internal/auth/services"
	"github.com/joseyanez/dummies-app/internal/database"
)

func main() {
	db, err := database.OpenSQLiteDatabase("data/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(*authRepo)

	adminReq := services.CreateAdminRequest{
		Name:     "Gamer64XD",
		Password: "STXD3t*#8484",
	}
	_, err = authService.CreateNewUser(&adminReq)
	if err != nil {
		log.Fatalf("Error to create admin: %v", err)
	}
	log.Println("Admin ", adminReq.Name+"created sucessfully")
}
