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

func (r *ProductRepository) SaveProduct(title, description string, price float64) (int, error) {
	query := `
		INSERT INTO products (title, description, price, created_at, updated_at) VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.DB.Exec(
		query,
		title, description, price, time.Now(), time.Now(),
	)

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r *ProductRepository) SaveProductImageFilenames(productId int, filenames []string) error {
	query := `
		INSERT INTO images (product_id, filename, created_at, updated_at) VALUES (?, ?, ?, ?)
	`

	for _, filename := range filenames {
		_, err := r.DB.Exec(query, productId, filename, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
