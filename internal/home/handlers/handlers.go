package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	query := r.URL.Query()
	page := 1
	if value := query.Get("page"); value != "" {
		if p, err := strconv.Atoi(value); err == nil && p >= 1 {
			page = p
		}
	}
	products, err := h.productService.GetProducts(page)
	if err != nil {
		log.Println("Error obtain products: " + err.Error())
		prods := []*models.Product{}
		public.PublicListProductsPage(prods, page).Render(r.Context(), w)
		return
	}
	public.PublicListProductsPage(products, page).Render(r.Context(), w)
}

func (h *HomeHandler) ProductDetailPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.productService.GetProduct(id)
	if err != nil {
		public.PublicViewProductPage(nil).Render(r.Context(), w)
		return
	}
	public.PublicViewProductPage(product).Render(r.Context(), w)
}
