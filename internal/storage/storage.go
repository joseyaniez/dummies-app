package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveImages(images []*multipart.FileHeader) ([]string, []string, error) {
	var filenames []string
	var errors []string

	uploadDir := "web/static/images/"

	for _, image := range images {
		src, err := image.Open()
		if err != nil {
			errors = append(errors, "Failed to open image: "+image.Filename)
			src.Close()
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
				errors = append(errors, "Failed to create destination file: "+filename)
				return
			}
			defer dst.Close()

			// Copy the uploaded image to the destination file
			_, err = io.Copy(dst, src)
			if err != nil {
				errors = append(errors, "Failed to save image: "+filename)
				return
			}

			filenames = append(filenames, filename)
		}()
	}
	return filenames, errors, nil
}
