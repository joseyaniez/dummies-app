package services

import "mime/multipart"

type ProductCreateRequest struct {
	Title       string
	Description string
	Price       string
	Images      []*multipart.FileHeader
}

type ProductEditRequest struct {
	Title           string
	Description     string
	Price           string
	ImagesForDelete []string
	Images          []*multipart.FileHeader
}
