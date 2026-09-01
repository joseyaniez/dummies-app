package helpers

import (
	"fmt"
	"strings"
)

func FilenamesToMap(filenames []string) map[string][]string {
	result := make(map[string][]string)

	for _, filename := range filenames {
		name := strings.TrimSuffix(filename, ".webp")

		index := strings.LastIndex(name, "_")
		if index == -1 {
			continue
		}

		baseName := name[:index]
		size := name[index+1:]

		result[baseName] = append(result[baseName], size)
	}

	return result
}

func MapImagesToFilenames(images map[string][]string) []string {
	var filenames []string

	for name, sizes := range images {
		for _, size := range sizes {
			filename := fmt.Sprintf("%s_%s.webp", name, size)
			filenames = append(filenames, filename)
		}
	}

	return filenames
}
