package handlers

import (
	"net/http"

	"github.com/joseyanez/dummies-app/internal/products/services"
	"github.com/joseyanez/dummies-app/internal/products/views/pages"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pages.ListProductsPage().Render(r.Context(), w)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	productRequest := services.ProductCreateRequest{}
	pages.CreateProductPage(productRequest, nil).Render(r.Context(), w)
}

func (h *ProductHandler) SaveProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errorForm := map[string]string{"form": "Error interno al enviar el formulario, intente de nuevo más tarde"}
		pages.CreateProductPage(services.ProductCreateRequest{}, errorForm).Render(r.Context(), w)
		return
	}

	productRequest := services.ProductCreateRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       r.FormValue("price"),
		Images:      r.MultipartForm.File["files"],
	}

	validationErrors, err := h.productService.SaveProduct(productRequest)
	if err != nil {
		errorForm := map[string]string{"form": "No se pudo guardar el producto, intente de nuevo más tarde"}
		pages.CreateProductPage(productRequest, errorForm).Render(r.Context(), w)
		return
	}

	if len(validationErrors) > 0 {
		pages.CreateProductPage(productRequest, validationErrors).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
