package handlers

import (
	"net/http"

	"github.com/joseyanez/dummies-app/internal/products/views/pages"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pages.ListProductsPage().Render(r.Context(), w)
}
