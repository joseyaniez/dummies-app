package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	products, err := h.productService.GetProducts()
	if err != nil {
		log.Println("Error obtain products: " + err.Error())
		pages.ListProductsPage(nil).Render(r.Context(), w)
		return
	}
	pages.ListProductsPage(products).Render(r.Context(), w)
}

func (h *ProductHandler) ViewProduct(w http.ResponseWriter, r *http.Request) {
	// obtener el producto con el id
	id := chi.URLParam(r, "id")
	prod, err := h.productService.GetProduct(id)
	if err != nil {
		log.Println(err)
		pages.ViewProductPage(nil).Render(r.Context(), w)
		return
	}
	pages.ViewProductPage(prod).Render(r.Context(), w)
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

func (h *ProductHandler) EditProductPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	errors := make(map[string]string)
	product, err := h.productService.GetProduct(id)
	if err != nil {
		pages.EditProductPage(nil, services.ProductEditRequest{}, errors)
	}
	pages.EditProductPage(product, services.ProductEditRequest{}, errors).Render(r.Context(), w)
}

func (h *ProductHandler) EditProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errorForm := map[string]string{"form": "Error interno al enviar el formulario, intente de nuevo más tarde"}
		pages.EditProductPage(nil, services.ProductEditRequest{}, errorForm).Render(r.Context(), w)
		return
	}

	id := chi.URLParam(r, "id")
	prod, err := h.productService.GetProduct(id)
	if err != nil {
		errorForm := map[string]string{"form": "El producto a eliminar no existe"}
		pages.EditProductPage(nil, services.ProductEditRequest{}, errorForm).Render(r.Context(), w)
		return
	}

	prodRequest := services.ProductEditRequest{
		Title:           r.FormValue("title"),
		Description:     r.FormValue("description"),
		Price:           r.FormValue("price"),
		ImagesForDelete: r.Form["images_filenames"],
		Images:          r.MultipartForm.File["files"],
	}

	validationErrors, err := h.productService.EditProduct(id, prodRequest)

	if len(validationErrors) > 0 {
		pages.EditProductPage(prod, prodRequest, nil).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
