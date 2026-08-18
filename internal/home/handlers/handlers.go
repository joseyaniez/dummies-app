package handlers

import (
	"log"
	"net/http"

	"github.com/joseyanez/dummies-app/internal/home/views/pages"
	"github.com/joseyanez/dummies-app/internal/products/models"
	"github.com/joseyanez/dummies-app/internal/products/services"
	"github.com/joseyanez/dummies-app/internal/products/views/pages/public"
)

type HomeHandler struct {
	productService *services.ProductService
}

func NewHomeHandler(prodService *services.ProductService) *HomeHandler {
	return &HomeHandler{
		productService: prodService,
	}
}

func (h *HomeHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	pages.HomePage().Render(r.Context(), w)
}

func (h *HomeHandler) ProductsPage(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetProducts()
	if err != nil {
		log.Println("Error obtain products: " + err.Error())
		prods := []*models.Product{}
		public.PublicListProductsPage(prods).Render(r.Context(), w)
	}
	public.PublicListProductsPage(products).Render(r.Context(), w)
}
