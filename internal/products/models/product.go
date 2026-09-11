package models

import (
	"fmt"
)

type Product struct {
	Id          string
	Title       string
	Description string
	Price       float64
	Available   bool
	Images      map[string][]string
}

func (p *Product) ImagesToSrcs(dir string) map[string]string {
	srcs := make(map[string]string)

	for name, sizes := range p.Images {
		srcImage := ""
		if len(sizes) == 3 {
			srcImage = fmt.Sprintf("%s_%s.webp", name, sizes[1])
		} else {
			srcImage = fmt.Sprintf("%s_%s.webp", name, sizes[len(sizes)-1])
		}

		srcset := ""
		for _, size := range sizes {
			filename := fmt.Sprintf("%s_%s.webp", name, size)
			if dir == "" {
				srcset += fmt.Sprintf("%s %sw, ", filename, size)
			} else {
				srcset += fmt.Sprintf("%s/%s %sw, ", dir, filename, size)
			}
		}
		srcs[srcImage] = srcset
	}
	return srcs
}

func (p *Product) ImagesMapToList() []string {
	imagesList := []string{}
	for imageName, sizes := range p.Images {
		for _, size := range sizes {
			filename := fmt.Sprintf("%s_%s.webp", imageName, size)
			imagesList = append(imagesList, filename)
		}
	}
	return imagesList
}

func (p *Product) GetOneAndOtherImages(dir string) (map[string]string, map[string]string) {
	firstImage := make(map[string]string)
	srcs := make(map[string]string)

	isFirst := true
	for name, sizes := range p.Images {
		srcImage := ""
		if len(sizes) == 3 {
			srcImage = fmt.Sprintf("%s_%s.webp", name, sizes[1])
		} else {
			srcImage = fmt.Sprintf("%s_%s.webp", name, sizes[len(sizes)-1])
		}

		srcset := ""
		for _, size := range sizes {
			filename := fmt.Sprintf("%s_%s.webp", name, size)
			if dir == "" {
				srcset += fmt.Sprintf("%s %sw, ", filename, size)
			} else {
				srcset += fmt.Sprintf("%s/%s %sw, ", dir, filename, size)
			}
		}
		srcs[srcImage] = srcset
		if isFirst {
			firstImage[srcImage] = srcset
			isFirst = false
		}
	}
	return firstImage, srcs
}
