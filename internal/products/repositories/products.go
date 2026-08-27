package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/joseyanez/dummies-app/internal/products/models"
)

var ErrProductNotFound = errors.New("No se encontró el producto")

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) FindProduct(id string) (*models.Product, error) {
	query := `
	  SELECT id, title, description, price, available FROM products WHERE id = ?
	`
	var prod models.Product

	row := r.DB.QueryRow(query, id)
	err := row.Scan(
		&prod.Id,
		&prod.Title,
		&prod.Description,
		&prod.Price,
		&prod.Available,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Obtener las imágenes
	queryImages := `
	  SELECT filename FROM images WHERE product_id = ?
	`
	rows, err := r.DB.Query(queryImages, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		err := rows.Scan(
			&filename,
		)
		if err != nil {
			return nil, err
		}
		prod.Images = append(prod.Images, filename)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &prod, nil
}

func (r *ProductRepository) GetProducts(page int) ([]*models.Product, error) {
	offset := (page - 1) * 10
	if page == 0 {
		offset = 0
	}
	query := `
		SELECT products.id, title, description, price, available, images.filename AS filename 
		FROM products 
		LEFT JOIN images 
	  	ON images.product_id = products.id
		LIMIT ? OFFSET ?
	`
	rows, err := r.DB.Query(query, 10, offset)
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

func (r *ProductRepository) SaveProductImageFilenames(productId string, filenames []string) error {
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

func (r *ProductRepository) DeleteProductImageFilenames(filenames []string) error {
	placeholders := make([]string, len(filenames))
	args := make([]any, len(filenames))

	for i, filename := range filenames {
		placeholders[i] = "?"
		args[i] = filename
	}

	query := fmt.Sprintf(
		"DELETE FROM images WHERE filename IN (%s)",
		strings.Join(placeholders, ","),
	)

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) EditProduct(id, title, description string, price float64) error {
	query := `
	  UPDATE products
		SET title = ?,
		    description = ?,
				price = ?,
				updated_at = ?
		WHERE id = ?
	`
	_, err := r.DB.Exec(query, title, description, price, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
	query := `
	  DELETE FROM products WHERE id = ?
	`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
