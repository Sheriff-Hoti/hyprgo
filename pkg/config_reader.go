package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sheriff-Hoti/hyprgo/consts"
)

type Backend string

const (
	Swaybg    Backend = "swaybg"
	Hyprpaper Backend = "hyprpaper"
)

type Config struct {
	backend      Backend
	wallpaperDir string
}

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

func ReadConfigFile() (*Config, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	config_path := fmt.Sprintf("%v/%v", home_dir, consts.CONFIG_PATH)
	log.Println(config_path)
	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		log.Print("congrats 404")
		return &Config{
			backend:      Swaybg,
			wallpaperDir: "./",
		}, nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		log.Print(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Config{
		backend:      Swaybg,
		wallpaperDir: "./",
	}, nil
}
