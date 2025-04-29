package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

func smth() {}

func GetWallpapers(dir string) (filenames []string, erro error) {
	file_names := make([]string, 0, 10)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filename := fmt.Sprintf("%v/%v", dir, info.Name())
			ext := filepath.Ext(filename)
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				file_names = append(file_names, filename)
			}
		}
		return nil
	})

	return file_names, err
}
