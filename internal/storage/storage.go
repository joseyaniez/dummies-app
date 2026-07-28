package storage

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

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
