package services

import (
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
	productId, err := s.productRepository.SaveProduct(productRequest.Title, productRequest.Description, price)

	filenames, errorsMap, err := storage.SaveImages(productRequest.Images)
	if err != nil {
		log.Println("Error for save images: " + err.Error())
		errors["image_form"] = "Error al guardar imágenes"
		return errors, nil
	}

	if len(errorsMap) > 0 {
		for _, errMap := range errorsMap {
			errors["image_form"] = errors["image_form"] + "; " + errMap
		}
		return errors, nil
	}

	err = s.productRepository.SaveProductImageFilenames(productId, filenames)
	if err != nil {
		log.Println("Error in SaveProductImageFilenames: " + err.Error())
		errors["image_form"] = errors["image_form"] + "; " + "Error al guadar ruta de imágenes"
	}

	return errors, nil
}
