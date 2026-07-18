package services

import "strconv"

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
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

	return nil, nil
}
