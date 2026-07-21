package services

import (
	"errors"
	"log"
	"strconv"

	"github.com/joseyanez/dummies-app/internal/products/repositories"
	"github.com/joseyanez/dummies-app/internal/storage"
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
	validationErrors := make(map[string]string)
	if productRequest.Title == "" {
		validationErrors["title"] = "Debes colocar un título"
	}

	price, err := strconv.ParseFloat(productRequest.Price, 64)
	if err != nil || price <= 0 {
		validationErrors["price"] = "Debes colocar un precio válido"
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	// Aquí iría la lógica para guardar el producto en la base de datos
	productId, err := s.productRepository.SaveProduct(productRequest.Title, productRequest.Description, price)
	if err != nil {
		log.Println("Error saving product in repository: " + err.Error())
		return nil, err
	}

	filenames, errorsStr := storage.SaveImages(productRequest.Images)

	if len(errorsStr) > 0 {
		error := ""
		for _, err := range errorsStr {
			error = error + err
		}
		return nil, errors.New(error)
	}

	err = s.productRepository.SaveProductImageFilenames(productId, filenames)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
