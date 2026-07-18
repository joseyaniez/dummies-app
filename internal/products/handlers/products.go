package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/joseyanez/dummies-app/internal/products/views/pages"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pages.ListProductsPage().Render(r.Context(), w)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	pages.CreateProductPage(nil, nil).Render(r.Context(), w)
}

func (h *ProductHandler) SaveProduct(w http.ResponseWriter, r *http.Request) {
	values := make(map[string]string)
	errors := make(map[string]string)

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors["form"] = "Error interno al enviar el formulario, intende de nuevo más tarde"
		pages.CreateProductPage(values, errors).Render(r.Context(), w)
		return
	}

	values["title"] = r.FormValue("title")
	values["description"] = r.FormValue("description")

	if values["title"] == "" {
		errors["title"] = "Debes colocar un título"
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil || price <= 0 {
		errors["price"] = "Debes colocar un precio válido"
	}
	values["price"] = r.FormValue("price")

	images := r.MultipartForm.File["files"]

	if len(errors) > 0 {
		pages.CreateProductPage(values, errors).Render(r.Context(), w)
		return
	}

	for _, image := range images {
		fmt.Printf("Received file: %s\n", image.Filename)
	}

	// Aquí iría la lógica para guardar el producto en la base de datos y manejar las imágenes

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
