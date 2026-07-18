package services

import "mime/multipart"

type ProductCreateRequest struct {
	Title       string
	Description string
	Price       string
	Images      []*multipart.FileHeader
}
