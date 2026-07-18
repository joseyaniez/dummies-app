package services

import (
	"strconv"

	"github.com/joseyanez/dummies-app/internal/products/repositories"
)

type ProductService struct {
	productRepository *repositories.ProductRepository
}

func NewProductService(productRepository *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) SaveProduct(productRequest ProductCreateRequest) (map[string]string, error) {
	errors := make(map[string]string)
	if productRequest.Title == "" {
		errors["title"] = "Debes colocar un título"
	}

	price, err := strconv.ParseFloat(productRequest.Price, 64)
	if err != nil || price <= 0 {
		errors["price"] = "Debes colocar un precio válido"
	}

	if len(errors) > 0 {
		return errors, nil
	}

	// Aquí iría la lógica para guardar el producto en la base de datos
	err = s.productRepository.SaveProduct(productRequest.Title, productRequest.Description, price)
	return nil, err
}
