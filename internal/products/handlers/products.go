package handlers

import (
	"fmt"
	"maps"
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
	errors := make(map[string]string)
	productRequest := services.ProductCreateRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Price:       r.FormValue("price"),
		Images:      []string{}, // Se asignará después de procesar las imágenes
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors["form"] = "Error interno al enviar el formulario, intende de nuevo más tarde"
		pages.CreateProductPage(productRequest, errors).Render(r.Context(), w)
		return
	}

	maperrors, err := h.productService.SaveProduct(productRequest)
	if err != nil {
		errors["form"] = "Error interno al guardar el producto, intente de nuevo más tarde"
		pages.CreateProductPage(productRequest, errors).Render(r.Context(), w)
		return
	}

	maps.Copy(errors, maperrors)

	if len(errors) > 0 {
		pages.CreateProductPage(productRequest, errors).Render(r.Context(), w)
		return
	}

	images := r.MultipartForm.File["files"]

	for _, image := range images {
		fmt.Printf("Received file: %s\n", image.Filename)
	}

	// Aquí iría la lógica para guardar el producto en la base de datos y manejar las imágenes

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
