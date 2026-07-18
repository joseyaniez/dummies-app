package repositories

import (
	"database/sql"
	"time"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) SaveProduct(title, description string, price float64) error {
	query := `
		INSERT INTO products (title, description, price, created_at, updated_at) VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(
		query,
		title, description, price, time.Now(), time.Now(),
	)

	return err
}
