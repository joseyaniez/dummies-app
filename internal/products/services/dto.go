package services

type ProductCreateRequest struct {
	Title       string
	Description string
	Price       string
	Images      []string
}
