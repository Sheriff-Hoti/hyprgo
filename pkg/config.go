package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Backend       string `json:"backend"`
	Wallpaper_dir string `json:"wallpaper_dir"`
	Data_dir      string `json:"data_dir"`
}

func GetWallpapers(dir string) (filenames []string, erro error) {
	dir_env_expanded := os.ExpandEnv(dir)
	entries, err := os.ReadDir(dir_env_expanded)
	if err != nil {
		return nil, err
	}

	file_names := make([]string, 0, 10)

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				fullPath := filepath.Join(dir_env_expanded, entry.Name())

				file_names = append(file_names, fullPath)
			}
		}
	}

	return file_names, nil
}

func ReadConfigFile(config_path string) (*Config, error) {

	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist and if it does not exists just return the defaults

		return GetDefaultConfigVals(), nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GetDefaultConfigPath() string {
	xdg_config_home := "XDG_CONFIG_HOME"
	home := "HOME"

	if _, ok := os.LookupEnv(xdg_config_home); ok {
		return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", xdg_config_home), "hyprgo", "config.json"))
	}
	return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", home), ".config", "hyprgo", "config.json"))

}

func GetDefaultConfigVals() *Config {

	return &Config{
		Backend:       "swaync",
		Wallpaper_dir: os.ExpandEnv("$HOME"),
		Data_dir:      GetDefaultDataPath(),
	}
}

// after this it will be ging to other validator metho
// first the default config and then each config function will get
// the map value and pull their own key
// also add a map func i guess ??
