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
		SELECT products.id, title, description, price, available, images.filename AS filename 
		FROM products 
		LEFT JOIN images 
	  	ON images.product_id = products.id
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var prod models.Product
		var filenameSql sql.NullString
		var filename string
		err = rows.Scan(
			&prod.Id,
			&prod.Title,
			&prod.Description,
			&prod.Price,
			&prod.Available,
			&filenameSql,
		)
		if err != nil {
			return nil, fmt.Errorf("Error scanning product: %w", err)
		}

		if filenameSql.Valid {
			filename = filenameSql.String
		}

		var exists bool = false
		for _, p := range products {
			if p.Id == prod.Id {
				exists = true
			}
		}
		if exists && filename != "" {
			products[len(products)-1].Images = append(products[len(products)-1].Images, filename)
		} else {
			if filename != "" {
				prod.Images = append(prod.Images, filename)
			}
			products = append(products, &prod)
		}
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
