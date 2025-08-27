package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sheriff-Hoti/hyprgo/data"
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

	config := GetDefaultConfigVals()

	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist and if it does not exists just return the defaults

		return config, nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func GetDefaultConfigPath() string {
	const (
		xdgConfigHome = "XDG_CONFIG_HOME"
	)

	if val, ok := os.LookupEnv(xdgConfigHome); ok {
		return filepath.Join(val, "hyprgo", "config.json")
	}

	// fallback to $HOME/.config/hyprgo/config.json
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		return filepath.Join(".", "config.json")
	}

	return filepath.Join(home, ".config", "hyprgo", "config.json")
}

func GetDefaultConfigVals() *Config {

	return &Config{
		Backend:       "swaybg",
		Wallpaper_dir: "",
		Data_dir:      data.GetDefaultDataPath(),
	}
}
