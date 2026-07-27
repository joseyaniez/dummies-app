package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/joseyanez/dummies-app/internal/products/models"
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

func (s *ProductService) GetProducts() ([]*models.Product, error) {
	products, err := s.productRepository.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	product, err := s.productRepository.FindProduct(id)
	if err != nil {
		return nil, err
	}
	return product, nil
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

	err = s.productRepository.SaveProductImageFilenames(fmt.Sprintf("%d", productId), filenames)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ProductService) EditProduct(id string, productRequest ProductEditRequest) (map[string]string, error) {
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

	err = s.productRepository.EditProduct(id, productRequest.Title, productRequest.Description, price)
	if err != nil {
		return nil, err
	}

	// eliminar físicamente las imágenes correspondientes
	errs := storage.DeleteImages(productRequest.ImagesForDelete)
	for _, err := range errs {
		log.Printf("Error al eliminar imágenes: %s", err)
	}

	// eliminar las rutas de imagen de la base de datos
	err = s.productRepository.DeleteProductImageFilenames(productRequest.ImagesForDelete)
	if err != nil {
		return nil, err
	}

	// insertar físicamente las nuevas imágenes
	filenames, errorsStr := storage.SaveImages(productRequest.Images)

	if len(errorsStr) > 0 {
		error := ""
		for _, err := range errorsStr {
			error = error + err
		}
		return nil, errors.New(error)
	}

	// agregar las rutas de imagen en la base de datos
	err = s.productRepository.SaveProductImageFilenames(id, filenames)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
