package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/joseyanez/dummies-app/internal/products/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProducts() ([]*models.Product, error) {
	query := `
		SELECT id, title, description, price, available FROM products
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var prod models.Product
		err = rows.Scan(
			&prod.Id,
			&prod.Title,
			&prod.Description,
			&prod.Price,
			&prod.Available,
		)
		if err != nil {
			return nil, fmt.Errorf("Error scanning product: %w", err)
		}
		products = append(products, &prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating scanning products: %w", err)
	}

	return products, nil
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
