package storage

import (
	"fmt"
	imagePkg "image"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

var widths = []int{400, 800, 1200}

func NewSaveImages(images []*multipart.FileHeader) (filenames, imageErrors []string, generalError error) {
	uploadDir := "web/static/images/"
	var wgImages sync.WaitGroup
	wgImages.Add(len(images))

	imageErrorsChan := make(chan string, len(images))

	for _, image := range images {
		imageFile, err := image.Open()
		if err != nil {
			imageErrors = append(imageErrors, fmt.Sprintf("Error al abrir la imagen: %s", image.Filename))
			log.Println("Failed to open image: ", err)
		}

		go func(imageFile multipart.File, filename string) {
			defer wgImages.Done()
			defer imageFile.Close()
			// Validar imagen: (formato permitido y dimensiones máximas)
			imgConfig, format, err := imagePkg.DecodeConfig(imageFile)
			if err != nil {
				imageErrorsChan <- fmt.Sprintf("Error al decodificar la imagen %s", filename)
				log.Println("Error while image decode: ", err)
				return
			}
			switch format {
			case "jpeg", "png", "webp":
				// Imagen válida
			default:
				imageErrorsChan <- fmt.Sprintf("Formato de imagen %s no válido", filename)
				return
			}
			// Obtener el ancho de la imagen
			imageWidth := imgConfig.Width

			// Determinar variantes de la imagen necesarias
			for _, targetWidth := range widths {
				if targetWidth < imageWidth {
					// generar nueva imagen
				}
			}
			// Redimensionar manteniendo proporción
			// Codificar como WebP
			// Guardar variantes en disco
		}(imageFile, image.Filename)
	}

	wgImages.Wait()
}

func SaveImages(images []*multipart.FileHeader) (filenames, errors []string) {
	uploadDir := "web/static/images/"

	for _, image := range images {
		src, err := image.Open()
		if err != nil {
			log.Println("Failed to open image: " + err.Error())
			errors = append(errors, "Failed to open image: "+image.Filename)
			continue
		}

		func() {
			defer src.Close()

			// Generate unique name for the image
			ext := filepath.Ext(image.Filename)
			filename := uuid.NewString() + ext
			dstPath := filepath.Join(uploadDir, filename)

			// Create a destination file
			dst, err := os.Create(dstPath)
			if err != nil {
				log.Println("Failed to create destination file: " + err.Error())
				errors = append(errors, "Failed to create destination file: "+filename)
				return
			}
			defer dst.Close()

			// Copy the uploaded image to the destination file
			_, err = io.Copy(dst, src)
			if err != nil {
				log.Println("Failed to copy image in destination: " + err.Error())
				errors = append(errors, "Failed to save image: "+filename)
				return
			}

			filenames = append(filenames, filename)
		}()
	}
	return
}

func DeleteImages(images []string) (errors []string) {
	uploadDir := "web/static/images/"

	for _, image := range images {
		path := filepath.Join(uploadDir, image)

		fmt.Println(path)

		err := os.Remove(path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Println("Image does not exist:", path)
				errors = append(errors, "Image does not exist: "+image)
				continue
			}

			log.Println("Failed to delete image:", err)
			errors = append(errors, "Failed to delete image: "+image)
		}
	}

	return
}
